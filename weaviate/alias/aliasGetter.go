package alias

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// AliasGetter builder object to get a schema class
type AliasGetter struct {
	connection *connection.Connection
	alias      string
}

// WithClassName specifies the class that will be fetched from schema
func (c *AliasGetter) WithAlias(aliasName string) *AliasGetter {
	c.alias = aliasName
	return c
}

// Do get a alias as specified in the builder
func (c *AliasGetter) Do(ctx context.Context) (*models.Alias, error) {
	responseData, err := c.connection.RunREST(ctx, fmt.Sprintf("/aliases/%s", c.alias), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var object models.Alias
		decodeErr := responseData.DecodeBodyIntoTarget(&object)
		return &object, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
