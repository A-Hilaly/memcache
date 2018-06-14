package memcache

import (
	"time"
)

// Auditor janitor, doctor is a multi purpose goroutine
// that handles Item expirations and errors
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

// Start the cache auditor
func (ch *cacheAuditor) Start(c *cache) {
	stop := make(chan struct{})
	errchan := make(chan error, 100)
	// Time ticker
	ticker := time.NewTicker(ch.interval)
	// go goroutine
	go func() {
		for {
			select {
			case <-ticker.C:
				err := ch.job(c)
				if err != nil {
					errchan <- err
				}
			case <-ch.stopChan:
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

// Pend stop signal to cacheAuditor
func (ch *cacheAuditor) Stop() {
	ch.stopChan <- struct{}{}
}

// CollectErrors all the errors
func (ch *cacheAuditor) CollectErrors(max int) (errs []error) {
	errs = make([]error, 0, max)
	func() {
		for i := 0; i < max; i++ {
			select {
			case e := <-ch.errChan:
				errs = append(errs, e)
				if len(errs) > max {
					return
				}
				continue
			case <-ch.stopChan:
				break
			default:
				return
			}
		}
	}()
	return errs
}

// Chans return stopChan and errChan
func (ch *cacheAuditor) Chans() (chan struct{}, chan error) {
	return ch.stopChan, ch.errChan
}
