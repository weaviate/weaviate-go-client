package data

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/pathbuilder"
)

// Checker builder to check data object's existence
type Checker struct {
	connection       *connection.Connection
	id               string
	className        string
	dbVersionSupport *db.VersionSupport
}

// WithID specifies the id of an data object to be checked
func (checker *Checker) WithID(id string) *Checker {
	checker.id = id
	return checker
}

// WithClassName specifies the class name of the object to be checked
func (checker *Checker) WithClassName(className string) *Checker {
	checker.className = className
	return checker
}

// Do check the specified data object if it exists in weaviate
func (checker *Checker) Do(ctx context.Context) (bool, error) {
	path := pathbuilder.ObjectsCheck(pathbuilder.Components{
		ID:        checker.id,
		Class:     checker.className,
		DBVersion: checker.dbVersionSupport,
	})
	responseData, err := checker.connection.RunREST(ctx, path, http.MethodHead, nil)
	exists := responseData.StatusCode == 204
	return exists, except.CheckResponseDataErrorAndStatusCode(responseData, err, 204, 404)
}
