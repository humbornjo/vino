package vino

import (
	"context"
	"time"
)

// contextGrace is a custom context wrapper that provides a graceful shutdown
// mechanism. It waits for the original context to be done, then allows a
// grace period before cancellation. The grace period can be interrupted if
// the provided grace function signals completion early.
type contextGrace struct {
	context.Context
	timeout time.Duration
	graceFn func() <-chan struct{}
	graceCh <-chan struct{}
}

// start listens for the original context's cancellation and then begins the
// grace period. If the grace function completes early, the cancellation is
// triggered before the timeout expires.
func (c *contextGrace) start(cancel context.CancelFunc) {
	select {
	case <-c.Context.Done():
	}

	timer := time.NewTimer(c.timeout)
	select {
	case <-timer.C:
	case <-c.graceFn():
	}
	cancel()
}

// Done returns the cancellation channel of the wrapped context.
func (c contextGrace) Done() <-chan struct{} {
	return c.graceCh
}

// WithGraceContext wraps a given context with a grace period before final
// cancellation. The provided grace function is executed and, if completed
// before the timeout, cancels early.
//
// Parameters:
//   - ctx: The parent context to wrap.
//   - timeout: The duration to wait before forcefully cancelling the context.
//   - graceFn: A function that runs during the grace period to perform cleanup
//     tasks.
//
// Returns:
//   - A new context that respects the grace period.
//   - A cancel function to manually trigger cancellation.
func WithGraceContext(ctx context.Context, timeout time.Duration, graceFn func()) (context.Context, context.CancelFunc) {
	cctx, cancel := context.WithCancel(ctx)

	wrappedGraceFn := func() <-chan struct{} {
		go func() {
			graceFn()
			cancel()
		}()
		return cctx.Done()
	}

	contextGrace := &contextGrace{
		Context: ctx,
		timeout: timeout,
		graceCh: cctx.Done(),
		graceFn: wrappedGraceFn,
	}

	go contextGrace.start(cancel)
	return contextGrace, cancel
}
