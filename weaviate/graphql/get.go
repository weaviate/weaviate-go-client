package graphql

import (
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
)

// Get data objects from weaviate using GraphQL
type Get struct {
	connection *connection.Connection
}

// Objects objects in result set
func (g *Get) Objects() *GetBuilder {
	return &GetBuilder{
		connection: g.connection,
	}
}
