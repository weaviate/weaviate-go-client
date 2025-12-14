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

// AwaitOption controls backup status polling.
type AwaitOption func(*pollingInterval)

// WithPollingInterval sets custom polling interval for checking backup status.
// Use [context.WithDeadline] to set await deadline.
func WithPollingInterval(d time.Duration) AwaitOption {
	return func(pi *pollingInterval) { *pi = pollingInterval(d) }
}

// AwaitCompletion is an AwaitStatus wrapper that awaits [StatusSucess].
func AwaitCompletion(ctx context.Context, c *Client, backup Info, options ...AwaitOption) (*Info, error) {
	return AwaitStatus(ctx, c, backup, StatusSuccess, options...)
}

// AwaitStatus blocks until backup reaches the desired state or otherwise completes.
//
// By default, AwaitStatus will poll backup status once per second and time out after 1 hour.
// The inverval can be adjusted via [WithPollingInterval]. Use [context.WithDeadline] to set a different deadline.
func AwaitStatus(ctx context.Context, c *Client, backup Info, want Status, options ...AwaitOption) (*Info, error) {
	if backup.Error != "" {
		return nil, errors.New(backup.Error)
	}

	if backup.Status == want {
		return &backup, nil
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
				return latest, fmt.Errorf("backup completed without reaching status %s", want)
			}

			// Sleep util the next poll interval. Respect context.
			select {
			case <-time.After(time.Duration(interval)):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}
