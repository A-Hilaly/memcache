package memcache

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

const (
	nanoSecond  = 1
	microSecond = 1000 * nanoSecond
	milliSecond = 1000 * microSecond
	second      = 1000 * milliSecond
	minute      = 60 * second
)

func debugAuditorDefault(job func(c *cache) error) Auditor {
	return &cacheAuditor{
		delay:    50 * milliSecond,
		interval: 150 * milliSecond,
		job:      job,
	}
}

func debugAuditorThirsty(job func(c *cache) error) Auditor {
	return &cacheAuditor{
		delay:    75 * milliSecond,
		interval: 225 * milliSecond,
		job:      job,
	}
}

func debugAuditorLight(job func(c *cache) error) Auditor {
	return &cacheAuditor{
		delay:    100 * milliSecond,
		interval: 300 * milliSecond,
		job:      job,
	}
}

func Test_cacheAuditor_Start(t *testing.T) {
	var mu sync.Mutex
	var incr []int = make([]int, 3)
	//incrDefault := 0
	//incrLight := 0
	//incrThirsty := 0

	auditDefault := debugAuditorDefault(func(c *cache) error {
		mu.Lock()
		defer mu.Unlock()
		incr[0] = incr[0] + 1
		return nil
	})
	auditThirsty := debugAuditorThirsty(func(c *cache) error {
		mu.Lock()
		defer mu.Unlock()
		incr[1] = incr[1] + 1
		return nil
	})
	auditLight := debugAuditorLight(func(c *cache) error {
		mu.Lock()
		defer mu.Unlock()
		incr[2] = incr[2] + 1
		return nil
	})

	tests := []struct {
		name   string
		fields Auditor
		want   int
	}{
		{
			name:   "test default auditor",
			fields: auditDefault,
			want:   6,
		},
		{
			name:   "test light auditor",
			fields: auditLight,
			want:   3,
		},
		{
			name:   "test thirsty auditor",
			fields: auditThirsty,
			want:   4,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Start(nil)
			time.Sleep(1 * time.Second)
			tt.fields.Stop()
			fmt.Printf("want %v, got: %v \n", tt.want, incr[i])
		})
	}
}

func Test_cacheAuditor_Stop(t *testing.T) {
	type fields struct {
		delay    time.Duration
		interval time.Duration
		job      func(c *cache) error
		stopChan chan struct{}
		errChan  chan error
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &cacheAuditor{
				delay:    tt.fields.delay,
				interval: tt.fields.interval,
				job:      tt.fields.job,
				stopChan: tt.fields.stopChan,
				errChan:  tt.fields.errChan,
			}
			ch.Stop()
		})
	}
}

func Test_cacheAuditor_CollectErrors(t *testing.T) {
	type fields struct {
		delay    time.Duration
		interval time.Duration
		job      func(c *cache) error
		stopChan chan struct{}
		errChan  chan error
	}
	type args struct {
		max int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErrs []error
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &cacheAuditor{
				delay:    tt.fields.delay,
				interval: tt.fields.interval,
				job:      tt.fields.job,
				stopChan: tt.fields.stopChan,
				errChan:  tt.fields.errChan,
			}
			if gotErrs := ch.CollectErrors(tt.args.max); !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("cacheAuditor.CollectErrors() = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

func Test_cacheAuditor_Chans(t *testing.T) {
	type fields struct {
		delay    time.Duration
		interval time.Duration
		job      func(c *cache) error
		stopChan chan struct{}
		errChan  chan error
	}
	tests := []struct {
		name   string
		fields fields
		want   chan struct{}
		want1  chan error
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &cacheAuditor{
				delay:    tt.fields.delay,
				interval: tt.fields.interval,
				job:      tt.fields.job,
				stopChan: tt.fields.stopChan,
				errChan:  tt.fields.errChan,
			}
			got, got1 := ch.Chans()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cacheAuditor.Chans() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("cacheAuditor.Chans() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
