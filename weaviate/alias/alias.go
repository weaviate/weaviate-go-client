package alias

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
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
func (schema *API) List() *AliasList {
	return &AliasList{connection: schema.connection}
}

// ClassGetter builder to get a weaviate schema class
func (schema *API) AliasGetter() *AliasGetter {
	return &AliasGetter{
		connection: schema.connection,
	}
}

// ClassCreator builder to create a weaviate schema class
func (schema *API) AliasCreator() *AliasCreator {
	return &AliasCreator{
		connection: schema.connection,
	}
}

// ClassUpdater builder to update a weaviate schema class
func (schema *API) AliasUpdater() *AliasUpdater {
	return &AliasUpdater{
		connection: schema.connection,
	}
}

// ClassDeleter builder to delete a weaviate schema class
func (schema *API) AliasDeleter() *AliasDeleter {
	return &AliasDeleter{
		connection: schema.connection,
	}
}
