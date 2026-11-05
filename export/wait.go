package export

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

// pollingInterval is currently the only configurable AwaitOption.
// We keep the type unexported to be able to extend it later.
type pollingInterval time.Duration

const (
	defaultTimeout                         = 1 * time.Hour
	defaultPollingInterval pollingInterval = pollingInterval(1 * time.Second)
)

var (
	// To be awaited, exports need to have a non-nil *Client in [Info.c]. Client correctly populates these
	// fields for Info returned from Create and  GetCreateStatus.
	errExportNotAwaitable = errors.New("only exports returned from Create and GetCreateStatus are awaitable")

	// Export completed without reaching desired state. For example, a export might
	// have been canceled while awaiting StatusSuccess, in which case it will never
	// reach state. The next status check should recognize this and return errExportStatusFallthrough.
	errExportStatusFallthrough = errors.New("export completed without reaching desired status")
)

// AwaitOption controls export status polling.
type AwaitOption func(*pollingInterval)

// WithPollingInterval sets custom polling interval for checking export status.
// Use [context.WithDeadline] to set await deadline.
func WithPollingInterval(d time.Duration) AwaitOption {
	return func(pi *pollingInterval) { *pi = pollingInterval(d) }
}

// AwaitCompletion is an AwaitStatus wrapper that awaits [StatusSucess].
func AwaitCompletion(ctx context.Context, export *Info, options ...AwaitOption) (*Info, error) {
	return AwaitStatus(ctx, export, StatusSuccess, options...)
}

// AwaitStatus blocks until export reaches the desired state or otherwise completes.
//
// By default, AwaitStatus will poll export status once per second and time out after 1 hour.
// The inverval can be adjusted via [WithPollingInterval]. Use [context.WithDeadline] to set a different deadline.
//
// AwaitStatus SHOULD only be called with [Info] obtained from either Create, or GetCreateStatus,
// as these will correcly populate the struct's private fields.
//
// Example:
//
//	// GOOD:
//	exp, _ := c.Export.Create(ctx, "exp-1", "filesystem")
//	export.AwaitStatus(ctx, exp, export.StatusCanceled)
//
//	// ALSO GOOD:
//	exp, _ := c.Export.GetCreateStatus(ctx, "exp-1", "s3")
//	export.AwaitStatus(ctx, exp, export.StatusCanceled)
//
//	// BAD:
//	export.AwaitStatus(ctx, &export.Info{ID: "exp-1"}, export.StatusCanceled)
//
// AwaitStatus returns immediately if the *Info is nil, the export's completed,
// or reached the desired status. In the first 2 cases the returned error is not nil.
func AwaitStatus(ctx context.Context, exp *Info, want Status, options ...AwaitOption) (*Info, error) {
	if exp == nil {
		return nil, fmt.Errorf("nil export")
	} else if exp.Status == want {
		return exp, nil
	} else if exp.IsCompleted() {
		return exp, fmt.Errorf("await %s: %w", want, errExportStatusFallthrough)
	}

	c := exp.c
	if c == nil {
		return nil, errExportNotAwaitable
	}

	interval := defaultPollingInterval
	for _, opt := range options {
		opt(&interval)
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(defaultTimeout))
		defer cancel()
	}

	_, hasDeadline := ctx.Deadline()
	dev.Assert(hasDeadline, "unbounded await context")

	conf := GetStatusOptions{Backend: exp.Backend, ID: exp.ID}
	current := exp
	for {
		select {
		case <-ctx.Done():
			return current, ctx.Err()
		default:
			latest, err := c.GetCreateStatus(ctx, conf)
			if err != nil {
				return current, err
			}
			current = latest

			dev.AssertNotNil(latest, "GetStatus returned nil export.Info")

			if latest.Status == want {
				return latest, nil
			} else if latest.IsCompleted() {
				return latest, fmt.Errorf("await %s: %w", want, errExportStatusFallthrough)
			}

			// Sleep util the next poll interval, respect context.
			select {
			case <-time.After(time.Duration(interval)):
			case <-ctx.Done():
				return latest, ctx.Err()
			}
		}
	}
}
