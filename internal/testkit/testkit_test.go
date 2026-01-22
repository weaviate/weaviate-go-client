package testkit_test

import (
	"context"
	"os"
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

func TestWithOnly(t *testing.T) {
	type test struct{ testkit.Only }

	// We disable testkit.WithOnly in CI, so this test will always fail.
	// To isolate it, we unset the variable and re-set it on cleanup.
	noWithOnly := os.Getenv(testkit.EnvNoWithOnly)
	require.NoErrorf(t, os.Unsetenv(testkit.EnvNoWithOnly), "unset %s", testkit.EnvNoWithOnly)
	t.Cleanup(func() { os.Setenv(testkit.EnvNoWithOnly, noWithOnly) })

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
