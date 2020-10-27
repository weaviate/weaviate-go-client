package batch

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
)

// API for batch requests
type API struct {
	Connection *connection.Connection
}

// ThingsBatcher get a builder to create things in a batch
func (batch *API) ThingsBatcher() *ThingsBatcher {
	return &ThingsBatcher{
		connection: batch.Connection,
	}
}

// ActionsBatcher get a builder to create actions in a batch
func (batch *API) ActionsBatcher() *ActionsBatcher {
	return &ActionsBatcher{
		connection: batch.Connection,
	}
}

// ReferencePayloadBuilder get a builder to create a reference payload for a reference batch
func (batch *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection:       batch.Connection,
		fromSemanticKind: paragons.SemanticKindThings,
		toSemanticKind:   paragons.SemanticKindThings,
	}
}

// ReferencesBatcher get a builder to add references in batch
func (batch *API) ReferencesBatcher() *ReferencesBatcher {
	return &ReferencesBatcher{
		connection: batch.Connection,
		references: []*models.BatchReference{},
	}
}
