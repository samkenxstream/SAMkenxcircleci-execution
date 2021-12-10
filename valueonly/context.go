package valueonly

import (
	"context"
	"time"
)

// Context that wraps another and suppresses its deadline or cancellation.
type valueOnlyContext struct{ context.Context }

func WithNewDeadline(parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(&valueOnlyContext{parent}, deadline)
}

func WithNewTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(&valueOnlyContext{parent}, timeout)
}

func (valueOnlyContext) Deadline() (deadline time.Time, ok bool) { return }
func (valueOnlyContext) Done() <-chan struct{}                   { return nil }
func (valueOnlyContext) Err() error                              { return nil }
