package misc

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
)

// API collection of endpoints that don't fit in other categories
type API struct {
	connection        *connection.Connection
	dbVersionProvider *db.VersionProvider
}

// New Misc (meta, .well-known) api group from connection
func New(con *connection.Connection, dbVersionProvider *db.VersionProvider) *API {
	return &API{connection: con, dbVersionProvider: dbVersionProvider}
}

// ReadyChecker retrieves weaviate ready status
func (misc *API) ReadyChecker() *ReadyChecker {
	return &ReadyChecker{connection: misc.connection, dbVersionProvider: misc.dbVersionProvider}
}

// LiveChecker retrieves weaviate live status
func (misc *API) LiveChecker() *LiveChecker {
	return &LiveChecker{connection: misc.connection, dbVersionProvider: misc.dbVersionProvider}
}

// OpenIDConfigurationGetter retrieves the Open ID configuration
// may be nil
func (misc *API) OpenIDConfigurationGetter() *OpenIDConfigGetter {
	return &OpenIDConfigGetter{connection: misc.connection}
}

// MetaGetter returns a builder to get the weaviate meta description
func (misc *API) MetaGetter() *MetaGetter {
	return &MetaGetter{connection: misc.connection}
}
