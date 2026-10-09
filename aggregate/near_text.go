package aggregate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	NearText     Request[query.NearText]
	NearTextFunc func(context.Context, NearText) (*Result, error)
)

// nearTextFunc encloses transport and request defaults into NearTextFunc scope.
func nearTextFunc(t internal.Transport, rd api.RequestDefaults) NearTextFunc {
	return func(ctx context.Context, nt NearText) (*Result, error) {
		return aggregate(ctx, t, rd, (*Request[query.NearText])(&nt), nt.Query.Search(), "near text")
	}
}

// GroupBy runs near text aggregation with a GroupBy clause.
func (ntf NearTextFunc) GroupBy(ctx context.Context, nt NearText, groupBy GroupBy) (*GroupByResult, error) {
	nt.groupBy = &groupBy
	return aggregateGroupBy(ctx, ntf, nt)
}
