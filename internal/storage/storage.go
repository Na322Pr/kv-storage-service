package storage

import "sync"

type KeyValueInMemoryStorage struct {
	items map[string]Item
	mu    *sync.RWMutex
}

type Item struct {
	Value      interface{}
	Expiration int64
}

func NewKeyValueInMemoryStorage() *KeyValueInMemoryStorage {
	return &KeyValueInMemoryStorage{
		items: make(map[string]Item),
		mu:    &sync.RWMutex{},
	}
}

func (s *KeyValueInMemoryStorage) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = Item{Value: value}
}

func (s *KeyValueInMemoryStorage) Get(key string) Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items[key]
}
