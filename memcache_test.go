package memcache

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func debugCache() CacheStore {
	return NewCacheStore(50, 10*time.Second, 500*time.Millisecond, 100*time.Millisecond)
}

func Test_cache_haveKey(t *testing.T) {
	type fields struct {
		mu        sync.Mutex
		capacity  uint64
		incr      uint64
		items     map[string]Item
		defaultlt time.Duration
		auditor   Auditor
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "case have key",
			fields: fields{
				items: map[string]Item{"test-key": Item{}},
			},
			args: args{
				key: "test-key",
			},
			want: true,
		},
		{
			name: "case key dosent exist",
			args: args{
				key: "test-key",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				mu:        tt.fields.mu,
				capacity:  tt.fields.capacity,
				incr:      tt.fields.incr,
				items:     tt.fields.items,
				defaultlt: tt.fields.defaultlt,
				auditor:   tt.fields.auditor,
			}
			if got := c.haveKey(tt.args.key); got != tt.want {
				t.Errorf("cache.haveKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_PutGet(t *testing.T) {
	c := debugCache()
	type args struct {
		key   string
		value interface{}
		tags  []uint16
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		wantErr bool
	}{
		{
			name:   "put new key",
			fields: c,
			args: args{
				key:   "test-key",
				value: 10,
			},
			wantErr: false,
		},
		{
			name:   "put the same key",
			fields: c,
			args: args{
				key:   "test-key",
				value: 10,
			},
			wantErr: true,
		},
		{
			name: "put another key",
			args: args{
				key:   "test-key-two",
				value: 777,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := c.Put(tt.args.key, tt.args.value, tt.args.tags...); (err != nil) != tt.wantErr {
				t.Errorf("cache.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
			if v, err := c.GetItem(tt.args.key); err != nil || v.Value != tt.args.value || (v.createdAt == time.Time{}) || (v.lifetime == time.Duration(0)) {
				t.Errorf("cache.Put() value = %v, want %v", tt.args.value, v)
			}
		})
	}
}

func Test_cache_Get(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "get existing key",
			fields: c,
			args: args{
				key: "10__one",
			},
			want:    10,
			wantErr: false,
		},
		{
			name:   "get existing key",
			fields: c,
			args: args{
				key: "10__two",
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "get non existing key",
			fields: c,
			args: args{
				key: "10__three",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("cache.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_GetItem(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "get existing item",
			fields: c,
			args: args{
				key: "10__one",
			},
			want:    10,
			wantErr: false,
		},
		{
			name:   "get existing item",
			fields: c,
			args: args{
				key: "10__two",
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "get non existing item",
			fields: c,
			args: args{
				key: "10__three",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetItem(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("cache.GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(got.Value == tt.want) || (got.lifetime == 0 && !tt.wantErr) {
				t.Errorf("cache.GetItem() = %v, want %v", got.Value, tt.want)
			}
		})
	}
}

func Test_cache_Update(t *testing.T) {
	c := debugCache()
	c.Put("10__one", 10, 2, 3)
	c.Put("10__two", true, 2)
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "update existing item",
			fields: c,
			args: args{
				key:   "10__one",
				value: 20,
			},
			want:    20,
			wantErr: false,
		},
		{
			name:   "update existing item",
			fields: c,
			args: args{
				key:   "10__two",
				value: 30,
			},
			want:    30,
			wantErr: false,
		},
		{
			name:   "update non existing item",
			fields: c,
			args: args{
				key:   "10__three",
				value: []byte("hello world"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Update(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("cache.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if v, _ := c.Get(tt.args.key); v != tt.want {
				t.Errorf("cache.Update() value = %v, want %v", v, tt.want)
			}

		})
	}
}

func Test_cache_Patch(t *testing.T) {
	c := debugCache()
	c.Put("test-1", 10, 2, 3)
	c.Put("test-2", true, 2)
	type args struct {
		key   string
		value interface{}
		tags  []uint16
	}
	tests := []struct {
		name   string
		fields CacheStore
		args   args
	}{
		{
			name:   "patch non existing item",
			fields: c,
			args: args{
				key:   "test-1",
				value: 333,
			},
		},
		{
			name:   "patch existing item",
			fields: c,
			args: args{
				key:   "test-2",
				value: 777,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Patch(tt.args.key, tt.args.value, tt.args.tags...)
			if v, err := c.GetItem(tt.args.key); err != nil || v.Value != tt.args.value {
				t.Errorf("cache.Update() value = %v, want %v", v.Value, tt.args.value)
			}
		})
	}
}

func Test_cache_Delete(t *testing.T) {
	c := debugCache()
	c.Put("test-1", 10, 2, 3)
	c.Put("test-2", true, 2)
	c.Put("test-3", true, 2)
	c.Put("test-4", true, 2)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		wantErr bool
	}{
		{
			name:   "delete existing item",
			fields: c,
			args: args{
				key: "test-1",
			},
			wantErr: false,
		},
		{
			name:   "delete existing item",
			fields: c,
			args: args{
				key: "test-2",
			},
			wantErr: false,
		},
		{
			name:   "delete deleted item",
			fields: c,
			args: args{
				key: "test-2",
			},
			wantErr: true,
		},
		{
			name:   "delete deleted item",
			fields: c,
			args: args{
				key: "test-000",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := c.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("cache.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if c.(*cache).haveKey(tt.args.key) {
				t.Errorf("cache.Delete() error = Item not deleted key = %v", tt.args.key)
			}
		})
	}
}

func Test_cache_Clear(t *testing.T) {
	//NOTE: Lesson learned NEVER IGNORE A TEST
	c := debugCache()
	c.Put("test-1", 10, 2, 3)
	c.Put("test-2", true, 2)
	c.Put("test-3", true, 2)
	c.Put("test-4", true, 2)

	tests := []struct {
		name string
	}{
		{
			name: "clear :)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c.Clear()
			if l := len(c.(*cache).items); l != 0 {
				t.Errorf("cache.Clear() cache not cleared, item size = %v", l)
			}
		})
	}
}

func Test_cache_List(t *testing.T) {
	c := debugCache()
	c.Put("test-1", 10, 2, 3)
	c.Put("test-2", true, 2)
	c.Put("test-3", true, 2)
	c.Put("test-4", true, 2)
	size := 4

	tests := []struct {
		name   string
		fields CacheStore
		want   []Item
	}{
		{
			name:   "test list",
			fields: c,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.List(); len(got) != size {
				t.Errorf("cache.List() error size = %v", len(got))
			}
		})
	}
}

func Test_cache_Filter(t *testing.T) {
	c := debugCache()
	c.Put("test-1", 10, 2, 3)
	c.Put("test-2", true, 2)
	c.Put("test-3", true, 2)
	c.Put("test-4", true, 2, 3)

	type args struct {
		f func(i Item) bool
	}
	tests := []struct {
		name   string
		fields CacheStore
		args   args
		want   []Item
	}{
		{
			name:   "one-filter",
			fields: c,
			args: args{
				f: func(i Item) bool {
					return i.Value == true
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Filter(tt.args.f); len(got) != 3 {
				t.Errorf("cache.Filter() = got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_ForEach(t *testing.T) {
	c := debugCache()
	c.Put("test-1", 10, 2, 3)
	c.Put("test-2", true, 2)
	c.Put("test-3", true, 2)
	c.Put("test-4", true, 2, 3)

	incr := 0

	type args struct {
		f func(i Item)
	}
	tests := []struct {
		name   string
		fields CacheStore
		args   args
	}{
		{
			name:   "for-each incr",
			fields: c,
			args: args{
				f: func(i Item) {
					incr++
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.ForEach(tt.args.f)
			if incr != 4 {
				t.Errorf("cache.ForEach() = got %v, want %v", incr, 4)
			}
		})
	}
}

func Test_cache_ListValues(t *testing.T) {
	c := debugCache()
	c2 := debugCache()
	c.Put("test-1", 10, 2, 3)

	tests := []struct {
		name   string
		fields CacheStore
		want   []interface{}
	}{
		{
			name:   "list empty",
			fields: c2,
			want:   []interface{}{},
		},
		{
			name:   "list non empty",
			fields: c,
			want:   []interface{}{10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.fields.ListValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.ListValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_ListKeys(t *testing.T) {
	c := debugCache()
	c2 := debugCache()
	c.Put("test-1", 10, 2, 3)

	tests := []struct {
		name   string
		fields CacheStore
		want   []string
	}{
		{
			name:   "list empty",
			fields: c2,
			want:   []string{},
		},
		{
			name:   "list non empty",
			fields: c,
			want:   []string{"test-1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.fields.ListKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cache.ListKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_ExtendLifetime(t *testing.T) {
	c := debugCache()
	c.Put("test-1", 10, 2, 3)

	type args struct {
		key string
		dur time.Duration
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		wantErr bool
	}{
		{
			name:   "extend with 0",
			fields: c,
			args: args{
				key: "test-1",
				dur: 0,
			},
			wantErr: false,
		},
		{
			name:   "extend with >0",
			fields: c,
			args: args{
				key: "test-1",
				dur: 3 * time.Hour,
			},
			wantErr: false,
		},
		{
			name:   "extend non existing key",
			fields: c,
			args: args{
				key: "test-non-exist",
				dur: 1 * time.Hour,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.ExtendLifetime(tt.args.key, tt.args.dur); (err != nil) != tt.wantErr {
				t.Errorf("cache.ExtendLifetime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cache_Immortalize(t *testing.T) {
	c := debugCache()
	c.Put("test-1", 10, 2, 3)

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  CacheStore
		args    args
		wantErr bool
	}{
		{
			name:   "extend with 0",
			fields: c,
			args: args{
				key: "test-1",
			},
			wantErr: false,
		},
		{
			name:   "extend non existing key",
			fields: c,
			args: args{
				key: "test-non-exist",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := c.Immortalize(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("cache.Immortalize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
