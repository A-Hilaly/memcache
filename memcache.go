package memcache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyAlreadyExist = errors.New("Key already exists")
	ErrKeyDoesntExist  = errors.New("Key doesnt exists")
	ErrNotImplemented  = errors.New("Not implemented")
)

type Item struct {
	createdAt time.Time
	lifetime  time.Duration
	Value     interface{}
	Tags      []uint16
}

type CacheStore interface {
	Auditor() Auditor
	Put(key string, value interface{}, tags ...uint16) error
	Get(key string) (interface{}, error)
	GetItem(key string) (Item, error)
	Update(key string, v interface{}, tags ...uint16) error
	Patch(key string, value interface{}, tags ...uint16) error
	Delete(key string) error
	Clear()
	List() []Item
	Filter(func(i Item) bool) []Item
	ForEach(func(i Item))
	ListKeys() []string
	ListValues() []interface{}
	ExtendLifetime(key string, dur time.Duration) error
	Immortalize(key string) error
}

type cache struct {
	mu        sync.Mutex
	capacity  uint64
	incr      uint64
	items     map[string]Item
	defaultlt time.Duration
	auditor   Auditor
}

func New(capacity uint64, defaultLifetime, interval, delay time.Duration) CacheStore {
	c := &cache{
		capacity:  capacity,
		defaultlt: defaultLifetime,
		items:     make(map[string]Item, capacity),
		auditor:   lifetimeAuditor(interval, delay),
	}
	c.start()
	return c
}

func Default() CacheStore {
	return New(10, 5*time.Second, 5*time.Second, 500*time.Millisecond)
}

func (c *cache) Auditor() Auditor {
	return c.auditor
}

func (c *cache) start() {
	c.auditor.Start(c)
}

func (c *cache) Stop() {
	c.auditor.Stop()
}

func (c *cache) haveKey(key string) bool {
	c.mu.Lock()
	_, exist := c.items[key]
	c.mu.Unlock()
	return exist
}

func (c *cache) Put(key string, value interface{}, tags ...uint16) error {
	if exist := c.haveKey(key); exist {
		return ErrKeyAlreadyExist
	}

	c.mu.Lock()
	c.items[key] = Item{
		createdAt: time.Now(),
		Value:     value,
		Tags:      tags,
		lifetime:  c.defaultlt,
	}
	c.incr++
	c.mu.Unlock()
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	if exist := c.haveKey(key); !exist {
		return nil, ErrKeyDoesntExist
	}

	c.mu.Lock()
	v := c.items[key].Value
	c.mu.Unlock()
	return v, nil
}

func (c *cache) GetItem(key string) (Item, error) {
	if exist := c.haveKey(key); !exist {
		return Item{}, ErrKeyDoesntExist
	}

	c.mu.Lock()
	v := c.items[key]
	c.mu.Unlock()
	return v, nil
}

func (c *cache) Update(key string, value interface{}, tags ...uint16) error {
	if exist := c.haveKey(key); !exist {
		return ErrKeyDoesntExist
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

func (c *cache) Patch(key string, value interface{}, tags ...uint16) error {
	if c.haveKey(key) {
		return c.Update(key, value, tags...)
	}
	c.mu.Lock()
	c.items[key] = Item{
		createdAt: time.Now(),
		Value:     value,
		Tags:      tags,
	}
	c.mu.Unlock()
	return nil
}

func (c *cache) Delete(key string) error {
	if exist := c.haveKey(key); !exist {
		return ErrKeyDoesntExist
	}
	c.mu.Lock()
	delete(c.items, key)
	c.incr--
	c.mu.Unlock()
	return nil
}

func (c *cache) Clear() {
	items := make(map[string]Item, c.capacity)
	c.mu.Lock()
	c.items = items
	c.incr = 0
	c.mu.Unlock()
}

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

func (c *cache) ForEach(f func(i Item)) {
	c.mu.Lock()
	items := c.items
	c.mu.Unlock()
	for _, v := range items {
		f(v)
	}
}

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
