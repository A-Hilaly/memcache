package memcache

import (
	"sync"
	"time"
)

func newItem(value interface{}, expireAt time.Time) *Item {
	return &Item{
		expireAt: expireAt,
		value:    value,
	}
}

// Item . Stores expiration time and value
// Is thread safe.
type Item struct {
	mu sync.RWMutex

	expireAt time.Time   // expire at timestamp
	value    interface{} // item value
}

// GetValue .
func (i *Item) GetValue() interface{} {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.value
}

// SetValue .
func (i *Item) SetValue(value interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.value = value
}

// GetExpireTime .
func (i *Item) GetExpireTime() time.Time {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.expireAt
}

// SetExpireTime .
func (i *Item) SetExpireTime(expireAt time.Time) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.expireAt = expireAt
}
