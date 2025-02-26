package data

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
)

// API Contains all the builders required to access the weaviate data API
type API struct {
	connection       *connection.Connection
	dbVersionSupport *db.VersionSupport
}

// New {semanticKind} api group from connection
func New(con *connection.Connection, dbVersionSupport *db.VersionSupport) *API {
	return &API{connection: con, dbVersionSupport: dbVersionSupport}
}

// Creator get a builder to create a data object
func (data *API) Creator() *Creator {
	return &Creator{
		connection: data.connection,
	}
}

// ObjectsGetter get a builder to get an Action
func (data *API) ObjectsGetter() *ObjectsGetter {
	return &ObjectsGetter{
		connection:           data.connection,
		additionalProperties: []string{},
		dbVersionSupport:     data.dbVersionSupport,
	}
}

// Deleter get a builder to delete data objects
func (data *API) Deleter() *Deleter {
	return &Deleter{
		connection:       data.connection,
		dbVersionSupport: data.dbVersionSupport,
	}
}

// Updater get a builder to update a data object
func (data *API) Updater() *Updater {
	return &Updater{
		connection:       data.connection,
		withMerge:        false,
		dbVersionSupport: data.dbVersionSupport,
	}
}

// Validator get a builder to validate a data object definition
func (data *API) Validator() *Validator {
	return &Validator{
		connection: data.connection,
	}
}

// Checker get a builder to check data object existence
func (data *API) Checker() *Checker {
	return &Checker{
		connection:       data.connection,
		dbVersionSupport: data.dbVersionSupport,
	}
}

// ReferencePayloadBuilder get a builder to create the payloads that reference an object
func (data *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection:       data.connection,
		dbVersionSupport: data.dbVersionSupport,
	}
}

// ReferenceCreator get a builder to add references to data objects
func (data *API) ReferenceCreator() *ReferenceCreator {
	return &ReferenceCreator{
		connection:       data.connection,
		dbVersionSupport: data.dbVersionSupport,
	}
}

// ReferenceReplacer get a builder to replace references on a data object
func (data *API) ReferenceReplacer() *ReferenceReplacer {
	return &ReferenceReplacer{
		connection:       data.connection,
		dbVersionSupport: data.dbVersionSupport,
	}
}

// ReferenceDeleter get a builder to delete references on a data object
func (data *API) ReferenceDeleter() *ReferenceDeleter {
	return &ReferenceDeleter{
		connection:       data.connection,
		dbVersionSupport: data.dbVersionSupport,
	}
}
