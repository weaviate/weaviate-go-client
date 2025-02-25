package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// ShardsGetter builder object to get a class' shards
type ShardsGetter struct {
	connection *connection.Connection
	className  string
}

// WithClassName specifies the class to which the shards belong
func (s *ShardsGetter) WithClassName(className string) *ShardsGetter {
	s.className = className
	return s
}

// Do get the status of the shards of the class specified in ShardsGetter
func (s *ShardsGetter) Do(ctx context.Context) ([]*models.ShardStatusGetResponse, error) {
	return getShards(ctx, s.connection, s.className)
}

func getShards(ctx context.Context, conn *connection.Connection, className string) ([]*models.ShardStatusGetResponse, error) {
	responseData, err := conn.RunREST(ctx, fmt.Sprintf("/schema/%s/shards", className), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var shards []*models.ShardStatusGetResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&shards)
		return shards, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
