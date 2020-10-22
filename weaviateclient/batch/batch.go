package batch

import "github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"

type API struct {
	Connection *connection.Connection
}

func (batch *API) ThingsBatcher() *ThingsBatcher {
	return &ThingsBatcher {
		connection: batch.Connection,
	}
}

func (batch *API) ActionsBatcher() *ActionsBatcher {
	return &ActionsBatcher{
		connection: batch.Connection,
	}
}