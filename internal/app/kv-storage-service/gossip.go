package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) Gossip(
	ctx context.Context,
	req *desc.GossipRequest,
) (*desc.GossipResponse, error) {
	s.nodeService.HandleGossip(req.GetNode())
	return &desc.GossipResponse{}, nil
}
