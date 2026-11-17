package aggregate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	Hybrid     Request[query.Hybrid]
	HybridFunc func(context.Context, Hybrid) (*Result, error)
)

// hybridFunc encloses transport and request defaults into HybridFunc scope.
func hybridFunc(t internal.Transport, rd api.RequestDefaults) HybridFunc {
	return func(ctx context.Context, h Hybrid) (*Result, error) {
		return aggregate(ctx, t, rd, (*Request[query.Hybrid])(&h), h.Query.Search(), "hybrid")
	}
}

// GroupBy runs hybrid aggregation with a GroupBy clause.
func (hf HybridFunc) GroupBy(ctx context.Context, h Hybrid, groupBy GroupBy) (*GroupByResult, error) {
	h.groupBy = &groupBy
	return aggregateGroupBy(ctx, hf, h)
}
