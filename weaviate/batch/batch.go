package batch

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/db"
	"github.com/weaviate/weaviate/entities/models"
)

// API for batch requests
type API struct {
	connection       *connection.Connection
	dbVersionSupport *db.VersionSupport
}

// New Batch api group from connection
func New(con *connection.Connection, dbVersionSupport *db.VersionSupport) *API {
	return &API{connection: con, dbVersionSupport: dbVersionSupport}
}

// ObjectsBatcher get a builder to create objects in a batch
func (batch *API) ObjectsBatcher() *ObjectsBatcher {
	return &ObjectsBatcher{
		connection: batch.connection,
	}
}

// ObjectsBatchDeleter returns a builder which deletes objects in bulk
func (batch *API) ObjectsBatchDeleter() *ObjectsBatchDeleter {
	return &ObjectsBatchDeleter{
		connection: batch.connection,
	}
}

// ReferencePayloadBuilder get a builder to create a reference payload for a reference batch
func (batch *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection: batch.connection,
		dbVersion:  batch.dbVersionSupport,
	}
}

// ReferencesBatcher get a builder to add references in batch
func (batch *API) ReferencesBatcher() *ReferencesBatcher {
	return &ReferencesBatcher{
		connection: batch.connection,
		references: []*models.BatchReference{},
	}
}
