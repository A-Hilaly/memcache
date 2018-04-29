package memcache

import (
	"fmt"
	"time"
)

type cacheHandler struct {
	interval time.Duration
	handler  func(c *cache) error
	stopChan chan struct{}
	errChan  chan error
}

func NewLifetimeHandler(interval time.Duration) *cacheHandler {
	return &cacheHandler{
		interval: interval,
		handler: func(c *cache) error {
			c.mu.Lock()
			items := c.items
			c.mu.Unlock()

			for key, item := range items {
				if item.lifetime == 0 {
					continue
				}
				if time.Since(item.createdAt) > item.lifetime {
					c.mu.Lock()
					delete(c.items, key)
					c.mu.Unlock()
					return fmt.Errorf("error lol %v", key)
				}
			}
			return nil
		},
	}
}

func (ch *cacheHandler) work(c *cache, stop chan struct{}, errchan chan error) {
	ticker := time.NewTicker(ch.interval)
	for {
		select {
		case <-ticker.C:
			err := ch.handler(c)
			if err != nil {
				errchan <- err
			}
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

func (ch *cacheHandler) Start(c *cache) {
	stop := make(chan struct{}, 1)
	errChannel := make(chan error)
	go ch.work(c, stop, errChannel)
	ch.errChan = errChannel
	ch.stopChan = stop
}

func (ch *cacheHandler) Stop() {
	ch.stopChan <- struct{}{}
}

func (ch *cacheHandler) HandleErrors(f func(err error)) {
	go func() {
		for {
			select {
			case e := <-ch.errChan:
				f(e)
			case <-ch.stopChan:
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
