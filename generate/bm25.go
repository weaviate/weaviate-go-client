package generate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	BM25     Request[*query.BM25]
	BM25Func func(context.Context, BM25) (*Result, error)
)

// bm25Func encloses transport and request defaults into BM25Func scope.
func bm25Func(t internal.Transport, rd api.RequestDefaults) BM25Func {
	return func(ctx context.Context, bm25 BM25) (*Result, error) {
		return generate(ctx, t, rd, (*Request[*query.BM25])(&bm25))
	}
}

// GroupBy runs BM25 generation with a GroupBy clause.
func (bf BM25Func) GroupBy(ctx context.Context, bm25 BM25, groupBy query.GroupBy) (*GroupByResult, error) {
	bm25.groupBy = &groupBy
	return generateGroupBy(ctx, bf, bm25)
}
