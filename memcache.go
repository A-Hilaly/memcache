package memcache

import (
	"errors"
	"sync"
	"time"
)

var (
	errKeyAlreadyExist    = errors.New("Key already exists")
	errKeyDoesntExist     = errors.New("Key doesnt exists")
	errNotImplemented     = errors.New("Not implemented")
	errMaxCapacityReached = errors.New("Reached max capacity")
)

// Item is cache store elements
type Item struct {
	// Default
	createdAt time.Time
	lifetime  time.Duration
	// Values
	Value interface{}
	Tags  []uint16
}

type cache struct {
	mu        sync.Mutex
	capacity  uint64
	incr      uint64
	items     map[string]Item
	defaultlt time.Duration
	auditor   Auditor
}

// Ensuring that cache implement the Cache
// interface
var _ Cache = &cache{}

// Cache store interface
type Cache interface {
	// puts a new value and it tag for specified key
	Put(key string, value interface{}, tags ...uint16) error
	// Get key value
	Get(key string) (interface{}, error)
	// Get key item
	GetItem(key string) (Item, error)
	// Update key value or tags
	Update(key string, v interface{}, tags ...uint16) error
	// Patch key if key dosent exist it will create it
	Patch(key string, value interface{}, tags ...uint16) error
	// Delete an item associated with key
	Delete(key string) error
	// List all the items without their keys
	List() []Item
	// return cache store map copy
	Items() map[string]Item
	// Filter some keys
	Filter(func(i Item) bool) []Item
	// Do for each item in store
	ForEach(func(i Item))
	// List store keys
	ListKeys() []string
	// List store values
	ListValues() []interface{}
	// extend item lifetime
	ExtendLifetime(key string, dur time.Duration) error
	// set item lifetime to 0 (cannot be deleted)
	Immortalize(key string) error
	// Clear all the memcache db
	Clear()
	// Close the memecache store
	Close()
}

// New memcache store
func New(capacity uint64, defaultLifetime, interval, delay time.Duration) *cache {
	c := &cache{
		capacity:  capacity,
		defaultlt: defaultLifetime,
		items:     make(map[string]Item, capacity),
		auditor:   lifetimeAuditor(interval, delay),
	}
	c.auditor.Start(c)
	return c
}

// Default memcache store
// item life time : 180 seconds
// auditor interval : 30 seconds
// auditor delay : 15 seconds
func Default() *cache {
	return New(10, 180*time.Second, 30*time.Second, 15*time.Second)
}

// Auditor
func (c *cache) Auditor() Auditor {
	return c.auditor
}

// Close cache store
func (c *cache) Close() {
	c.auditor.Stop()
}

// haveKey return true if key exist
// false if it doesnt
func (c *cache) haveKey(key string) bool {
	c.mu.Lock()
	_, exist := c.items[key]
	c.mu.Unlock()
	return exist
}

// Put a key with a value and list of tags
// return an error if element already exist
// or cache store is Full
func (c *cache) Put(key string, value interface{}, tags ...uint16) error {
	if exist := c.haveKey(key); exist {
		return errKeyAlreadyExist
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.incr >= c.capacity {
		return errMaxCapacityReached
	}
	c.items[key] = Item{
		createdAt: time.Now(),
		Value:     value,
		Tags:      tags,
		lifetime:  c.defaultlt,
	}
	c.incr++
	return nil
}

// Get a key item value if item exists exists,
// will return error if not
func (c *cache) Get(key string) (interface{}, error) {
	if exist := c.haveKey(key); !exist {
		return nil, errKeyDoesntExist
	}

	c.mu.Lock()
	v := c.items[key].Value
	c.mu.Unlock()
	return v, nil
}

// GetItem key with key, will return error
// if item doesnt exist
func (c *cache) GetItem(key string) (Item, error) {
	if exist := c.haveKey(key); !exist {
		return Item{}, errKeyDoesntExist
	}

	c.mu.Lock()
	v := c.items[key]
	c.mu.Unlock()
	return v, nil
}

// Update try to to update an existing element
// within the memcache
func (c *cache) Update(key string, value interface{}, tags ...uint16) error {
	if exist := c.haveKey(key); !exist {
		return errKeyDoesntExist
	}

	c.mu.Lock()
	v := c.items[key]
	if len(tags) == 0 {
		c.items[key] = Item{
			createdAt: v.createdAt,
			Value:     value,
			Tags:      v.Tags,
			lifetime:  v.lifetime,
		}
	} else {
		c.items[key] = Item{
			createdAt: v.createdAt,
			Value:     value,
			Tags:      tags,
			lifetime:  v.lifetime,
		}
	}

	c.mu.Unlock()
	return nil
}

// Patch: see https://en.wikipedia.org/wiki/Patch_verb
// Will try to create the item if he doesnt exist
// Else it will update the existing item
func (c *cache) Patch(key string, value interface{}, tags ...uint16) error {
	if c.haveKey(key) {
		return c.Update(key, value, tags...)
	}
	return c.Put(key, value, tags...)
}

// Delete an item with it key value
// return error if item doesnt exist
func (c *cache) Delete(key string) error {
	if exist := c.haveKey(key); !exist {
		return errKeyDoesntExist
	}
	c.mu.Lock()
	delete(c.items, key)
	c.incr--
	c.mu.Unlock()
	return nil
}

// Clear clear all cache by reinitializing
// items map
func (c *cache) Clear() {
	items := make(map[string]Item, c.capacity)
	c.mu.Lock()
	c.items = items
	c.incr = 0
	c.mu.Unlock()
}

// List eturn a list containing all the items
func (c *cache) List() []Item {
	c.mu.Lock()
	items := c.items
	size := c.incr
	c.mu.Unlock()

	it := make([]Item, 0, size)
	for _, v := range items {
		it = append(it, v)
	}
	return it
}

// Items return the original map of items
// strings -> Item
func (c *cache) Items() map[string]Item {
	c.mu.Lock()
	items := c.items
	c.mu.Unlock()
	return items
}

// Return a list of items that satisfy the
// filter f
func (c *cache) Filter(f func(i Item) bool) []Item {
	c.mu.Lock()
	items := c.items
	size := c.incr
	c.mu.Unlock()

	it := make([]Item, 0, size)
	for _, v := range items {
		if f(v) {
			it = append(it, v)
		}
	}
	return it
}

// ForEach apply for each element f
func (c *cache) ForEach(f func(i Item)) {
	c.mu.Lock()
	items := c.items
	c.mu.Unlock()
	for _, v := range items {
		f(v)
	}
}

// ListValues list only the values of items
// in cache items map
func (c *cache) ListValues() []interface{} {
	c.mu.Lock()
	size := c.incr
	c.mu.Unlock()

	array := make([]interface{}, 0, size)
	var collectValue = func(i Item) {
		array = append(array, i.Value)
	}

	c.ForEach(collectValue)
	return array
}

// ListKeys list only the keys of items
// in cache item map
func (c *cache) ListKeys() []string {
	c.mu.Lock()
	items := c.items
	size := c.incr
	c.mu.Unlock()

	array := make([]string, 0, size)
	for key := range items {
		array = append(array, key)
	}

	return array
}

// ExtendLifetime will add lifetime duration
// to the wanted object
func (c *cache) ExtendLifetime(key string, dur time.Duration) error {
	item, err := c.GetItem(key)
	if err != nil {
		return err
	}

	item.lifetime = time.Duration(item.lifetime + dur)
	c.mu.Lock()
	c.items[key] = item
	c.mu.Unlock()
	return nil
}

// Imorialize set the lifetime key to 0
// Which always ignored by cache auditor
func (c *cache) Immortalize(key string) error {
	item, err := c.GetItem(key)
	if err != nil {
		return err
	}

	item.lifetime = 0
	c.mu.Lock()
	c.items[key] = item
	c.mu.Unlock()
	return nil
}
