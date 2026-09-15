package store

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/util"
)

type Storage interface {
	Set(key string, value util.Redisobj) (string, bool)
	Get(key string) (util.Redisobj, bool)
	Delete(key string) bool
	Exists(key string) bool
	Keys() []string
	GetAddress(key string)(*util.Redisobj,bool)
	AddBlocking(key string,waiter *util.Waiter)()
	GetBlocking(key string)([]*util.Waiter, bool)
	Deleteblocking(key string)()
}


type MemoryStorage struct {
	data map[string]util.Redisobj
	blocked map[string][]*util.Waiter
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]util.Redisobj),
		blocked: make(map[string][]*util.Waiter),
	}
}

func (s *MemoryStorage) Set(key string, value util.Redisobj) (string, bool) {
	s.data[key] = value
	return "+OK\r\n", true
}

func (s *MemoryStorage) Get(key string) (util.Redisobj, bool) {
	obj, exists := s.data[key]
	if !exists {
		return util.Redisobj{}, false
	}
	if obj.Expiry && obj.Expiryat.Before(time.Now()) {
		delete(s.data, key)
		return util.Redisobj{}, false
	}
	return obj, true
}
func (s *MemoryStorage) GetAddress(key string) (*util.Redisobj, bool) {
	obj, exists := s.data[key]
	if !exists {
		return &util.Redisobj{}, false
	}
	if obj.Expiry && obj.Expiryat.Before(time.Now()) {
		delete(s.data, key)
		return &util.Redisobj{}, false
	}
	return &obj, true
}

func (s *MemoryStorage) Delete(key string) bool {
	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
		return true
	}
	return false
}

func (s *MemoryStorage) Exists(key string) bool {
	_, exists := s.data[key]
	return exists
}

func (s *MemoryStorage) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

func (s *MemoryStorage) AddBlocking(key string, waiter *util.Waiter) {
    s.blocked[key] = append(s.blocked[key], waiter)
}
func (s *MemoryStorage) GetBlocking(key string) ([]*util.Waiter, bool) {
    waiters, exists := s.blocked[key]
    return waiters, exists
}

func (s *MemoryStorage)Deleteblocking(key string){
	_ , exists := s.blocked[key]
	if exists{
		delete(s.blocked,key)
	}
}