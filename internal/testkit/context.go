package testkit

import (
	"context"
	"time"
)

// NewTickingContext returns a context which expires after N calls to Done.
func NewTickingContext(ticks int) *tickingContext {
	return &tickingContext{
		ticks: ticks - 1, // ticks are "0-indexed" for convenience
		done:  nil,
	}
}

// tickingContext models a context which expires after a number of "ticks".
// A tick happens every time its Done method gets called. On the Nth call
// to Done, the returned channel is closed. The next call to [Context.Err]
// will return [context.DeadlineExceeded].
//
// Ticking context is useful for forcing scenarios which would otherwise be
// depending on real time. Consider the following loop:
//
//	for {
//		select {
//		case <-ctx.Done():
//			return nil
//		default:
//			select {
//			case <-time.After(5*time.Millisecond):
//			case <-ctx.Done():
//				return errors.New("done while sleeping")
//			}
//		}
//	}
//
// YMMV when trying the context to expire in the second select statement
// using a time-based deadline. With tickingContext we can force the execution
// to reach it by saying it should expire "after 2 ticks".
// Similarly, we can unambiguously model a scenario where the context expires
// exactly on the 6 iteration without introducing additional wait.
type tickingContext struct {
	ticks int           // remaining ticks, "0-indexed" (ticks==1 means there are 2 ticks remaining)
	done  chan struct{} // done channel is closed when ticks expire.
	err   error         // not safe for concurrent access.
}

var _ context.Context = (*tickingContext)(nil)

// Deadline always returns ok==true, so that from the outside it looks like this
// this context has a deadline. In does, in fact, but it's tick-based and not time-based.
func (ctx *tickingContext) Deadline() (deadline time.Time, ok bool) { return time.Time{}, true }
func (ctx *tickingContext) Value(key any) any                       { return nil }

func (ctx *tickingContext) Done() <-chan struct{} {
	if ctx.ticks == 0 {
		ctx.err = context.DeadlineExceeded
		if ctx.done == nil {
			ctx.done = make(chan struct{})
			close(ctx.done)
		}
	} else {
		ctx.ticks--
	}
	return ctx.done
}

func (ctx *tickingContext) Err() error { return ctx.err }
