package data

import (
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
)

// API Contains all the builders required to access the weaviate data API
type API struct {
	Connection *connection.Connection
}

func (data *API) Creator() *Creator {
	return &Creator{
		connection:   data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

func (data *API) ActionGetter() *ActionGetter {
	return &ActionGetter{
		connection:           data.Connection,
		underscoreProperties: &underscoreProperties{},
	}
}

func (data *API) ThingGetter() *ThingGetter {
	return &ThingGetter{
		connection: data.Connection,
		underscoreProperties: &underscoreProperties{},
	}
}

func (data *API) Deleter() *Deleter {
	return &Deleter{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

func (data *API) Updater() *Updater {
	return &Updater{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
		withMerge: false,
	}
}

func (data *API) Validator() *Validator {
	return &Validator{
		connection:     data.Connection,
		semanticKind:   paragons.SemanticKindThings,
	}
}

func (data *API) ReferencePayloadBuilder() *ReferencePayloadBuilder {
	return &ReferencePayloadBuilder{
		connection:   data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

func (data *API) ReferenceCreator() *ReferenceCreator {
	return &ReferenceCreator{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}

func (data *API) ReferenceReplacer() *ReferenceReplacer {
	return &ReferenceReplacer{
		connection: data.Connection,
		semanticKind: paragons.SemanticKindThings,
	}
}