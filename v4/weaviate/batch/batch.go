package batch

import (
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/util"
	"github.com/semi-technologies/weaviate/entities/models"
)

// API for batch requests
type API struct {
	connection *connection.Connection
	version    string
}

// New Batch api group from connection
func New(con *connection.Connection, version string) *API {
	return &API{connection: con, version: version}
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
		dbVersion:  util.NewDBVersionSupport(batch.version),
	}
}

// ReferencesBatcher get a builder to add references in batch
func (batch *API) ReferencesBatcher() *ReferencesBatcher {
	return &ReferencesBatcher{
		connection: batch.connection,
		references: []*models.BatchReference{},
	}
}
