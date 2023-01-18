package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/util"
)

// Deleter builder to delete a data object
type Deleter struct {
	connection       *connection.Connection
	id               string
	className        string
	dbVersionSupport *util.DBVersionSupport
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

// Do delete the specified data object from weaviate
func (deleter *Deleter) Do(ctx context.Context) error {
	path := buildObjectsDeletePath(deleter.id, deleter.className, deleter.dbVersionSupport)
	responseData, err := deleter.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 204)
}
