package service

import (
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
	"sync"
)

type ConnectionManagerService struct {
	connections map[string]desc.KeyValueStorage_SetStreamClient
	mu          sync.RWMutex
}

func NewConnectionManagerService() *ConnectionManagerService {
	return &ConnectionManagerService{
		connections: make(map[string]desc.KeyValueStorage_SetStreamClient),
	}
}

func (cm *ConnectionManagerService) AddConnection(id string, stream desc.KeyValueStorage_SetStreamClient) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[id] = stream
}

func (cm *ConnectionManagerService) RemoveConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, id)
}

func (cm *ConnectionManagerService) GetConnection(id string) (desc.KeyValueStorage_SetStreamClient, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	stream, ok := cm.connections[id]
	return stream, ok
}

func (cm *ConnectionManagerService) Broadcast(msg *desc.SetRequest) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, stream := range cm.connections {
		if err := stream.Send(msg); err != nil {
			// Handle error - maybe remove this connection
			//go cm.RemoveConnection(id)
		}
	}
}
