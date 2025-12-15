package backup

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
	// To be awaited, backups need to have a valid [Info.operation] field
	// and a non-nil *Client in [Info.c]. Client correctly populates these
	// fields for Info returned from Create,  Restore,  GetCreateStatus, and GetRestoreStatus.
	errBackupNotAwaitable = errors.New("only backups returned from Create / Restore and get-status are awaitable")

	// Backup completed without reaching desired state. For example, a backup might
	// have been canceled while awaiting StatusSuccess, in which case it will never
	// reach state. The next status check should recognize this and return errBackupStatusFallthrough.
	errBackupStatusFallthrough = errors.New("backup completed without reaching desired status")
)

// AwaitOption controls backup status polling.
type AwaitOption func(*pollingInterval)

// WithPollingInterval sets custom polling interval for checking backup status.
// Use [context.WithDeadline] to set await deadline.
func WithPollingInterval(d time.Duration) AwaitOption {
	return func(pi *pollingInterval) { *pi = pollingInterval(d) }
}

// AwaitCompletion is an AwaitStatus wrapper that awaits [StatusSucess].
func AwaitCompletion(ctx context.Context, backup *Info, options ...AwaitOption) (*Info, error) {
	return AwaitStatus(ctx, backup, StatusSuccess, options...)
}

// AwaitStatus blocks until backup reaches the desired state or otherwise completes.
//
// By default, AwaitStatus will poll backup status once per second and time out after 1 hour.
// The inverval can be adjusted via [WithPollingInterval]. Use [context.WithDeadline] to set a different deadline.
//
// AwaitStatus SHOULD only be called with [Info] obtained from either Create / Restore,
// or GetCreateStatus / GetRestoreStatus, as these will correcly populate the struct's private fields.
//
// Example:
//
//	// GOOD:
//	bak, _ := c.Backup.Create(ctx, "bak-1", "filesystem")
//	backup.AwaitStatus(ctx, bak, backup.StatusCanceled)
//
//	// ALSO GOOD:
//	bak, _ := c.Backup.GetCreateStatus(ctx, "bak-1", "filesystem")
//	backup.AwaitStatus(ctx, bak, backup.StatusCanceled)
//
//	// BAD:
//	backup.AwaitStatus(ctx, &backup.Info{ID: "bak-1"}, backup.StatusCanceled)
//
// AwaitStatus returns immediately if the *Info is nil, the backup's completed,
// or reached the desired status. In the first 2 cases the returned error is not nil.
func AwaitStatus(ctx context.Context, backup *Info, want Status, options ...AwaitOption) (*Info, error) {
	if backup == nil {
		return nil, fmt.Errorf("nil backup")
	} else if backup.Status == want {
		return backup, nil
	} else if backup.IsCompleted() {
		return nil, fmt.Errorf("await %s: %w", want, errBackupStatusFallthrough)
	}

	c := backup.c
	if c == nil {
		return nil, errBackupNotAwaitable
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

	latest, err := c.getStatus(ctx, backup.ID, backup.Backend, backup.operation)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			latest, err = c.getStatus(ctx, backup.ID, backup.Backend, backup.operation)
			if err != nil {
				return nil, err
			}

			dev.Assert(latest != nil, "getStatus returned nil backup.Info")

			if latest.Status == want {
				return latest, nil
			} else if latest.IsCompleted() {
				return latest, fmt.Errorf("await %s: %w", want, errBackupStatusFallthrough)
			}

			// Sleep util the next poll interval, respect context.
			select {
			case <-time.After(time.Duration(interval)):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}
