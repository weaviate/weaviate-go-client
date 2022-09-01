package backup

import (
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/util"
)

// API for batch requests
type API struct {
	connection       *connection.Connection
	dbVersionSupport *util.DBVersionSupport
}

// New Batch api group from connection
func New(con *connection.Connection, dbVersionSupport *util.DBVersionSupport) *API {
	return &API{connection: con, dbVersionSupport: dbVersionSupport}
}

// builder to create a backup creator
func (batch *API) CreateRequester() *CreateRequest {
	return &CreateRequest{
		connection: batch.connection,
	}
}

func (batch *API) RestoreRequester() *RestoreRequest {
	return &RestoreRequest{
		connection: batch.connection,
	}
}

func (batch *API) RestoreStatusRequester() *RestoreStatusRequest {
	return &RestoreStatusRequest{
		connection: batch.connection,
	}
}

func (batch *API) CreateStatusRequester() *CreateStatusRequest {
	return &CreateStatusRequest{
		connection: batch.connection,
	}
}
