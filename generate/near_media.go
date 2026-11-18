package generate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	NearMedia     Request[*query.NearMedia]
	NearMediaFunc func(context.Context, NearMedia) (*Result, error)
)

// nearMediaFunc encloses transport and request defaults into NearMediaFunc scope.
func nearMediaFunc(t internal.Transport, rd api.RequestDefaults) NearMediaFunc {
	return func(ctx context.Context, nm NearMedia) (*Result, error) {
		return generate(ctx, t, rd, (*Request[*query.NearMedia])(&nm))
	}
}

// GroupBy runs NearMedia generation with a GroupBy clause.
func (nmf NearMediaFunc) GroupBy(ctx context.Context, nm NearMedia, groupBy query.GroupBy) (*GroupByResult, error) {
	nm.groupBy = &groupBy
	return generateGroupBy(ctx, nmf, nm)
}
