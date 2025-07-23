package service

import (
	"context"
	"github.com/Na322Pr/kv-storage-service/internal/model"
	"github.com/Na322Pr/kv-storage-service/internal/storage"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

type Operation string

const (
	operationEmpty  Operation = "empty"
	operationSet    Operation = "set"
	operationDelete Operation = "delete"
)

type SetMessage struct {
	Key       string
	Value     string
	Operation Operation
}

type StorageService struct {
	store *storage.KeyValueInMemoryStorage
	node  *model.Node
	cm    *ConnectionManagerService
}

func NewStorageService(
	store *storage.KeyValueInMemoryStorage,
	node *model.Node,
	cm *ConnectionManagerService,
) *StorageService {
	return &StorageService{
		store: store,
		node:  node,
		cm:    cm,
	}
}

func (s *StorageService) Set(_ context.Context, msg SetMessage) error {
	switch msg.Operation {
	case operationSet:
		s.store.Set(msg.Key, msg.Value)
	case operationDelete:
		s.store.Delete(msg.Key)
	}

	if !s.node.IsLeader() {
		return nil
	}

	operationString := string(msg.Operation)
	broadcastMsg := &desc.SetRequest{
		Key:       msg.Key,
		Value:     msg.Value,
		Operation: &operationString,
	}

	s.cm.Broadcast(broadcastMsg)

	return nil
}

func (s *StorageService) Get(_ context.Context, key string) (string, bool) {
	value, ok := s.store.Get(key)
	return string(value.Value), ok
}

func (s *StorageService) GetDataVersion(_ context.Context) int64 {
	return s.store.GetDataVersion()
}
