package contextionary

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
)

// API for the contextionary endpoint
type API struct {
	Connection *connection.Connection
}

// ConceptsGetter get builder to query weaviate concepts
func (c11y *API) ConceptsGetter() *ConceptGetter {
	return &ConceptGetter{
		connection: c11y.Connection,
	}
}

// ExtensionCreator get a builder to extend weaviates contextionary
func (c11y *API) ExtensionCreator() *ExtensionCreator {
	return &ExtensionCreator{
		connection: c11y.Connection,
		extension: &models.C11yExtension{
			Weight: 1.0,
		},
	}
}
