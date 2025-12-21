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
