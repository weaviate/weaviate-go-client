package internal

import (
	"context"
)

// Custom context key type for internal application.
//
// Empty struct avoids heap allocation when passing it to
// any / interface{} in context.WithValue().
//
// Creating the key at call-site is a mistake -- no 2 instances of ContextKey
// are equal. Declare a shared package-level context key per value type and use
// it across the package.
//
//	package other
//
//	var numberKey internal.ContextKey
//	internal.ContextWithPlaceholder[int](ctx, numberKey)
//
//	var textKey internal.ContextKey
//	internal.ContextWithPlaceholder[string](ctx, textKey)
type ContextKey struct{}

// ContextWithPlaceholder derives a new context, preserving the deadlines
// and cancelation behaviour of the original context.
//
// It creates a nil-pointer to T in the context.Values. That pointer acts
// as a placeholder, which can later be replaced using SetContextValue.
func ContextWithPlaceholder[T any](ctx context.Context, key ContextKey) context.Context {
	placeholder := (*T)(nil)
	return context.WithValue(ctx, key, &placeholder)
}

// Extract value T from the context.
func ValueFromContext[T any](ctx context.Context, key ContextKey) *T {
	v := ctx.Value(key).(**T)
	if v == nil {
		return nil
	}
	return *v
}

// SetContextValue updates value placeholder created by ContextWithPlaceholder.
//
// We want to update the context passed to us in the request,
// rather than derive a new one. In the latter case the original
// context will stay unchanged and the caller will not see the value.
func SetContextValue[T any](ctx context.Context, key ContextKey, v *T) {
	value := ctx.Value(key).(**T)
	*value = v
}
