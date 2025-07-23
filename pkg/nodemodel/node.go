package nodemodel

import (
	"strings"
	"sync"
	"time"
)

type Node struct {
	ID           int
	NomadID      string
	Address      string
	IsLeader     bool
	Peers        map[string]time.Time
	LastElection int
	mu           sync.Mutex
}

func NewNode(id int, address string) *Node {
	return &Node{
		ID:           id,
		Address:      address,
		Peers:        make(map[string]time.Time),
		LastElection: 0,
		mu:           sync.Mutex{},
	}
}

func (n *Node) GetID() int {
	return n.ID
}

func (n *Node) GetAddress() string {
	return n.Address
}

func (n *Node) GetIsLeader() int {
	return n.LastElection
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

func (n *Node) BecomeLeader() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.IsLeader = true
}

func (n *Node) BecomeReplica() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.IsLeader = false
}
