package kv_storage_service

import (
	"context"
	"fmt"
	"github.com/Na322Pr/kv-storage-service/internal/service"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) SetStream(stream desc.KeyValueStorage_SetStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		operation := service.Operation("empty")
		if req.Operation != nil {
			operation = service.Operation(*req.Operation)
		}

		msg := service.SetMessage{
			Key:       req.Key,
			Value:     req.Value,
			Operation: operation,
		}
		s.logger.Debug(fmt.Sprintf("Received stream request: key=%s, value=%s, operation=%s", msg.Key, msg.Value, msg.Operation))

		if err := s.storageService.Set(context.Background(), msg); err != nil {
			return err
		}

		if err := stream.Send(&desc.SetResponse{}); err != nil {
			return err
		}
	}
}
