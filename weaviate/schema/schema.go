package schema

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
func (schema *API) Getter() *Getter {
	return &Getter{connection: schema.connection}
}

// ClassGetter builder to get a weaviate schema class
func (schema *API) ClassGetter() *ClassGetter {
	return &ClassGetter{
		connection: schema.connection,
	}
}

// ClassExistenceChecker builder to check if a class is part of a weaviate schema
func (schema *API) ClassExistenceChecker() *ClassExistenceChecker {
	return &ClassExistenceChecker{
		connection: schema.connection,
	}
}

// ClassCreator builder to create a weaviate schema class
func (schema *API) ClassCreator() *ClassCreator {
	return &ClassCreator{
		connection: schema.connection,
	}
}

// ClassUpdater builder to update a weaviate schema class
func (schema *API) ClassUpdater() *ClassUpdater {
	return &ClassUpdater{
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

// PropertyCreator builder to add a property to an existing schema class
func (schema *API) VectorAdder() *VectorAdder {
	return &VectorAdder{
		connection:   schema.connection,
		classGetter:  schema.ClassGetter(),
		classUpdater: schema.ClassUpdater(),
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

// TenantsCreator builder to add tenants to Class
func (schema *API) TenantsCreator() *TenantsCreator {
	return &TenantsCreator{
		connection: schema.connection,
	}
}

func (schema *API) TenantsUpdater() *TenantsUpdater {
	return &TenantsUpdater{
		connection: schema.connection,
	}
}

// TenantsDeleter builder to delete tenants from Class
func (schema *API) TenantsDeleter() *TenantsDeleter {
	return &TenantsDeleter{
		connection: schema.connection,
	}
}

// TenantsGetter builder to get tenants of Class
func (schema *API) TenantsGetter() *TenantsGetter {
	return &TenantsGetter{
		connection: schema.connection,
	}
}

// TenantsExists builder to check Class's tenants
func (schema *API) TenantsExists() *TenantsExists {
	return &TenantsExists{
		connection: schema.connection,
	}
}
