package nodemodel

import (
	"strings"
	"sync"
	"time"
)

type Node struct {
	ID              int
	Address         string
	LeaderAddress   string
	IsLeader        bool
	IsSeed          bool
	Peers           map[string]time.Time
	ElectionTerm    int
	LastUpdateIndex int64
	VotedFor        string
	mu              sync.Mutex
}

func NewNode(id int, address string, isSeed bool) *Node {
	return &Node{
		ID:      id,
		Address: address,
		IsSeed:  isSeed,
		Peers:   make(map[string]time.Time),
	}
}

func (n *Node) GetAddress() string {
	return n.Address
}

func (n *Node) GetPeers() []string {
	n.mu.Lock()
	defer n.mu.Unlock()

	peers := make([]string, 0, len(n.Peers))
	for key := range n.Peers {
		peers = append(peers, key)
	}

	return peers
}

func (n *Node) GetPeersString() string {
	peers := n.GetPeers()
	return strings.Join(peers, ",")
}

func (n *Node) UpdatePeer(peer string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Peers[peer] = time.Now()
}

func (n *Node) UpdatePeers(peer []string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, peer := range peer {
		n.Peers[peer] = time.Now()
	}
}

func (n *Node) RemovePeer(peer string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.Peers, peer)
}

func (n *Node) CleanUpStalePeers(duration time.Duration) {
	n.mu.Lock()
	defer n.mu.Unlock()

	for peer := range n.Peers {
		if time.Since(n.Peers[peer]) > duration {
			delete(n.Peers, peer)
		}
	}
}

func (n *Node) GetTerm() int {
	return n.ElectionTerm
}

func (n *Node) CheckVoteAvailable(term int) bool {
	if n.ElectionTerm >= term {
		return false
	}

	return true
}

func (n *Node) GrantVote(votedFor string, term int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.VotedFor = votedFor
	n.ElectionTerm = term
}

func (n *Node) StartLeaderElection() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.ElectionTerm = n.ElectionTerm + 1
	n.VotedFor = n.GetAddress()
}

func (n *Node) BecomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.IsLeader = true
	n.LeaderAddress = n.GetAddress()
}

func (n *Node) GetLeaderAddress() string {
	return n.LeaderAddress
}
