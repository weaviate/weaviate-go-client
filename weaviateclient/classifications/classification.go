package classifications

import "github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"

// API classifications API
type API struct {
	Connection *connection.Connection
}

// Scheduler get a builder to schedule a classification
func (api *API) Scheduler() *Scheduler {
	return &Scheduler{connection: api.Connection}
}

// Getter get a builder to retrieve a classification
func (api *API) Getter() *Getter {
	return &Getter{connection: api.Connection}
}
