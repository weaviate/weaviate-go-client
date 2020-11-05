package graphql

import (
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
)

// Get data objects from weaviate using GraphQL
type Get struct {
	connection *connection.Connection
}

// Things objects in result set
func (g *Get) Things() *GetBuilder {
	return &GetBuilder{
		connection:   g.connection,
		semanticKind: semantics.Things,
	}
}

// Actions objects in result set
func (g *Get) Actions() *GetBuilder {
	return &GetBuilder{
		connection:   g.connection,
		semanticKind: semantics.Actions,
	}
}

