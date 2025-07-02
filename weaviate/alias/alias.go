package alias

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

// API Conntains all the builder objects required to access the weaviate alias API.
type API struct {
	connection *connection.Connection
}

// New alias api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Getter builder to get a weaviate aliases
func (schema *API) Getter() *Getter {
	return &Getter{connection: schema.connection}
}

// AliasGetter builder to get a weaviate alias
func (schema *API) AliasGetter() *AliasGetter {
	return &AliasGetter{
		connection: schema.connection,
	}
}

// AliasCreator builder to create a weaviate alias
func (schema *API) AliasCreator() *AliasCreator {
	return &AliasCreator{
		connection: schema.connection,
	}
}

// AliasUpdater builder to update a weaviate alias
func (schema *API) AliasUpdater() *AliasUpdater {
	return &AliasUpdater{
		connection: schema.connection,
	}
}

// AliasDeleter builder to delete a weaviate alias
func (schema *API) AliasDeleter() *AliasDeleter {
	return &AliasDeleter{
		connection: schema.connection,
	}
}

// Alias represents the alias(softlink) to a collection in weaviate.
type Alias struct {
	// The name of the alias.
	Alias string `json:"alias,omitempty"`

	// class (collection) to which alias is assigned.
	Class string `json:"class,omitempty"`
}
