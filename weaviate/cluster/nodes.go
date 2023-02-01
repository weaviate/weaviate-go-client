package cluster

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type NodesStatusGetter struct {
	connection *connection.Connection
}

// Do get the nodes endpoint
func (nsg *NodesStatusGetter) Do(ctx context.Context) (*models.NodesStatusResponse, error) {
	responseData, responseErr := nsg.connection.RunREST(ctx, "/nodes", http.MethodGet, nil)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var nodesStatus models.NodesStatusResponse
	parseErr := responseData.DecodeBodyIntoTarget(&nodesStatus)
	return &nodesStatus, parseErr
}
