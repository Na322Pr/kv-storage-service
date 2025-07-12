package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) Set(ctx context.Context, req *desc.SetRequest) (*desc.SetResponse, error) {
	return &desc.SetResponse{}, nil
}
