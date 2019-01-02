package memcache

import (
	"strconv"
	"testing"
	"time"
)

func TestMemStore_Put(t *testing.T) {

	type args struct {
		key  string
		item *Item
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "new key 1",
			args: args{
				"key 1",
				&Item{value: 1},
			},
		},
		{
			name: "new key 2",
			args: args{
				"key 2",
				&Item{value: 1},
			},
		},
		{
			name: "existing key",
			args: args{
				"key 2",
				&Item{value: 3},
			},
		},
	}

	ms := NewMemStore(10)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.Put(tt.args.key, tt.args.item)
			if _, ok := ms.Get(tt.args.key); !ok {
				t.Errorf("MemStore.Put(%v) got = %v, want %v", tt.args.key, !ok, ok)
			}
		})
	}
}

func TestMemStore_Get(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "key exist",
			args: args{
				"key one",
			},
			want: true,
		},
		{
			name: "key doesn't exist",
			args: args{
				"another key",
			},
			want: false,
		},
	}
	ms := NewMemStore(10)
	ms.Put("key one", &Item{value: "tag A"})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := ms.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("MemStore.Get(%v) got = %v, want %v", tt.args.key, got, tt.want)
			}
		})
	}
}

func TestMemStore_Delete(t *testing.T) {
	ms := NewMemStore(10)
	ms.Put("key one", &Item{value: "tag A"})
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "new key 1",
			args: args{
				"key one",
			},
		},
		{
			name: "new key 2",
			args: args{
				"key two",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.Delete(tt.args.key)
			if _, ok := ms.Get(tt.args.key); ok {
				t.Errorf("MemStore.Delete(%v) got = %v, want %v", tt.args.key, ok, !ok)
			}
		})
	}
}

func TestMemStore_Audit(t *testing.T) {
	type args struct {
		key                   string
		shouldExistAfterAudit bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "expired key",
			args: args{
				key:                   "key one",
				shouldExistAfterAudit: true,
			},
		},
		{
			name: "non expired key",
			args: args{
				key:                   "key two",
				shouldExistAfterAudit: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStore(10)
			ms.Put("key one", &Item{
				value:    "tag A",
				expireAt: time.Now().Add(time.Hour),
			})
			ms.Audit()
			ms.Put("key one", &Item{
				value:    "tag A",
				expireAt: time.Now(),
			})
			ms.Audit()
			if _, exist := ms.Get(tt.args.key); exist && !tt.args.shouldExistAfterAudit {
				t.Errorf("MemStore.Audit(%v) shouldExist = %v, exist = %v", tt.args.key, tt.args.shouldExistAfterAudit, exist)
			}
		})
	}
}

func BenchmarkMemStorePutPlat(b *testing.B) {
	ms := NewMemStore(uint64(b.N))
	keys := []string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Put(keys[i], nil)
	}
}

func BenchmarkMemStorePutPreAlloc(b *testing.B) {
	ms := NewMemStore(uint64(b.N))
	keys := []string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Put(keys[i], nil)
	}
}

func BenchmarkMemStoreGetExist(b *testing.B) {
	ms := NewMemStore(1)
	ms.Put("key A", &Item{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Get("key A")
	}
}

func BenchmarkMemStoreGetNonExist(b *testing.B) {
	ms := NewMemStore(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Get("key A")
	}
}

func BenchmarkMemStoreDeleteExist(b *testing.B) {
	ms := NewMemStore(uint64(b.N))
	keys := []string{}
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		keys = append(keys, key)
		ms.Put(key, nil)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Delete(keys[i])
	}
}

func BenchmarkMemStoreDeleteNonExist(b *testing.B) {
	ms := NewMemStore(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ms.Delete("key A")
	}
}
