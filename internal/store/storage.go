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
	AddStreamBlocking(streamkey string, waiter *util.StreamWaiter)
    GetStreamBlocking(streamkey string) ([]*util.StreamWaiter, bool)
    RemoveStreamBlocking(streamkey string, waiter *util.StreamWaiter)
	GetVersion(key string)uint64
	SetVersion(key string , version uint64 )
}


type MemoryStorage struct {
	data map[string]util.Redisobj
	blocked map[string][]*util.Waiter
	streamBlocked map[string][]*util.StreamWaiter
	versions map[string]uint64
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]util.Redisobj),
		blocked: make(map[string][]*util.Waiter),
		streamBlocked: make(map[string][]*util.StreamWaiter),
		versions: make(map[string]uint64),
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

func (s *MemoryStorage)AddStreamBlocking(streamkey  string, waiter *util.StreamWaiter){

	streamobj,exists := s.streamBlocked[streamkey]

	if !exists {
		s.streamBlocked[streamkey]= make([]*util.StreamWaiter, 0)
	}

	streamobj = append(streamobj, waiter)
	s.streamBlocked[streamkey]=streamobj

}

func (s *MemoryStorage) RemoveStreamBlocking(streamKey string,target *util.StreamWaiter,) {
    waiters, exists := s.streamBlocked[streamKey]
    if !exists {
        return
    }

    for i, waiter := range waiters {
        if waiter == target {
            waiters = append(
                waiters[:i],
                waiters[i+1:]...,
            )
            break
        }
    }

    if len(waiters) == 0 {
        delete(s.streamBlocked, streamKey)
    } else {
        s.streamBlocked[streamKey] = waiters
    }
}

func (s *MemoryStorage)GetStreamBlocking(streamkey string )([]*util.StreamWaiter,bool){
	
	 waiters, exists := s.streamBlocked[streamkey]

    if !exists {
        return nil, false
    }

    return waiters, true

}
func (s *MemoryStorage)GetVersion(key string)uint64{
	obj , exists := s.versions[key]

	if exists{
		return obj
	}
	return 0
}


func (s *MemoryStorage)SetVersion(key string,version uint64){
	s.versions[key]=version
}




