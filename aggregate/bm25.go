package aggregate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	BM25     Request[query.BM25]
	BM25Func func(context.Context, BM25) (*Result, error)
)

// bm25Func encloses transport and request defaults into BM25Func scope.
func bm25Func(t internal.Transport, rd api.RequestDefaults) BM25Func {
	return func(ctx context.Context, bm25 BM25) (*Result, error) {
		return aggregate(ctx, t, rd, (*Request[query.BM25])(&bm25), bm25.Query.Search(), "bm25")
	}
}

// GroupBy runs bm25 aggregation with a GroupBy clause.
func (bf BM25Func) GroupBy(ctx context.Context, bm25 BM25, groupBy GroupBy) (*GroupByResult, error) {
	bm25.groupBy = &groupBy
	return aggregateGroupBy(ctx, bf, bm25)
}
