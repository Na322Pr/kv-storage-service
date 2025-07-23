package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) SetLeader(ctx context.Context, req *desc.NewLeaderRequest) (*desc.NewLeaderResponse, error) {
	s.leService.SetLeader(int(req.Id))
	return &desc.NewLeaderResponse{}, nil
}
