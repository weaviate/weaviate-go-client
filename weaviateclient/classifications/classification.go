package classifications

import "github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"

// API classifications API
type API struct {
	connection *connection.Connection
}

// New Classification api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// Scheduler get a builder to schedule a classification
func (api *API) Scheduler() *Scheduler {
	return &Scheduler{connection: api.connection}
}

// Getter get a builder to retrieve a classification
func (api *API) Getter() *Getter {
	return &Getter{connection: api.connection}
}
