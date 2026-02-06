package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

var contextKey = internal.ContextKey{}

func TestContext(t *testing.T) {
	ctx := internal.ContextWithPlaceholder[int](t.Context(), contextKey)
	require.Nil(t, internal.ValueFromContext[int](ctx, contextKey), "initial value")

	want := testkit.Ptr(92)
	internal.SetContextValue(ctx, contextKey, want)

	got := internal.ValueFromContext[int](ctx, contextKey)
	require.Equal(t, want, got, "retrieved value")
}
