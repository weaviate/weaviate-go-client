package schema

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
)

// API Conntains all the builder objects required to access the weaviate schema API.
type API struct {
	Connection *connection.Connection
}

// Getter builder to get a weaviate schema
func (schema *API) Getter() *Getter {
	return &Getter{connection: schema.Connection}
}

// ClassCreator builder to create a weaviate schema class
func (schema *API) ClassCreator() *ClassCreator {
	return &ClassCreator{
		connection:   schema.Connection,
		semanticKind: paragons.SemanticKindThings, // Set the default
	}
}

// ClassDeleter builder to delete a weaviate schema class
func (schema *API) ClassDeleter() *ClassDeleter {
	return &ClassDeleter{
		connection:   schema.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

// AllDeleter builder to delete an entire schema from a weaviate
func (schema *API) AllDeleter() *AllDeleter {
	return &AllDeleter{
		connection: schema.Connection,
		schemaAPI:  schema,
	}
}

// PropertyCreator builder to add a property to an existing schema class
func (schema *API) PropertyCreator() *PropertyCreator {
	return &PropertyCreator{
		connection:   schema.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}
