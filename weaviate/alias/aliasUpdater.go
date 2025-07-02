package alias

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// AliasUpdater builder object to update a schema class
type AliasUpdater struct {
	connection *connection.Connection
	alias      *models.Alias
}

// WithClass specifies the class properties that will be added to the schema
func (cu *AliasUpdater) WithAlias(alias *models.Alias) *AliasUpdater {
	cu.alias = alias
	return cu
}

// Do create a class in the schema as specified in the builder
func (cu *AliasUpdater) Do(ctx context.Context) error {
	if cu.alias == nil {
		return except.NewWeaviateClientError(0, "A alias must be provided")
	}
	path := fmt.Sprintf("/aliases/%v", cu.alias.Alias)
	updatePaylod := struct {
		Class string `json:"class"`
	}{
		Class: cu.alias.Class,
	}

	responseData, err := cu.connection.RunREST(ctx, path, http.MethodPut, updatePaylod)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
