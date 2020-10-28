package graphql

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
)

type Get struct {
	connection *connection.Connection
}

func (g *Get) Things() *GetBuilder {
	return &GetBuilder{
		connection:   g.connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

func (g *Get) Actions() *GetBuilder {
	return &GetBuilder{
		connection:   g.connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

