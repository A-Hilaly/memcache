package memcache

import (
	"errors"
	"fmt"
	"reflect"
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
		delay:    5 * milliSecond,
		interval: 15 * milliSecond,
		job:      job,
	}
}

func debugAuditorThirsty(job func(c *cache) error) Auditor {
	return &cacheAuditor{
		delay:    7500 * microSecond,
		interval: 22500 * microSecond,
		job:      job,
	}
}

func debugAuditorLight(job func(c *cache) error) Auditor {
	return &cacheAuditor{
		delay:    10 * milliSecond,
		interval: 30 * milliSecond,
		job:      job,
	}
}

func Test_cacheAuditor_Start(t *testing.T) {
	incr := 0

	auditDefault := debugAuditorDefault(func(c *cache) error {
		incr++
		return nil
	})
	auditThirsty := debugAuditorThirsty(func(c *cache) error {
		incr++
		return nil
	})
	auditLight := debugAuditorLight(func(c *cache) error {
		incr++
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
			want:   9,
		},
		{
			name:   "test thirsty auditor",
			fields: auditThirsty,
			want:   13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Start(nil)
			time.Sleep(100 * milliSecond)
			tt.fields.(*cacheAuditor).stopChan <- struct{}{}
			fmt.Println(incr)
		})
	}
	if incr != 13 {
		t.Errorf("cacheAuditor.Start() want = %v, have = %v", 13, incr)
	}
}

func Test_cacheAuditor_Stop(t *testing.T) {
	incr := 0

	auditDefault := debugAuditorDefault(func(c *cache) error {
		incr++
		return nil
	})
	auditThirsty := debugAuditorThirsty(func(c *cache) error {
		incr++
		return nil
	})
	auditLight := debugAuditorLight(func(c *cache) error {
		incr++
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
			want:   9,
		},
		{
			name:   "test thirsty auditor",
			fields: auditThirsty,
			want:   13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Start(nil)
			tt.fields.Stop()
			time.Sleep(100 * milliSecond)
			fmt.Println(incr)
		})
	}
	if incr != 0 {
		t.Errorf("cacheAuditor.Stop() want = %v, have = %v", 0, incr)
	}
}

func Test_cacheAuditor_CollectErrors(t *testing.T) {
	auditDefault := debugAuditorDefault(func(c *cache) error {
		return nil
	})
	auditThirsty := debugAuditorThirsty(func(c *cache) error {
		return errors.New("test-error-1")
	})
	auditLight := debugAuditorLight(func(c *cache) error {
		return errors.New("test-error-2")
	})

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
