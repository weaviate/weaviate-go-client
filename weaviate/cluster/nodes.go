package cluster

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type NodesStatusGetter struct {
	connection *connection.Connection
	class      string
	output     string
}

// Do get the nodes endpoint
func (nsg *NodesStatusGetter) Do(ctx context.Context) (*models.NodesStatusResponse, error) {
	path := "/nodes"
	if nsg.class != "" {
		path += "/" + nsg.class
	}
	if nsg.output != "" {
		path += "?output=" + nsg.output
	}

	responseData, responseErr := nsg.connection.RunREST(ctx, path, http.MethodGet, nil)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var nodesStatus models.NodesStatusResponse
	parseErr := responseData.DecodeBodyIntoTarget(&nodesStatus)
	return &nodesStatus, parseErr
}

func (nsg *NodesStatusGetter) WithClass(className string) *NodesStatusGetter {
	nsg.class = className
	return nsg
}

// WithOutput specifies nodes status response verbosity.
// Can be "minimal" or "verbose", defaults to "minimal".
func (nsg *NodesStatusGetter) WithOutput(output string) *NodesStatusGetter {
	nsg.output = output
	return nsg
}
