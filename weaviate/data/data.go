package data

import (
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
)

// API Contains all the builders required to access the weaviate data API
type API struct {
	connection *connection.Connection
}

// New {semanticKind} api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Creator get a builder to create a data object
func (data *API) Creator() *Creator {
	return &Creator{
		connection:   data.connection,
		semanticKind: semantics.Objects,
	}
}

// ObjectsGetter get a builder to get an Action
func (data *API) ObjectsGetter() *ObjectsGetter {
	return &ObjectsGetter{
		connection:           data.connection,
		additionalProperties: &additionalProperties{},
	}
}

// Deleter get a builder to delete data objects
func (data *API) Deleter() *Deleter {
	return &Deleter{
		connection:   data.connection,
		semanticKind: semantics.Objects,
	}
}

// Updater get a builder to update a data object
func (data *API) Updater() *Updater {
	return &Updater{
		connection:   data.connection,
		semanticKind: semantics.Objects,
		withMerge:    false,
	}
}

// Validator get a builder to validate a data object definition
func (data *API) Validator() *Validator {
	return &Validator{
		connection:   data.connection,
		semanticKind: semantics.Objects,
	}
}

// ReferencePayloadBuilder get a builder to create the payloads that reference an object
func (data *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection:   data.connection,
		semanticKind: semantics.Objects,
	}
}

// ReferenceCreator get a builder to add references to data objects
func (data *API) ReferenceCreator() *ReferenceCreator {
	return &ReferenceCreator{
		connection:   data.connection,
		semanticKind: semantics.Objects,
	}
}

// ReferenceReplacer get a builder to replace references on a data object
func (data *API) ReferenceReplacer() *ReferenceReplacer {
	return &ReferenceReplacer{
		connection:   data.connection,
		semanticKind: semantics.Objects,
	}
}

// ReferenceDeleter get a builder to delete references on a data object
func (data *API) ReferenceDeleter() *ReferenceDeleter {
	return &ReferenceDeleter{
		connection:   data.connection,
		semanticKind: semantics.Objects,
	}
}
