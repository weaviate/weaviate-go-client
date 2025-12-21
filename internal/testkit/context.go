package testkit

import (
	"context"
	"time"
)

// NewTickingContext returns a context which expires after N calls to Done.
func NewTickingContext(ticks int) *TickingContext {
	done := make(chan struct{}, 1)
	done <- struct{}{}
	return &TickingContext{
		ticks:   ticks - 1, // ticks are "0-indexed" for convenience
		done:    done,
		notDone: make(<-chan struct{}),
	}
}

// TickingContext models a context which expires after a number of "ticks".
// A tick happens every time its Done method gets called. On the Nth call
// to Done, it returns a buffered channel, which can be read from immediately.
// While the context is not expired, Done returns an empty channel which blocks
// indefinitely.
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
// using a time-based deadline. With TickingContext we can force the execution
// to reach it by saying it should expire "after 2 ticks".
// Similarly, we can unambiguously model a scenario where the context expires
// exactly on the 6 iteration without introducing additional wait.
type TickingContext struct {
	ticks   int             // remaining ticks, "0-indexed" (ticks==1 means there are 2 ticks remaining)
	done    <-chan struct{} // done channel is returned when ticks expire.
	notDone <-chan struct{} // notDone channel is never sent on.
}

var _ context.Context = (*TickingContext)(nil)

// Deadline always returns ok==true, so that from the outside it looks like this
// this context has a deadline. In does, in fact, but it's tick-based and not time-based.
func (t *TickingContext) Deadline() (deadline time.Time, ok bool) { return time.Time{}, true }
func (t *TickingContext) Value(key any) any                       { return nil }

func (t *TickingContext) Done() <-chan struct{} {
	if t.ticks == 0 {
		return t.notDone
	}
	t.ticks--
	return t.done
}

func (t *TickingContext) Err() error {
	if t.ticks == 0 {
		return context.DeadlineExceeded
	}
	return nil
}
