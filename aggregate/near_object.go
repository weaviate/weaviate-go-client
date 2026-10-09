package aggregate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	NearObject     Request[query.NearObject]
	NearObjectFunc func(context.Context, NearObject) (*Result, error)
)

// nearObjectFunc encloses transport and request defaults into NearObjectFunc scope.
func nearObjectFunc(t internal.Transport, rd api.RequestDefaults) NearObjectFunc {
	return func(ctx context.Context, no NearObject) (*Result, error) {
		return aggregate(ctx, t, rd, (*Request[query.NearObject])(&no), no.Query.Search(), "near object")
	}
}

// GroupBy runs near object aggregation with a GroupBy clause.
func (nof NearObjectFunc) GroupBy(ctx context.Context, no NearObject, groupBy GroupBy) (*GroupByResult, error) {
	no.groupBy = &groupBy
	return aggregateGroupBy(ctx, nof, no)
}
