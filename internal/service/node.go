package service

import (
	"context"
	"fmt"
	pb "github.com/Na322Pr/kv-storage-service/pkg/api"
	"github.com/Na322Pr/kv-storage-service/pkg/nodemodel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
	"sync"
	"time"
)

type NodeService struct {
	node   *nodemodel.Node
	logger *zap.Logger
}

func NewNodeService(node *nodemodel.Node, logger *zap.Logger) *NodeService {
	return &NodeService{
		node:   node,
		logger: logger,
	}
}

func (s *NodeService) Run(ctx context.Context, seeds []string) error {
	if err := s.JoinCluster(ctx, seeds); err != nil {
		s.logger.Fatal("failed to join cluster", zap.Error(err))
		return err
	}
	
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

	for _, seed := range seeds {
		peers := s.FetchPeersFromSeed(ctx, seed)
		s.logger.Info(fmt.Sprintf("joining cluster with seed: %s", strings.Join(seeds, ",")))
		if len(peers) > 0 {
			s.node.UpdatePeers(peers)
			break
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

func (s *NodeService) GetPeers() []string {
	return s.node.GetPeers()
}

func (s *NodeService) SeedHandleAddPeer(address string) []string {
	peers := s.node.GetPeers()
	s.node.UpdatePeer(address)
	return peers
}

func (s *NodeService) becomeLeader() {
	s.node.BecomeLeader()
	s.logger.Info("became leader")
}

func (s *NodeService) becomeReplica() {
	s.node.BecomeReplica()
	s.logger.Info("became replica")
}
