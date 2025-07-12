package service

import (
	"context"
	"fmt"
	"github.com/Na322Pr/kv-storage-service/internal/model"
	pb "github.com/Na322Pr/kv-storage-service/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
	"sync"
	"time"
)

type NodeService struct {
	node   *model.Node
	logger *zap.Logger
}

func NewNodeService(id int, address string, isSeed bool, logger *zap.Logger) *NodeService {
	return &NodeService{
		node:   model.NewNode(id, address, isSeed),
		logger: logger,
	}
}

func (s *NodeService) Run(ctx context.Context, seeds []string) error {
	if err := s.JoinCluster(ctx, seeds); err != nil {
		s.logger.Fatal("failed to join cluster", zap.Error(err))
		return err
	}

	//go func() {
	//	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond) // Random delay
	//	s.StartLeaderElection()
	//}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from gossip loop panic: %v", r)
			}
		}()
		s.StartGossipLoop(ctx)
	}()

	return nil
}

func (s *NodeService) JoinCluster(ctx context.Context, seeds []string) error {
	if !s.node.IsSeed {
		for _, seed := range seeds {
			peers := s.FetchPeersFromSeed(ctx, seed)
			s.logger.Info(fmt.Sprintf("joining cluster with seed: %s", strings.Join(seeds, ",")))
			if len(peers) > 0 {
				s.node.UpdatePeers(peers)
				break
			}
		}
	}

	return nil
}

func (s *NodeService) StartGossipLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		s.logger.Info("Starting gossip loop...")
		s.BroadcastPresence(ctx)
		s.node.CleanUpStalePeers(30 * time.Second)

		activePeers := s.node.GetPeersString()
		if len(activePeers) == 0 {
			s.logger.Info("no active peers")
		} else {
			s.logger.Info(fmt.Sprintf("active peers: %s", activePeers))
		}
	}
}

func (s *NodeService) BroadcastPresence(ctx context.Context) {
	peers := s.node.GetPeers()

	wg := &sync.WaitGroup{}
	for _, peer := range peers {
		wg.Add(1)
		go s.SendGossip(ctx, wg, peer)
	}
	wg.Wait()
}

func (s *NodeService) SendGossip(ctx context.Context, wg *sync.WaitGroup, peer string) {
	defer wg.Done()

	if peer == "" {
		s.logger.Warn("attempted to gossip with empty peer address")
		return
	}

	s.logger.Info("sending gossip message", zap.String("peer", peer))

	conn, err := grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.logger.Error("failed to connect to peer", zap.String("peer", peer), zap.Error(err))
		return
	}
	defer conn.Close()

	client := pb.NewKeyValueStorageClient(conn)
	_, err = client.Gossip(ctx, &pb.GossipRequest{
		Node: s.node.GetAddress(),
	})
	if err != nil {
		s.logger.Error(
			"failed to send gossip message",
			zap.String("peer", peer),
			zap.Error(err),
		)
		return
	}

	s.node.UpdatePeer(peer)
}

func (s *NodeService) FetchPeersFromSeed(ctx context.Context, seed string) []string {
	conn, _ := grpc.NewClient(seed, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := pb.NewKeyValueStorageClient(conn)
	resp, err := client.FetchFromSeed(ctx, &pb.FetchFromSeedRequest{
		Address: s.node.GetAddress(),
	})
	if err != nil {
		//TODO:
	}

	if resp != nil {
		return resp.Peers
	}

	return nil
}

func (s *NodeService) HandleGossip(peer string) {
	s.logger.Info(fmt.Sprintf("handling gossip message from %s", peer))
	s.node.UpdatePeer(peer)
}

func (s *NodeService) GetTerm() int {
	return s.node.ElectionTerm
}

func (s *NodeService) GetPeers() []string {
	return s.node.GetPeers()
}

func (s *NodeService) SeedHandleAddPeer(address string) []string {
	peers := s.node.GetPeers()
	s.node.UpdatePeer(address)
	return peers
}

func (s *NodeService) HandleVoteRequest(candidateAddress string, term int) bool {
	if !s.node.CheckVoteAvailable(term) {
		return false
	}

	s.node.GrantVote(candidateAddress, term)
	return true
}

func (s *NodeService) StartLeaderElection() {
	s.node.StartLeaderElection()

	s.logger.Info("starting leader election",
		zap.Int("term", s.node.GetTerm()),
		zap.String("candidate", s.node.GetAddress()))

	voteCh := make(chan bool, 1)

	var (
		votesNeeded   = len(s.node.GetPeers())/2 + 1
		votesReceived = 1
		wg            sync.WaitGroup
	)

	for _, peer := range s.node.GetPeers() {
		wg.Add(1)
		go func(peer string) {
			defer wg.Done()

			conn, err := grpc.NewClient(
				peer,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				s.logger.Error("failed to connect to peer for vote request",
					zap.String("peer", peer),
					zap.Error(err))
				return
			}
			defer conn.Close()

			client := pb.NewKeyValueStorageClient(conn)
			resp, err := client.LeaderVote(context.Background(), &pb.LeaderVoteRequest{
				CandidateAddress: s.node.GetAddress(),
				Term:             int64(s.node.ElectionTerm),
			})

			if err != nil {
				s.logger.Error("vote request failed",
					zap.String("peer", peer),
					zap.Error(err))
				return
			}

			if resp.VoteGranted {
				voteCh <- true
			}
		}(peer)
	}

	// Wait for votes or timeout
	go func() {
		for vote := range voteCh {
			if vote {
				votesReceived++
			}

			if votesReceived >= votesNeeded {
				s.becomeLeader()
				return
			}
		}
	}()
}

func (s *NodeService) becomeLeader() {
	s.node.BecomeLeader()

	s.logger.Info("became leader",
		zap.Int("term", s.node.GetTerm()),
		zap.String("leader", s.node.GetAddress()),
	)
}
