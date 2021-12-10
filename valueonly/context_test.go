package valueonly

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestContextValues(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testkey", "testvalue")
	derivedCtx, _ := WithNewDeadline(ctx, time.Now())
	assert.Equal(t, derivedCtx.Value("testkey"), "testvalue")
}

func TestNewContextExpiration(t *testing.T) {
	oldDeadline := time.Now().Add(-time.Minute)
	oldCtx, _ := context.WithDeadline(context.Background(), oldDeadline)

	t.Run("deadline", func(t *testing.T) {
		newDeadline := time.Now().Add(time.Minute)
		derivedCtx, _ := WithNewDeadline(oldCtx, newDeadline)

		actualDeadline, _ := derivedCtx.Deadline()
		assert.Equal(t, actualDeadline, newDeadline)
		assert.Check(t, actualDeadline != oldDeadline)
	})

	t.Run("timeout", func(t *testing.T) {
		now := time.Now()
		timeout := time.Second * 100
		derivedCtx, _ := WithNewTimeout(oldCtx, timeout)

		actualDeadline, _ := derivedCtx.Deadline()
		expectedDeadline := now.Add(timeout)
		delta := actualDeadline.Sub(expectedDeadline)
		assert.Check(t, delta < time.Millisecond,
			fmt.Sprintf("real deadline: %v must be within 1ms since the expected deadline: %v, is %v",
				actualDeadline, expectedDeadline, delta))
		assert.Check(t, actualDeadline != oldDeadline)
	})
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	derivedCtx, derivedCancel := WithNewTimeout(ctx, time.Second*10)

	assert.Check(t, !isDone(ctx))
	assert.Check(t, !isDone(derivedCtx))

	cancel()

	assert.Check(t, isDone(ctx))
	assert.Check(t, !isDone(derivedCtx))

	assert.Equal(t, ctx.Err().Error(), "context canceled", "original context cancelled")
	assert.NilError(t, derivedCtx.Err(), "derived context not cancelled")

	derivedCancel()
	assert.Equal(t, derivedCtx.Err().Error(), "context canceled", "derived context cancelled")
}

func isDone(ctx context.Context) bool {
	done := false
	select {
	case <-ctx.Done():
		done = true
	default:
		done = false
	}
	return done
}

//func TestContextCancellation(t *testing.T) {
//	// Creating original context and deadlining it.
//	ctx := context.Background()
//	setDeadline := time.Now().Add(-time.Second)
//	ctx, cancel := context.WithDeadline(ctx, setDeadline)
//	defer cancel()
//
//	done := false
//	select {
//	case <-ctx.Done():
//		done = true
//	default:
//		done = false
//	}
//	actualDeadline, deadlineSet := ctx.Deadline()
//	assert.Equal(t, actualDeadline, setDeadline, "original context's deadline is set")
//	assert.Equal(t, deadlineSet, true, "original context's deadline is set")
//	assert.Equal(t, done, true, "original context's done channel resolved")
//	assert.Equal(t, ctx.Err().Error(), "context deadline exceeded",
//		"original context cancelled by deadline")
//
//	// Checking that the derived context is not cancelled.
//	derivedCtx := &Context{ctx}
//	derivedDone := false
//	select {
//	case <-derivedCtx.Done():
//		derivedDone = true
//	default:
//		derivedDone = false
//	}
//
//	_, derivedDeadlineSet := derivedCtx.Deadline()
//	assert.Equal(t, derivedDeadlineSet, false, "derived context has no deadline set")
//	assert.Equal(t, derivedDone, false, "derived context's done channel not resolved")
//	assert.Equal(t, derivedCtx.Err(), nil, "derived context's error is nil")
//
//	// Checking that we can cancel the derived context.
//	derivedCtxTryCancel, cancel := context.WithCancel(derivedCtx)
//	cancel()
//	assert.Equal(t, derivedCtxTryCancel.Err().Error(), "context canceled",
//		"can cancel a derived context")
//}
