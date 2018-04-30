package memcache

import (
	"time"
)

type Auditor interface {
	Start(c *cache)
	Stop()
	CollectErrors(max int) []error
	Chans() (chan struct{}, chan error)
}

type cacheAuditor struct {
	delay    time.Duration // Time
	interval time.Duration
	job      func(c *cache) error
	stopChan chan struct{}
	errChan  chan error
}

func (ch *cacheAuditor) Start(c *cache) {
	stop := make(chan struct{}, 2)
	errchan := make(chan error)
	ticker := time.NewTicker(ch.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := ch.job(c)
				if err != nil {
					errchan <- err
				}
			case <-stop:
				ticker.Stop()
				return
			default:
				time.Sleep(ch.delay)
			}
		}
	}()

	ch.errChan = errchan
	ch.stopChan = stop
}

func (ch *cacheAuditor) Stop() {
	ch.stopChan <- struct{}{}
}

func (ch *cacheAuditor) CollectErrors(max int) (errs []error) {
	errs = make([]error, 0, 100)
	func() {
		for {
			select {
			case e := <-ch.errChan:
				errs = append(errs, e)
				if len(errs) == max {
					return
				}
				continue
			case <-ch.stopChan:
				return
			default:
				return
			}
		}
	}()
	return errs
}

func (ch *cacheAuditor) Chans() (chan struct{}, chan error) {
	return ch.stopChan, ch.errChan
}

func lifetimeAuditor(interval time.Duration, delay time.Duration) *cacheAuditor {
	return &cacheAuditor{
		interval: interval,
		delay:    delay,
		job: func(c *cache) error {
			c.mu.Lock()
			items := c.items
			c.mu.Unlock()

			for key, item := range items {

				if item.lifetime != 0 && time.Since(item.createdAt) > item.lifetime {
					c.mu.Lock()
					delete(c.items, key)
					c.mu.Unlock()
				}
			}
			return nil
		},
	}
}
