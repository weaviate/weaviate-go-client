package aggregate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type (
	NearVector     Request[query.NearVector]
	NearVectorFunc func(context.Context, NearVector) (*Result, error)
)

// nearVectorFunc encloses transport and request defaults into NearVectorFunc scope.
func nearVectorFunc(t internal.Transport, rd api.RequestDefaults) NearVectorFunc {
	return func(ctx context.Context, nv NearVector) (*Result, error) {
		return aggregate(ctx, t, rd, (*Request[query.NearVector])(&nv), nv.Query.Search(), "near vector")
	}
}
