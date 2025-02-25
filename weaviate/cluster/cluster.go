package cluster

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

// API collection of cluster related endpoints
type API struct {
	connection *connection.Connection
}

// New Cluster (nodes) api group from connection
func New(con *connection.Connection) *API {
	return &API{connection: con}
}

// NodesStatusGetter returns a builder to get the weaviate nodes status
func (cluster *API) NodesStatusGetter() *NodesStatusGetter {
	return &NodesStatusGetter{connection: cluster.connection}
}
