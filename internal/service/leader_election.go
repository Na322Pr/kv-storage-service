package service

import (
	"github.com/Na322Pr/kv-storage-service/pkg/nodemodel"
	"go.uber.org/zap"
)

type Meta struct {
	NodeID      int
	DataVersion int64
}

type LeService struct {
	node           *nodemodel.Node
	storageService *StorageService

	logger *zap.Logger
}

func NewLeService(
	node *nodemodel.Node,
	storageService *StorageService,
	logger *zap.Logger,
) *LeService {
	return &LeService{
		node:           node,
		storageService: storageService,
		logger:         logger,
	}
}

func (s *LeService) Meta() *Meta {
	return &Meta{
		NodeID:      s.node.ID,
		DataVersion: s.storageService.GetDataVersion(nil),
	}
}

func (s *LeService) SetLeader(leaderID int) {
	if leaderID == s.node.ID {
		s.node.BecomeLeader()
		s.logger.Sugar().Infof("Leader %d is become the leader", s.node.ID)
	} else {
		s.node.BecomeReplica()
		s.logger.Sugar().Infof("Leader %d is become the replica", s.node.ID)
	}
}
