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
	return &ActionGetter{connection: data.Connection}
}

func (data *API) ThingGetter() *ThingGetter {
	return &ThingGetter{connection: data.Connection}
}