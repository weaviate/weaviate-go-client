package tokenize

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

// API is the tokenize API group.
type API struct {
	connection *connection.Connection
}

// New creates a new tokenize API group.
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Text returns a builder that tokenizes arbitrary text with a chosen
// tokenization method.
func (api *API) Text() *TextTokenizer {
	return &TextTokenizer{
		connection: api.connection,
	}
}

// Property returns a builder that tokenizes text using an existing property's
// tokenization configuration.
func (api *API) Property() *PropertyTokenizer {
	return &PropertyTokenizer{
		connection: api.connection,
	}
}
