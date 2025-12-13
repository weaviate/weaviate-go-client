package internal

import "context"

// Custom context key type for internal application.
//
// Empty struct avoids heap allocation when passing it to
// any / interface{} in context.WithValue().
type contextKey struct{}

// groupByKey is used to pass grouped query results from the transport layer.
var groupByResultKey = contextKey{}

// WithGroupByResult derives a new context, preserving the deadlines
// and cancelation behaviour of the original context.
func ContextWithGroupByResult(ctx context.Context) context.Context {
	placeholder := (*GroupByResult)(nil)
	return context.WithValue(ctx, groupByResultKey, &placeholder)
}

// GroupByResult is a placeholder for transport-layer grouped query response.
type GroupByResult struct {
	Objects []any
	Groups  map[string]any
}

// Extract GroupByResult from a context.
func GroupByResultFromContext(ctx context.Context) *GroupByResult {
	v := ctx.Value(groupByResultKey).(**GroupByResult)
	if v == nil {
		return nil
	}
	return *v
}

// Set GroupByResult in the context to another value.
//
// We want to update the context passed to us in the request,
// rather than derive a new one. In the latter case the original
// context will stay unchanged and the caller will not see the value.
//
// Populating api.GroupByResult is NOT a part of the Transport contract,
// but rather a responsibility of the layer using ContextWithGroupByResult.
func SetGroupByResult(ctx context.Context, r *GroupByResult) {
	value := ctx.Value(groupByResultKey).(**GroupByResult)
	*value = r
}
