package weaviateclient

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"time"
)

const apiVersion = "v1"
const defaultTimeout = 15000 * time.Millisecond

// Config of the client endpoint
type Config struct {
	Host   string
	Scheme string
}

// WeaviateClient implementing the weaviate API
type WeaviateClient struct {
	connection *connection.Connection
	Misc       Misc
}

// New weaviate client from config
func New(config Config) *WeaviateClient {
	con := connection.NewConnection(config.Scheme, config.Host)

	return &WeaviateClient{
		connection: con,
		Misc:       Misc{connection: con},
	}
}
