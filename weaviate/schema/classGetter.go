package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// ClassGetter builder object to get a schema class
type ClassGetter struct {
	connection *connection.Connection
	className  string
}

// WithClassName specifies the class that will be fetched from schema
func (c *ClassGetter) WithClassName(className string) *ClassGetter {
	c.className = className
	return c
}

// Do get a class from schema as specified in the builder
func (c *ClassGetter) Do(ctx context.Context) (*models.Class, error) {
	responseData, err := c.connection.RunREST(ctx, fmt.Sprintf("/schema/%s", c.className), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var object models.Class
		decodeErr := responseData.DecodeBodyIntoTarget(&object)
		return &object, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
