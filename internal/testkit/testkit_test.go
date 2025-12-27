package testkit_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestTickingContext(t *testing.T) {
	ctx := testkit.NewTickingContext(2)
	require.Implements(t, (*context.Context)(nil), ctx)

	_, ok := ctx.Deadline()
	assert.True(t, ok, "must report that deadline is set")
	assert.NoError(t, ctx.Err(), "context is initially valid")

	select {
	case <-ctx.Done():
		require.FailNow(t, "context expired after 1 tick, want 2")
	default:
	}
	require.NoError(t, ctx.Err(), "context error after 1 tick")

	select {
	case <-ctx.Done():
		assert.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)
	case <-time.After(5 * time.Millisecond):
		// When multiple channels can be read from, select will fire a case
		// at random. To avoid flakiness, we block the second channel for a
		// short while, such that it won't stall the test suite even if our
		// tick logic fails.
		require.FailNow(t, "context not done after 2 ticks")
	}
}

func TestMockTransport(t *testing.T) {
	// Test Do:
	// 	- compares non-nil request
	// 	- writes to non-nil dest
	// Test Done => true when done, but not before.

	t.Run("respects context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		transport := testkit.NewTransport(t, make([]testkit.Stub[any, any], 1))
		err := transport.Do(ctx, nil, nil)
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("returns error from stubs", func(t *testing.T) {
		transport := testkit.NewTransport(t, []testkit.Stub[any, any]{
			{Err: testkit.ErrWhaam},
		})
		err := transport.Do(t.Context(), nil, nil)
		require.ErrorIs(t, err, testkit.ErrWhaam)
	})

	t.Run("writes to non-nil dest", func(t *testing.T) {
		var dest bool
		transport := testkit.NewTransport(t, []testkit.Stub[any, bool]{
			{Response: true},
		})

		err := transport.Do(t.Context(), nil, &dest)
		assert.NoError(t, err)
		assert.Equal(t, true, dest, "dest not updated")
	})

	t.Run("done when all requests consumed", func(t *testing.T) {
		n := 10
		transport := testkit.NewTransport(t, make([]testkit.Stub[any, any], n))

		for range n - 1 {
			err := transport.Do(t.Context(), nil, nil)
			n-- // keep our own tally to as a sanity check
			require.NoError(t, err, "mock transport error")
			require.False(t, transport.Done(), "done: %d requests remaining", n)
		}
		err := transport.Do(t.Context(), nil, nil)
		n--
		require.NoError(t, err, "mock transport error")

		require.Equal(t, n, 0, "mistake in test code")
		require.True(t, transport.Done(), "done: all requests consumed")
	})
}
