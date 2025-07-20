package storage

import (
	"sync"
	"time"
)

type KeyValueInMemoryStorage struct {
	items sync.Map
}

type Item struct {
	Value      interface{}
	Expiration int64
}

func NewKeyValueInMemoryStorage() *KeyValueInMemoryStorage {
	return &KeyValueInMemoryStorage{
		items: sync.Map{},
	}
}

func (s *KeyValueInMemoryStorage) Set(key string, value interface{}) {
	s.items.Store(key, Item{Value: value})
}

func (s *KeyValueInMemoryStorage) Get(key string) (Item, bool) {
	val, ok := s.items.Load(key)
	if !ok {
		return Item{}, false
	}
	return val.(Item), true
}

func (s *KeyValueInMemoryStorage) SetWithExpiration(key string, value interface{}, expiration time.Duration) {
	expirationTime := time.Now().Add(expiration).UnixNano()
	s.items.Store(key, Item{
		Value:      value,
		Expiration: expirationTime,
	})
}

func (s *KeyValueInMemoryStorage) Delete(key string) {
	s.items.Delete(key)
}
