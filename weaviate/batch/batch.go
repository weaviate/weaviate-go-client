package batch

import (
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/paragons"
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

// ThingsBatcher get a builder to create things in a batch
func (batch *API) ThingsBatcher() *ThingsBatcher {
	return &ThingsBatcher{
		connection: batch.connection,
	}
}

// ActionsBatcher get a builder to create actions in a batch
func (batch *API) ActionsBatcher() *ActionsBatcher {
	return &ActionsBatcher{
		connection: batch.connection,
	}
}

// ReferencePayloadBuilder get a builder to create a reference payload for a reference batch
func (batch *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection:       batch.connection,
		fromSemanticKind: paragons.SemanticKindThings,
		toSemanticKind:   paragons.SemanticKindThings,
	}
}

// ReferencesBatcher get a builder to add references in batch
func (batch *API) ReferencesBatcher() *ReferencesBatcher {
	return &ReferencesBatcher{
		connection: batch.connection,
		references: []*models.BatchReference{},
	}
}
