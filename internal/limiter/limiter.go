package limiter

import (
	"context"
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

func (l *Limiter) update(operation string, increase bool) {
	l.counterMapLock.Lock()
	defer l.counterMapLock.Unlock()

	if increase {
		l.counterMap[operation]++
	} else {
		if l.counterMap[operation] > l.backoffAfter {
			l.counterMap[operation] = l.backoffAfter - 1
		} else {
			l.counterMap[operation]--
		}
	}
}

func (l *Limiter) Do(operation string, fn func(h *Helper)) {
	h := &Helper{
		limiter:   l,
		operation: operation,
		increase:  false,
	}

	fn(h)

	l.update(operation, h.increase)
}

type Helper struct {
	limiter   *Limiter
	operation string
	increase  bool
}

func (h *Helper) Backoff() time.Duration {
	return h.limiter.backoff(h.operation)
}

func (h *Helper) Increase() {
	h.increase = true
}

func (h *Helper) Sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
	}
	return nil
}
