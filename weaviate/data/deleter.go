package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/pathbuilder"
)

// Deleter builder to delete a data object
type Deleter struct {
	connection       *connection.Connection
	id               string
	className        string
	consistencyLevel string
	tenant           string
	dbVersionSupport *db.VersionSupport
}

// WithID specifies the uuid of the object about to be deleted
func (deleter *Deleter) WithID(uuid string) *Deleter {
	deleter.id = uuid
	return deleter
}

// WithClassName specifies the class name of the object about to be deleted
func (deleter *Deleter) WithClassName(className string) *Deleter {
	deleter.className = className
	return deleter
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (deleter *Deleter) WithConsistencyLevel(cl string) *Deleter {
	deleter.consistencyLevel = cl
	return deleter
}

// WithTenant sets tenant, object should be deleted from
func (d *Deleter) WithTenant(tenant string) *Deleter {
	d.tenant = tenant
	return d
}

// Do delete the specified data object from weaviate
func (deleter *Deleter) Do(ctx context.Context) error {
	path := pathbuilder.ObjectsDelete(pathbuilder.Components{
		ID:               deleter.id,
		Class:            deleter.className,
		DBVersion:        deleter.dbVersionSupport,
		ConsistencyLevel: deleter.consistencyLevel,
		Tenant:           deleter.tenant,
	})
	responseData, err := deleter.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 204)
}
