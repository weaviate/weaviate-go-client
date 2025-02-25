package data

import (
	"context"
	"net/http"
	"net/url"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/pathbuilder"
)

// Checker builder to check data object's existence
type Checker struct {
	connection       *connection.Connection
	id               string
	className        string
	tenant           string
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

// WithTenant sets tenant, object should be checked for
func (c *Checker) WithTenant(tenant string) *Checker {
	c.tenant = tenant
	return c
}

// Do check the specified data object if it exists in weaviate
func (checker *Checker) Do(ctx context.Context) (bool, error) {
	responseData, err := checker.connection.RunREST(ctx, checker.buildPath(), http.MethodHead, nil)
	exists := responseData.StatusCode == 204
	return exists, except.CheckResponseDataErrorAndStatusCode(responseData, err, 204, 404)
}

func (c *Checker) buildPath() string {
	endpoint := c.getPath()
	query := c.buildPathParams().Encode()
	if query == "" {
		return endpoint
	}
	return endpoint + "?" + query
}

func (c *Checker) getPath() string {
	return pathbuilder.ObjectsCheck(pathbuilder.Components{
		ID:        c.id,
		Class:     c.className,
		DBVersion: c.dbVersionSupport,
	})
}

func (c *Checker) buildPathParams() url.Values {
	pathParams := url.Values{}
	if c.tenant != "" {
		pathParams.Set("tenant", c.tenant)
	}
	return pathParams
}
