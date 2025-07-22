package service

import (
	"context"
	"github.com/Na322Pr/kv-storage-service/internal/storage"
)

type StorageService struct {
	store *storage.KeyValueInMemoryStorage
}

func NewStorageService(store *storage.KeyValueInMemoryStorage) *StorageService {
	return &StorageService{
		store: store,
	}
}

func (s *StorageService) Set(_ context.Context, key, value string) error {
	s.store.Set(key, value)
	return nil
}

func (s *StorageService) Get(_ context.Context, key string) (storage.Item, bool) {
	value, ok := s.store.Get(key)
	return value, ok
}

func (s *StorageService) GetDataVersion(_ context.Context) int64 {
	return s.store.GetDataVersion()
}
