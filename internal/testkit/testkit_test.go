package testkit_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("heterogenous requests", func(t *testing.T) {
		transport := testkit.NewTransport(t, []testkit.Stub[any, any]{
			{
				Request:  testkit.Ptr[any](testkit.Ptr(5)),
				Response: 92,
			},
			{
				Request:  testkit.Ptr[any](testkit.Ptr("hello, transport")),
				Response: true,
			},
		})

		{
			req := 5
			var dest int
			err := transport.Do(t.Context(), &req, &dest)
			require.NoError(t, err, "mock transport error")
			require.Equal(t, 92, dest, "int dest")
		}
		{
			req := "hello, transport"
			var dest bool
			err := transport.Do(t.Context(), &req, &dest)
			require.NoError(t, err, "mock transport error")
			require.Equal(t, true, dest, "bool dest")
		}
	})
}

func TestWithOnly(t *testing.T) {
	type test struct{ testkit.Only }

	// We disable testkit.WithOnly in CI, so this test will always fail.
	// To isolate it, we unset the variable and re-set it on cleanup.
	noWithOnly := os.Getenv(testkit.EnvNoWithOnly)
	require.NoErrorf(t, os.Unsetenv(testkit.EnvNoWithOnly), "unset %s", testkit.EnvNoWithOnly)
	t.Cleanup(func() { os.Setenv(testkit.EnvNoWithOnly, noWithOnly) }) //nolint:errcheck

	for _, tt := range []struct {
		name  string // Test case name.
		tests []test // Exclusive test cases.
		want  int    // How many tt.tests should actually run.
	}{
		{
			name: "all tests",
			tests: []test{
				{}, {}, {},
			},
			want: 3,
		},
		{
			name: "only 1",
			tests: []test{
				{}, {Only: true}, {},
			},
			want: 1,
		},
		{
			name: "only 2",
			tests: []test{
				{}, {Only: true}, {Only: true},
			},
			want: 2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := testkit.WithOnly(t, tt.tests)
			require.Len(t, got, tt.want, "wrong number of tests")
		})
	}
}

func TestError(t *testing.T) {
	require.NotNil(t, testkit.ErrWhaam, "invalid test: ErrWhaam is nil")

	t.Run("nil value", func(t *testing.T) {
		var e testkit.Error

		t.Run("assert", func(t *testing.T) {
			var ok bool
			require.NotPanics(t, func() {
				ok = e.Assert(t, nil)
			})
			require.True(t, ok, "return value")
		})

		t.Run("require", func(t *testing.T) {
			require.NotPanics(t, func() {
				e.Require(t, nil)
			})
		})
	})

	t.Run("ExpectError", func(t *testing.T) {
		t.Run("assert", func(t *testing.T) {
			var ok bool
			require.NotPanics(t, func() {
				ok = testkit.ExpectError.Assert(t, testkit.ErrWhaam)
			})
			require.True(t, ok, "return value")
		})

		t.Run("require", func(t *testing.T) {
			require.NotPanics(t, func() {
				testkit.ExpectError.Require(t, testkit.ErrWhaam)
			})
		})
	})

	t.Run("ErrorIs", func(t *testing.T) {
		t.Run("assert", func(t *testing.T) {
			var ok bool
			require.NotPanics(t, func() {
				ok = testkit.ErrorIs(testkit.ErrWhaam).Assert(t, testkit.ErrWhaam)
			})
			require.True(t, ok, "return value")
		})
		t.Run("require", func(t *testing.T) {
			require.NotPanics(t, func() {
				testkit.ErrorIs(testkit.ErrWhaam).Require(t, testkit.ErrWhaam)
			})
		})
	})
}

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
		assert.NotPanics(t, func() { ctx.Done() }, "close of closed channel")
	case <-time.After(5 * time.Millisecond):
		// When multiple channels can be read from, select will fire a case
		// at random. To avoid flakiness, we block the second channel for a
		// short while, such that it won't stall the test suite even if our
		// tick logic fails.
		require.FailNow(t, "context not done after 2 ticks")
	}
}
