package contextionary

import "github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"

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