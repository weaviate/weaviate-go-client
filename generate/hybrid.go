package generate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	Hybrid     Request[*query.Hybrid]
	HybridFunc func(context.Context, Hybrid) (*Result, error)
)

// hybridFunc encloses transport and request defaults into HybridFunc scope.
func hybridFunc(t internal.Transport, rd api.RequestDefaults) HybridFunc {
	return func(ctx context.Context, h Hybrid) (*Result, error) {
		return generate(ctx, t, rd, (*Request[*query.Hybrid])(&h))
	}
}

// GroupBy runs hybrid generation with a GroupBy clause.
func (hf HybridFunc) GroupBy(ctx context.Context, h Hybrid, groupBy query.GroupBy) (*GroupByResult, error) {
	h.groupBy = &groupBy
	return generateGroupBy(ctx, hf, h)
}
