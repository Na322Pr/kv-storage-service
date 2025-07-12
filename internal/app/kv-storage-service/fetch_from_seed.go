package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) FetchFromSeed(
	ctx context.Context,
	req *desc.FetchFromSeedRequest,
) (*desc.FetchFromSeedResponse, error) {
	peers := s.nodeService.SeedHandleAddPeer(req.GetAddress())

	return &desc.FetchFromSeedResponse{
		Peers: peers,
	}, nil
}
