package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	value, ok := s.storageService.Get(ctx, req.Key)

	resp := &desc.GetResponse{
		Value: value,
		Found: ok,
	}

	return resp, nil
}
