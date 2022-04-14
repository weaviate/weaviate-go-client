package batch

import (
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate/entities/models"
)

// API for batch requests
type API struct {
	connection *connection.Connection
}

// New Batch api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// ObjectsBatcher get a builder to create objects in a batch
func (batch *API) ObjectsBatcher() *ObjectsBatcher {
	return &ObjectsBatcher{
		connection: batch.connection,
	}
}

// ReferencePayloadBuilder get a builder to create a reference payload for a reference batch
func (batch *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection: batch.connection,
	}
}

// ReferencesBatcher get a builder to add references in batch
func (batch *API) ReferencesBatcher() *ReferencesBatcher {
	return &ReferencesBatcher{
		connection: batch.connection,
		references: []*models.BatchReference{},
	}
}
