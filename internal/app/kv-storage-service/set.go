package kv_storage_service

import (
	"context"
	"fmt"
	"github.com/Na322Pr/kv-storage-service/internal/service"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) Set(ctx context.Context, req *desc.SetRequest) (*desc.SetResponse, error) {
	operation := service.Operation("empty")
	if req.Operation != nil {
		operation = service.Operation(*req.Operation)
	}

	msg := service.SetMessage{
		Key:       req.Key,
		Value:     req.Value,
		Operation: operation,
	}

	s.logger.Debug(fmt.Sprintf("Received request: key=%s, value=%s, operation=%s", msg.Key, msg.Value, msg.Operation))

	if err := s.storageService.Set(ctx, msg); err != nil {
		return &desc.SetResponse{}, err
	}

	return &desc.SetResponse{}, nil
}
