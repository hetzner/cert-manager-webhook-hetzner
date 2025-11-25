package limiter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type Limiter struct {
	backoffAfter int
	backoffFunc  hcloud.BackoffFunc

	counterMapLock sync.Mutex
	counterMap     map[string]int
}

type Opts struct {
	// Number of attempts after which the backoff function starts.
	BackoffAfter int
	// Returns a sleep duration based on the number of attempts.
	BackoffFunc hcloud.BackoffFunc
}

func New(opts Opts) *Limiter {
	return &Limiter{
		backoffAfter: opts.BackoffAfter,
		backoffFunc:  opts.BackoffFunc,

		counterMapLock: sync.Mutex{},
		counterMap:     make(map[string]int),
	}
}

func (l *Limiter) backoff(operation string) time.Duration {
	l.counterMapLock.Lock()
	defer l.counterMapLock.Unlock()

	count := l.counterMap[operation]

	if count < l.backoffAfter {
		return time.Duration(0)
	}

	// Start at the bottom of the exponential curve.
	return l.backoffFunc(max(count-l.backoffAfter, 0))
}

func (l *Limiter) update(id string, increase bool) {
	l.counterMapLock.Lock()
	defer l.counterMapLock.Unlock()

	if increase {
		l.counterMap[id]++
	} else {
		if l.counterMap[id] > l.backoffAfter {
			l.counterMap[id] = l.backoffAfter - 1
		} else {
			l.counterMap[id]--
		}
	}
}

func (l *Limiter) Operation(id string) *Operation {
	return &Operation{
		limiter: l,
		id:      id,
	}
}

type Operation struct {
	id      string
	limiter *Limiter
}

func (o *Operation) Backoff() time.Duration {
	return o.limiter.backoff(o.id)
}

func (o *Operation) Sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
	}
	return nil
}

func (o *Operation) Limit(ctx context.Context, logger *slog.Logger) error {
	if duration := o.Backoff(); duration > 0 {
		logger.Warn("too many failures, limiting request rate", "duration", duration)

		if err := o.Sleep(ctx, duration); err != nil {
			return err
		}
	}
	return nil
}

func (o *Operation) Increase(increase bool) {
	o.limiter.update(o.id, increase)
}
