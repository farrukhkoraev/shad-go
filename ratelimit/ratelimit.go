//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
	"sync"
)

// Limiter is precise rate limiter with context support.
type req struct {
	ctx context.Context
	resp chan error 
}

type Limiter struct {
	requests chan req
	stop chan int

	once sync.Once
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxCount events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {

	stop := make(chan int, 1)
	reqCh   := make(chan req, 1)
    timestamps  := []time.Time{}

	l := Limiter {reqCh, stop, sync.Once{}}

	go func() {
		for {
			var waitCh <-chan time.Time
			if len(timestamps) == maxCount {
				// timer for expiry of oldest slot
				waitCh = time.After(timestamps[0].Add(interval).Sub(time.Now()))
			}

			select {
			case <-stop:
				// send ErrStoped to waiting goroutines
				for {
					select {
					case r := <-reqCh:
						r.resp <- ErrStopped
					default:
						return
					}
				}
			case <-waitCh:
				// oldest slot expired, loop again to accept request
			case r := <-reqCh:
				select {
				case <-r.ctx.Done():
					r.resp <- r.ctx.Err()
					continue
				default:
				}

				now := time.Now()
				j := -1
				for i, ts := range timestamps {
					if now.Sub(ts) > interval {
						j = i
					}
				}
				if j >= 0 {
					timestamps = timestamps[j+1:]
				}

				if len(timestamps) == maxCount {
					dt := timestamps[0].Add(interval).Sub(now)
					select {
					case <-time.After(dt):
						now = time.Now()
						timestamps = timestamps[1:]
					case <-r.ctx.Done():
						r.resp <- r.ctx.Err()
						continue
					case <-stop:
						r.resp <- ErrStopped
						return
					}
				}

				timestamps = append(timestamps, now)
				r.resp <- nil
			}
		}
	}()
	
	return &l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	r := req { ctx, make(chan error, 1) }

	select {
	case <-l.stop:
		return ErrStopped 
	case l.requests <- r:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <- r.resp:
		return err
	}

	return nil
}

func (l *Limiter) Stop() {
	l.once.Do(func() {
		l.stop <- 1
		close(l.stop)
	})
}
