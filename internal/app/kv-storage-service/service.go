package kv_storage_service

import (
	"github.com/Na322Pr/kv-storage-service/internal/service"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

type Implementation struct {
	desc.UnimplementedKeyValueStorageServer

	nodeService    *service.NodeService
	storageService *service.StorageService
	leService      *service.LeService
}

func NewImplementation(
	nodeService *service.NodeService,
	storeService *service.StorageService,
	leService *service.LeService,
) *Implementation {
	return &Implementation{
		nodeService:    nodeService,
		storageService: storeService,
	}
}
