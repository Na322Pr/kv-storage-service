package model

type Node struct {
	id       string
	nomadID  string
	address  string
	isLeader bool
}

func NewNode(id, nomadID, address string) *Node {
	return &Node{
		id:       id,
		nomadID:  nomadID,
		address:  address,
		isLeader: true,
	}
}

func (node *Node) ID() string {
	return node.id
}

func (node *Node) NomadID() string {
	return node.nomadID
}

func (node *Node) Address() string {
	return node.address
}

func (node *Node) IsLeader() bool {
	return node.isLeader
}

func (node *Node) SetLeader(isLeader bool) {
	node.isLeader = isLeader
}
