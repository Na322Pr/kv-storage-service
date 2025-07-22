package storage

import (
	"sync"
	"time"
)

type KeyValueInMemoryStorage struct {
	items   sync.Map
	version int64
}

type Item struct {
	Value      interface{}
	Expiration int64
}

func NewKeyValueInMemoryStorage() *KeyValueInMemoryStorage {
	return &KeyValueInMemoryStorage{
		items:   sync.Map{},
		version: 0,
	}
}

func (s *KeyValueInMemoryStorage) Set(key string, value interface{}) {
	s.items.Store(key, Item{Value: value})
	s.version++
}

func (s *KeyValueInMemoryStorage) Get(key string) (Item, bool) {
	val, ok := s.items.Load(key)
	if !ok {
		return Item{}, false
	}
	return val.(Item), true
}

func (s *KeyValueInMemoryStorage) GetDataVersion() int64 {
	return s.version
}

func (s *KeyValueInMemoryStorage) SetWithExpiration(key string, value interface{}, expiration time.Duration) {
	expirationTime := time.Now().Add(expiration).UnixNano()
	s.items.Store(key, Item{
		Value:      value,
		Expiration: expirationTime,
	})
	s.version++
}

func (s *KeyValueInMemoryStorage) Delete(key string) {
	s.items.Delete(key)
	s.version++
}
