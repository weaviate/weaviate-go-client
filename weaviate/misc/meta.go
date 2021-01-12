package misc

import (
	"context"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// MetaGetter builder to get meta endpoint
type MetaGetter struct {
	connection *connection.Connection
}

// Do get the meta endpoint
func (mg *MetaGetter) Do(ctx context.Context) (*models.Meta, error) {
	responseData, responseErr := mg.connection.RunREST(ctx, "/meta", http.MethodGet, nil)
	err := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var meta models.Meta
	parseErr := responseData.DecodeBodyIntoTarget(&meta)
	return &meta, parseErr
}
