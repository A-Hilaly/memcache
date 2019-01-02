package memcache

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestMemCache_Renew(t *testing.T) {
	type args struct {
		key           string
		allocCapacity uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "key exist",
			args: args{
				key:           "key one",
				allocCapacity: 10,
			},
		},
		{
			name: "key doesn't exist",
			args: args{
				key:           "unknown key",
				allocCapacity: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := New(100, 100, 2*time.Second)
			mc.Put("key one", 0)

			mc.Renew(tt.args.allocCapacity)
			if _, err := mc.Get(tt.args.key); err == nil {
				t.Errorf("MemCache.Renew(%v) error = %v, wantErr = %v", tt.args.key, err, true)
			}
		})
	}
}

func TestMemCache_Put(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "existing key",
			args: args{
				key:   "key one",
				value: "value one",
			},
			wantErr: true,
		},
		{
			name: "non existing key",
			args: args{
				key:   "key two",
				value: "value one",
			},
			wantErr: false,
		},
		{
			name: "non existing key, on max capacity reached",
			args: args{
				key:   "key three",
				value: "value one",
			},
			wantErr: true,
		},
	}
	mc := New(2, 2, 2*time.Second)
	mc.Put("key one", 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := mc.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemCache.Put(%v, %v) error = %v, wantErr %v", tt.args.key, tt.args.value, err, tt.wantErr)
			}
		})
	}
}

func TestMemCache_Get(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "get value of existing key",
			args: args{
				key: "key one",
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "get value of non existing key",
			args: args{
				key: "key two",
			},
			want:    nil,
			wantErr: true,
		},
	}

	mc := New(2, 2, 2*time.Second)
	mc.Put("key one", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := mc.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemCache.Get(%v) error = %v, wantErr %v", tt.args.key, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MemCache.Get(%v) = %v, want %v", tt.args.key, got, tt.want)
			}
		})
	}
}

func TestMemCache_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "delete existing key",
			args: args{
				key: "key one",
			},
			wantErr: false,
		},
		{
			name: "delete existing key",
			args: args{
				key: "key two",
			},
			wantErr: true,
		},
	}

	mc := New(2, 2, 2*time.Second)
	mc.Put("key one", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mc.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("MemCache.Delete(%v) error = %v, wantErr %v", tt.args.key, err, tt.wantErr)
			}
		})
	}
}

func TestMemCache_Update(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update existing key",
			args: args{
				key: "key one",
			},
			wantErr: false,
		},
		{
			name: "update existing key",
			args: args{
				key: "key two",
			},
			wantErr: true,
		},
	}

	mc := New(2, 2, 2*time.Second)
	mc.Put("key one", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mc.Update(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemCache.Update(%v, %v) error = %v, wantErr %v", tt.args.key, tt.args.value, err, tt.wantErr)
			}
		})
	}
}

func TestMemCache_Patch(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "patch existing key",
			args: args{
				key:   "key one",
				value: 7,
			},
			wantErr: false,
		},
		{
			name: "patch non existing key",
			args: args{
				key:   "key two",
				value: 8,
			},
			wantErr: false,
		},
		{
			name: "patch non existing key",
			args: args{
				key:   "key three",
				value: 9,
			},
			wantErr: true,
		},
	}

	mc := New(2, 2, 2*time.Second)
	mc.Put("key one", 3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := mc.Patch(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MemCache.Patch(%v, %v) error = %v, wantErr %v", tt.args.key, tt.args.value, err, tt.wantErr)
			}
		})
	}
}

func BenchmarkMemCachePutPlat(b *testing.B) {
	mc := New(uint64(b.N), 0, 0*time.Second)
	keys := []string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Put(keys[i], nil)
	}
}

func BenchmarkMemCachePutPreAlloc(b *testing.B) {
	mc := New(uint64(b.N), uint64(b.N), 0*time.Second)
	keys := []string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Put(keys[i], nil)
	}
}

func BenchmarkMemCacheGetExist(b *testing.B) {
	mc := New(1, 1, 0*time.Second)
	mc.Put("key A", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Get("key A")
	}
}

func BenchmarkMemCacheGetNonExist(b *testing.B) {
	mc := New(uint64(b.N), 0, 0*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Get("key X")
	}
}

func BenchmarkMemCacheDeleteExist(b *testing.B) {
	mc := New(uint64(b.N), uint64(b.N), 0*time.Second)
	keys := []string{}
	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		keys = append(keys, key)
		mc.Put(key, nil)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Delete(keys[i])
	}
}

func BenchmarkMemCacheDeleteNonExist(b *testing.B) {
	mc := New(uint64(b.N), 0, 0*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mc.Delete("key A")
	}
}
