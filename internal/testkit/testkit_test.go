package testkit_test

import (
	"context"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestNopTransport(t *testing.T) {
	require.NotNil(t, testkit.NopTransport, "testkit.NopTransport")
	require.NoError(t, testkit.NopTransport.Do(t.Context(), nil, nil), "testkit.NopTransport.Do()")
}

func TestMockTransport(t *testing.T) {
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
