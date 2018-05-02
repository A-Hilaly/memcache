package memcache

import (
	"errors"
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

func Test_cacheAuditor_StartJob(t *testing.T) {
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
		fields   Auditor
		args     args
		wantErrs int // n error
	}{
		{
			name:     "test default",
			fields:   auditDefault,
			args:     args{10},
			wantErrs: 0,
		},
		{
			name:     "test light",
			fields:   auditThirsty,
			args:     args{4},
			wantErrs: 4,
		},
		{
			name:     "test thirsty",
			fields:   auditLight,
			args:     args{10},
			wantErrs: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.Start(nil)
			time.Sleep(100 * milliSecond)
			tt.fields.Stop()
			if gotErrs := tt.fields.CollectErrors(tt.args.max); len(gotErrs) != tt.wantErrs {
				t.Errorf("cacheAuditor.CollectErrors() = %v, want %v", len(gotErrs), tt.wantErrs)
			}
		})
	}
}
