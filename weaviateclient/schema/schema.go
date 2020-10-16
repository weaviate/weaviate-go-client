package schema

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	clientModels "github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
)

type SchemaAPI struct {
	Connection *connection.Connection
}

func (schema *SchemaAPI) Getter() *SchemaGetter {
	return &SchemaGetter{connection: schema.Connection}
}

func (schema *SchemaAPI) ClassCreator() *ClassCreator {
	return &ClassCreator{
		connection:   schema.Connection,
		semanticKind: clientModels.SemanticKindThings, // Set the default
	}
}

func (schema *SchemaAPI) ClassDeleter() *ClassDeleter {
	return &ClassDeleter{
		connection: schema.Connection,
		semanticKind: clientModels.SemanticKindThings,
	}
}

func (schema *SchemaAPI) AllDeleter() *AllDeleter {
	return &AllDeleter{
		connection: schema.Connection,
		schemaAPI: schema,
	}
}

