package weaviateclient

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
)

type SemanticKind string

const SemanticKindThings SemanticKind = "things"
const SemanticKindActions SemanticKind = "actions"

// Config of the client endpoint
type Config struct {
	Host   string
	Scheme string
}

// WeaviateClient implementing the weaviate API
type WeaviateClient struct {
	connection *connection.Connection
	Misc       MiscAPI
	Schema     SchemaAPI
}

// New weaviate client from config
func New(config Config) *WeaviateClient {
	con := connection.NewConnection(config.Scheme, config.Host)

	return &WeaviateClient{
		connection: con,
		Misc:       MiscAPI{connection: con},
		Schema:     SchemaAPI{connection: con},
	}
}
