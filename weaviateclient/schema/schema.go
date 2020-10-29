package schema

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
)

// API Conntains all the builder objects required to access the weaviate schema API.
type API struct {
	connection *connection.Connection
}

// New Schema api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Getter builder to get a weaviate schema
func (schema *API) Getter() *Getter {
	return &Getter{connection: schema.connection}
}

// ClassCreator builder to create a weaviate schema class
func (schema *API) ClassCreator() *ClassCreator {
	return &ClassCreator{
		connection:   schema.connection,
		semanticKind: paragons.SemanticKindThings, // Set the default
	}
}

// ClassDeleter builder to delete a weaviate schema class
func (schema *API) ClassDeleter() *ClassDeleter {
	return &ClassDeleter{
		connection:   schema.connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

// AllDeleter builder to delete an entire schema from a weaviate
func (schema *API) AllDeleter() *AllDeleter {
	return &AllDeleter{
		connection: schema.connection,
		schemaAPI:  schema,
	}
}

// PropertyCreator builder to add a property to an existing schema class
func (schema *API) PropertyCreator() *PropertyCreator {
	return &PropertyCreator{
		connection:   schema.connection,
		semanticKind: paragons.SemanticKindThings,
	}
}
