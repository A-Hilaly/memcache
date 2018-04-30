package memcache

import (
	"reflect"
	"testing"
	"time"
)

func Test_cacheAuditor_Start(t *testing.T) {
	type fields struct {
		delay    time.Duration
		interval time.Duration
		job      func(c *cache) error
		stopChan chan struct{}
		errChan  chan error
	}
	type args struct {
		c *cache
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
			ch.Start(tt.args.c)
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
