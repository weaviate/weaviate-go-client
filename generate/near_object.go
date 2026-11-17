package generate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	NearObject     Request[*query.NearObject]
	NearObjectFunc func(context.Context, NearObject) (*Result, error)
)

// nearObjectFunc encloses transport and request defaults into NearObjectFunc scope.
func nearObjectFunc(t internal.Transport, rd api.RequestDefaults) NearObjectFunc {
	return func(ctx context.Context, no NearObject) (*Result, error) {
		return generate(ctx, t, rd, (*Request[*query.NearObject])(&no))
	}
}

// GroupBy runs NearObject generation with a GroupBy clause.
func (nof NearObjectFunc) GroupBy(ctx context.Context, no NearObject, groupBy query.GroupBy) (*GroupByResult, error) {
	no.groupBy = &groupBy
	return generateGroupBy(ctx, nof, no)
}
