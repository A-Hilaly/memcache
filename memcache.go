package memcache

import (
	"errors"
	"sync"
	"time"
)

var (
	errKeyDoesntExist     = errors.New("key doesnt exists")
	errNotImplemented     = errors.New("not implemented")
	errKeyAlreadyExist    = errors.New("key already exists")
	errMaxCapacityReached = errors.New("reached max capacity")
)

// MemCache .
type MemCache struct {
	mu             sync.RWMutex  // global m utex
	imu            sync.RWMutex  // increment mutex
	capacity       uint64        // capacity
	numElements    uint64        // number of items
	store          *MemStore     // memory store
	defaultExpTime time.Duration // default live duration
}

// New return a new MemCache
// capacity: max items
// allocCapacity: startup allocation size cant be > capacity
// defaultExpireTime: item liveness duration
func New(capacity uint64, allocCapacity uint64, defaultExpireTime time.Duration) *MemCache {
	return &MemCache{
		capacity:       capacity,
		numElements:    0,
		store:          NewMemStore(allocCapacity),
		defaultExpTime: defaultExpireTime,
	}
}

// Audit memstore
// Will delete all expire items
func (mc *MemCache) Audit() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.store.Audit()
}

// Renew will clear cache and reallocate new space for allocCapacity
func (mc *MemCache) Renew(allocCapacity uint64) {
	mc.mu.Lock()
	mc.mu.Unlock()
	mc.store = NewMemStore(allocCapacity)
}

// size return numElements
func (mc *MemCache) size() uint64 {
	mc.imu.RLock()
	defer mc.imu.RUnlock()
	return mc.numElements
}

// increment memcache num items
func (mc *MemCache) incr() {
	mc.imu.Lock()
	defer mc.imu.Unlock()
	mc.numElements++
}

// decrement memcache num items
func (mc *MemCache) decr() {
	mc.imu.Lock()
	defer mc.imu.Unlock()
	mc.numElements--
}

// Put will add a new item to the cache store
// It will store item if doesn't exist and max
// capacity isn't reached elsewill return err
// errKeyAlreadyExist
func (mc *MemCache) Put(key string, value interface{}) error {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// if key already exist return err
	if _, exist := mc.store.Get(key); exist {
		return errKeyAlreadyExist
	}

	// check capacity
	if mc.size() >= mc.capacity {
		return errMaxCapacityReached
	}

	mc.incr()

	// compute expiration time and put new item
	tExp := time.Now().Add(mc.defaultExpTime)
	mc.store.Put(key, newItem(value, tExp))

	return nil
}

// Get will return element of the given key and a nil
// error if key already exists
func (mc *MemCache) Get(key string) (interface{}, error) {
	item, exist := mc.store.Get(key)
	if !exist {
		return nil, errKeyDoesntExist
	}

	return item.GetValue(), nil
}

// Delete will try to delete key if it exists else will
// return a err
func (mc *MemCache) Delete(key string) error {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	if _, exist := mc.store.Get(key); !exist {
		return errKeyDoesntExist
	}

	mc.store.Delete(key)
	mc.decr()
	return nil
}

// Update will try to update item if it exists
func (mc *MemCache) Update(key string, value interface{}) error {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// exist
	item, exist := mc.store.Get(key)
	if !exist {
		return errKeyDoesntExist
	}

	item.SetValue(value)
	return nil
}

// Patch will update item if it exist else will create item
// with given value
func (mc *MemCache) Patch(key string, value interface{}) error {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	item, exist := mc.store.Get(key)
	if !exist {
		if mc.size() >= mc.capacity {
			return errMaxCapacityReached
		}

		mc.incr()
		tExp := time.Now().Add(mc.defaultExpTime)
		mc.store.Put(key, newItem(value, tExp))
		return nil
	}

	item.SetValue(value)
	return nil
}
