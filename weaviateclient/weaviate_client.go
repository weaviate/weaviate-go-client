package weaviateclient

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/batch"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/data"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/misc"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/schema"
)



// Config of the client endpoint
type Config struct {
	Host   string
	Scheme string
}

// WeaviateClient implementing the weaviate API
type WeaviateClient struct {
	connection *connection.Connection
	Misc       misc.API
	Schema     schema.API
	Data data.API
	Batch batch.API
}

// New weaviate client from config
func New(config Config) *WeaviateClient {
	con := connection.NewConnection(config.Scheme, config.Host)

	return &WeaviateClient{
		connection: con,
		Misc:       misc.API{Connection: con},
		Schema:     schema.API{Connection: con},
		Data: data.API{Connection: con},
		Batch: batch.API{Connection: con},
	}
}
