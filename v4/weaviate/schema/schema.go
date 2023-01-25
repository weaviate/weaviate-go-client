package schema

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
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

// ClassGetter builder to get a weaviate schema class
func (schema *API) ClassGetter() *ClassGetter {
	return &ClassGetter{
		connection: schema.connection,
	}
}

// ClassCreator builder to create a weaviate schema class
func (schema *API) ClassCreator() *ClassCreator {
	return &ClassCreator{
		connection: schema.connection,
	}
}

// ClassDeleter builder to delete a weaviate schema class
func (schema *API) ClassDeleter() *ClassDeleter {
	return &ClassDeleter{
		connection: schema.connection,
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
		connection: schema.connection,
	}
}

// ShardsGetter builder to get a weaviate class' shards
func (schema *API) ShardsGetter() *ShardsGetter {
	return &ShardsGetter{
		connection: schema.connection,
	}
}

// ShardUpdater builder to update a single shard
func (schema *API) ShardUpdater() *ShardUpdater {
	return &ShardUpdater{
		connection: schema.connection,
	}
}

// ShardsUpdater builder to update all shards within a class
func (schema *API) ShardsUpdater() *ShardsUpdater {
	return &ShardsUpdater{
		connection: schema.connection,
	}
}
