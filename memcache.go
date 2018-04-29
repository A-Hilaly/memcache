package memcache

import (
	"errors"
	"fmt"
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

type DataStore interface {
	Start()
	Stop()

	Put(key string, value interface{}, tags ...uint16) error
	Get(key string) (Item, error)

	Update(key string, f func(i Item)) error
	Patch(i Item)

	Delete(key string) error
	ClearAll()

	List() []Item
	Filter(func(key string, i Item) bool) []Item
	ForEach(func())

	ListKeys() []string
	ListValues() []interface{}

	ExtendLifetime(key string, dur time.Duration)
	SetImmortal(key string)
}

type cache struct {
	mu        sync.Mutex
	capacity  uint64
	incr      uint64
	items     map[string]Item
	defaultlt time.Duration
	handler   *cacheHandler
}

func New(capacity uint64, defaultLifetime time.Duration, handler *cacheHandler) *cache {
	return &cache{
		capacity:  capacity,
		defaultlt: defaultLifetime,
		items:     make(map[string]Item, capacity),
		handler:   handler,
	}
}

func Default() *cache {
	return New(10, 5*time.Second, NewLifetimeHandler(2))
}

func Debug() *cache {
	return New(10, 500*time.Millisecond, NewLifetimeHandler(400*time.Millisecond))
}

func (c *cache) Start() {
	c.handler.Start(c)
	c.handler.HandleErrors(func(err error) {
		fmt.Println(err.Error())
	})
}

func (c *cache) Stop() {
	c.handler.Stop()
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

func (c *cache) Update(key string, value interface{}) error {
	if exist := c.haveKey(key); !exist {
		return ErrKeyDoesntExist
	}

	c.mu.Lock()
	v := c.items[key]
	c.items[key] = Item{
		createdAt: v.createdAt,
		Value:     value,
		Tags:      v.Tags,
		lifetime:  v.lifetime,
	}
	c.mu.Unlock()
	return nil
}

func (c *cache) Patch(key string, value interface{}, tags ...uint16) {
	c.mu.Lock()
	c.items[key] = Item{
		createdAt: time.Now(),
		Value:     value,
		Tags:      tags,
	}
	c.mu.Unlock()
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

	it := make([]Item, size)
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
