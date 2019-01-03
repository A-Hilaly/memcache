package memcache

import "time"

type auditor interface {
	Audit()
}

type AuditorFunc func(auditor) (stop <-chan struct{}, done chan<- struct{})

func NewAuditorFunc(checkDelay time.Duration, auditDelay time.Duration) AuditorFunc {
	return func(a auditor) (<-chan struct{}, chan<- struct{}) {
		stop := make(chan struct{}, 1)
		done := make(chan struct{}, 1)
		ticker := time.NewTicker(auditDelay)
		go func() {
			for {
				select {
				case <-stop:
					ticker.Stop()
					done <- struct{}{}
					return
				case <-ticker.C:
					a.Audit()
				default:
					time.Sleep(checkDelay)
				}
			}
		}()
		return done, stop
	}

}
