package classifications

import "github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"

// API classifications API
type API struct {
	Connection *connection.Connection
}


func (api *API) Scheduler() *Scheduler {
	return &Scheduler{connection: api.Connection}
}

func (api *API) Getter() *Getter {
	return &Getter{connection: api.Connection}
}
