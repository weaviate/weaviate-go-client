package data

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
)

// API Contains all the builders required to access the weaviate data API
type API struct {
	Connection *connection.Connection
}

// Creator get a builder to create a data object
func (data *API) Creator() *Creator {
	return &Creator{
		connection:   data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

// ActionsGetter get a builder to get an Action
func (data *API) ActionsGetter() *ActionsGetter {
	return &ActionsGetter{
		connection:           data.Connection,
		underscoreProperties: &underscoreProperties{},
	}
}

// ThingsGetter get a builder to get a Thing
func (data *API) ThingsGetter() *ThingsGetter {
	return &ThingsGetter{
		connection: data.Connection,
		underscoreProperties: &underscoreProperties{},
	}
}

// Deleter get a builder to delete data objects
func (data *API) Deleter() *Deleter {
	return &Deleter{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

// Updater get a builder to update a data object
func (data *API) Updater() *Updater {
	return &Updater{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
		withMerge: false,
	}
}

// Validator get a builder to validate a data object definition
func (data *API) Validator() *Validator {
	return &Validator{
		connection:     data.Connection,
		semanticKind:   paragons.SemanticKindThings,
	}
}

// ReferencePayloadBuilder get a builder to create the payloads that reference an object
func (data *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection:   data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

// ReferenceCreator get a builder to add references to data objects
func (data *API) ReferenceCreator() *ReferenceCreator {
	return &ReferenceCreator{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

// ReferenceReplacer get a builder to replace references on a data object
func (data *API) ReferenceReplacer() *ReferenceReplacer {
	return &ReferenceReplacer{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

// ReferenceDeleter get a builder to delete references on a data object
func (data *API) ReferenceDeleter() *ReferenceDeleter {
	return &ReferenceDeleter{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}
