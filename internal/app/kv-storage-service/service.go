package kv_storage_service

import (
	"github.com/Na322Pr/kv-storage-service/internal/service"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
	"go.uber.org/zap"
)

type Implementation struct {
	desc.UnimplementedKeyValueStorageServer

	nodeService    *service.NodeService
	storageService *service.StorageService
	leService      *service.LeService

	logger *zap.Logger
}

func NewImplementation(
	nodeService *service.NodeService,
	storeService *service.StorageService,
	leService *service.LeService,
	logger *zap.Logger,
) *Implementation {
	return &Implementation{
		nodeService:    nodeService,
		storageService: storeService,
		leService:      leService,
		logger:         logger,
	}
}
