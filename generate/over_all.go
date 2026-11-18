package generate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	OverAll     Request[*query.OverAll]
	OverAllFunc func(context.Context, OverAll) (*Result, error)
)

// overAllFunc encloses transport and request defaults into OverAllFunc scope.
func overAllFunc(t internal.Transport, rd api.RequestDefaults) OverAllFunc {
	return func(ctx context.Context, oa OverAll) (*Result, error) {
		return generate(ctx, t, rd, (*Request[*query.OverAll])(&oa))
	}
}

// GroupBy runs over all generation with a GroupBy clause.
func (oaf OverAllFunc) GroupBy(ctx context.Context, oa OverAll, groupBy query.GroupBy) (*GroupByResult, error) {
	oa.groupBy = &groupBy
	return generateGroupBy(ctx, oaf, oa)
}
