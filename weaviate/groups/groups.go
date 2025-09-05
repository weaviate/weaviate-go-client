package groups

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

type API struct {
	connection *connection.Connection
}

func New(connection *connection.Connection) *API {
	return &API{connection}
}

func (api *API) OIDC() *GroupsOIDC {
	return &GroupsOIDC{api.connection}
}
