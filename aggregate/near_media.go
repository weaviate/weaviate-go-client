package aggregate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	NearMedia     Request[query.NearMedia]
	NearMediaFunc func(context.Context, NearMedia) (*Result, error)
)

// nearMediaFunc encloses transport and request defaults into NearMediaFunc scope.
func nearMediaFunc(t internal.Transport, rd api.RequestDefaults) NearMediaFunc {
	return func(ctx context.Context, nm NearMedia) (*Result, error) {
		return aggregate(ctx, t, rd, (*Request[query.NearMedia])(&nm), nm.Query.Search(), "near media")
	}
}

// GroupBy runs near media aggregation with a GroupBy clause.
func (nmf NearMediaFunc) GroupBy(ctx context.Context, nm NearMedia, groupBy GroupBy) (*GroupByResult, error) {
	nm.groupBy = &groupBy
	return aggregateGroupBy(ctx, nmf, nm)
}
