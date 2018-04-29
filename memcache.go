package memcache

import (
	"sync"
	"time"
)

type Cache interface {
	Put(i Item) error
	Get(key string) (Item, error)

	Del(key string) error
	Update(key string, f func(i Item)) error
	Patch(i Item)

	List() []Item
	Filter(func(key string, i Item) bool) []Item

	ListKeys() []string
	ListValues() []interface{}
	ListFlags() []uint16
	DoFor(func(i Item) error) error

	ExtendLifetime(key string, dur time.Duration)
	SetImmortal(key string)
}

type Item struct {
	sync.RWMutex

	CreatedAt time.Time
	Key       string
	Value     interface{}
	Flag      uint16
	Lifetime  time.Duration
}

func (i *Item) lifetimeExpired() bool {
	return false
}

type store struct {
	sync.Mutex

	open     bool
	id       string
	capacity int64
	items    []Item
}

type StoreManager interface {
	NewStore()

	GetStore()
	SetStore()
	DropStore()
	ListStores()

	OpenStore()
	CloseStore()
}

func NewMemcaceStore() {

}

type MemcacheStores struct {
	stores map[string]Cache
}
