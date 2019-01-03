package memcache

import (
	"sync"
	"time"
)

// MemStore .
type MemStore struct {
	// Read write mutex locker
	mu sync.RWMutex

	// items map
	items map[string]*Item
}

// NewMemStore return a new memory store
// it is safe by multi threaded programs
func NewMemStore(initCapacity uint64) *MemStore {
	c := &MemStore{
		items: make(map[string]*Item, initCapacity),
	}
	return c
}

// Put key: *Item
func (s *MemStore) Put(key string, item *Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = item
}

// Get item with key
func (s *MemStore) Get(key string) (*Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, exist := s.items[key]
	return item, exist
}

// Delete an item with key
func (s *MemStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

func (s *MemStore) Audit() {
	s.mu.RLock()

	// collect to delete keys
	delKeys := make([]string, len(s.items))
	for key, item := range s.items {
		if time.Now().After(item.GetExpireTime()) {
			delKeys = append(delKeys, key)
		}
	}

	// unlock read only mutex
	s.mu.RUnlock()
	if len(delKeys) == 0 {
		return
	}

	// lock write mutex
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range delKeys {
		delete(s.items, key)
	}
}
