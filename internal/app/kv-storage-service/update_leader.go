package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) UpdateLeader(ctx context.Context, req *desc.UpdateLeaderRequest) (*desc.UpdateLeaderResponse, error) {
	s.leService.SetLeader(req.NomadId)
	return &desc.UpdateLeaderResponse{}, nil
}
