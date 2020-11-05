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
		semanticKind: semantics.Things,
	}
}

// ActionsGetter get a builder to get an Action
func (data *API) ActionsGetter() *ActionsGetter {
	return &ActionsGetter{
		connection:           data.connection,
		underscoreProperties: &underscoreProperties{},
	}
}

// ThingsGetter get a builder to get a Thing
func (data *API) ThingsGetter() *ThingsGetter {
	return &ThingsGetter{
		connection:           data.connection,
		underscoreProperties: &underscoreProperties{},
	}
}

// Deleter get a builder to delete data objects
func (data *API) Deleter() *Deleter {
	return &Deleter{
		connection:   data.connection,
		semanticKind: semantics.Things,
	}
}

// Updater get a builder to update a data object
func (data *API) Updater() *Updater {
	return &Updater{
		connection:   data.connection,
		semanticKind: semantics.Things,
		withMerge:    false,
	}
}

// Validator get a builder to validate a data object definition
func (data *API) Validator() *Validator {
	return &Validator{
		connection:   data.connection,
		semanticKind: semantics.Things,
	}
}

// ReferencePayloadBuilder get a builder to create the payloads that reference an object
func (data *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection:   data.connection,
		semanticKind: semantics.Things,
	}
}

// ReferenceCreator get a builder to add references to data objects
func (data *API) ReferenceCreator() *ReferenceCreator {
	return &ReferenceCreator{
		connection:   data.connection,
		semanticKind: semantics.Things,
	}
}

// ReferenceReplacer get a builder to replace references on a data object
func (data *API) ReferenceReplacer() *ReferenceReplacer {
	return &ReferenceReplacer{
		connection:   data.connection,
		semanticKind: semantics.Things,
	}
}

// ReferenceDeleter get a builder to delete references on a data object
func (data *API) ReferenceDeleter() *ReferenceDeleter {
	return &ReferenceDeleter{
		connection:   data.connection,
		semanticKind: semantics.Things,
	}
}
