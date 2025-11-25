package limiter

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestLimiterBackoff(t *testing.T) {
	l := New(Opts{
		BackoffAfter: 1,
		BackoffFunc: hcloud.ExponentialBackoffWithOpts(hcloud.ExponentialBackoffOpts{
			Base:       time.Second,
			Multiplier: 2,
			Cap:        25 * time.Second,
		}),
	})

	assert.Equal(t, time.Duration(0), l.backoff("test"))
	l.update("test", true)
	assert.Equal(t, 1*time.Second, l.backoff("test"))
	l.update("test", true)
	assert.Equal(t, 2*time.Second, l.backoff("test"))

	l.update("test", false)
	assert.Equal(t, time.Duration(0), l.backoff("test"))
	l.update("test", true)
	assert.Equal(t, 1*time.Second, l.backoff("test"))

	assert.Equal(t, time.Duration(0), l.backoff("unknown"))
}

func TestLimiterDo(t *testing.T) {
	l := New(Opts{
		BackoffAfter: 1,
		BackoffFunc: hcloud.ExponentialBackoffWithOpts(hcloud.ExponentialBackoffOpts{
			Base:       time.Second,
			Multiplier: 2,
			Cap:        25 * time.Second,
		}),
	})

	ctx := context.Background()

	assert.Equal(t, 0, l.counterMap["test"])

	{
		op := l.Operation("test")

		duration := op.Backoff()
		assert.Equal(t, time.Duration(0), duration)

		err := op.Limit(ctx, slog.New(slog.DiscardHandler))
		assert.NoError(t, err)

		op.Increase(true)
	}

	assert.Equal(t, 1, l.counterMap["test"])

	{
		op := l.Operation("test")

		duration := op.Backoff()
		assert.Equal(t, 1*time.Second, duration)

		// Skip sleep

		op.Increase(true)
	}

	assert.Equal(t, 2, l.counterMap["test"])

	{
		op := l.Operation("test")

		duration := op.Backoff()
		assert.Equal(t, 2*time.Second, duration)

		// With cancelled context
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		<-ctx.Done()

		err := op.Limit(ctx, slog.New(slog.DiscardHandler))
		assert.EqualError(t, err, "context canceled")

		// No failure => reset to 0
		op.Increase(false)
		assert.Equal(t, 0, l.counterMap["test"])

		// Min count is 0
		op.Increase(false)
		assert.Equal(t, 0, l.counterMap["test"])
	}
}
