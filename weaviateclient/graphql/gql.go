package graphql

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
)

type API struct {
	Connection *connection.Connection
}

func (api *API) Get() *Get {
	return &Get{
		connection: api.Connection,
	}
}

//funnc ()


type rest interface {
	RunREST(ctx context.Context, path string, restMethod string, requestBody interface{}) (*connection.ResponseData, error)
}