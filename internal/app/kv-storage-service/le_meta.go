package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) LeMeta(ctx context.Context, req *desc.LeMetaRequest) (*desc.LeMetaResponse, error) {
	meta := s.leService.Meta()
	return &desc.LeMetaResponse{
		Node:        int64(meta.NodeID),
		DataVersion: meta.DataVersion,
	}, nil
}
