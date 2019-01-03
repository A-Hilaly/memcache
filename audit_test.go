package memcache

import (
	"testing"
	"time"
)

func TestNewAuditorFunc(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name        string
		args        args
		shouldExist bool
	}{
		{
			name: "key should exist after AuditFunc auditing",
			args: args{
				key: "key 1",
			},
			shouldExist: false,
		},
		{
			name: "key shouldn't exist after AuditFunc auditing",
			args: args{
				key: "key 2",
			},
			shouldExist: true,
		},
	}
	mc := New(5, 5, 75*time.Millisecond)
	auditFn := NewAuditorFunc(20*time.Millisecond, 30*time.Millisecond)
	done, stop := auditFn(mc)
	mc.Put("key 1", "tag 1")
	time.Sleep(80 * time.Millisecond)
	mc.Put("key 2", "tag 1")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := mc.Get(tt.args.key); err != nil && tt.shouldExist {
				t.Errorf("MemCache.Renew(%v) error = %v, should exist = %v", tt.args.key, err, tt.shouldExist)
			}
		})
	}
	stop <- struct{}{}
	<-done
}
