package alias

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// AliasUpdater builder object to update a alias
type AliasUpdater struct {
	connection *connection.Connection
	alias      *Alias
}

// WithAlias specifies the alias that will be updated to the schema
func (cu *AliasUpdater) WithAlias(alias *Alias) *AliasUpdater {
	cu.alias = alias
	return cu
}

// Do update a alias in the schema as specified in the builder
func (cu *AliasUpdater) Do(ctx context.Context) error {
	if cu.alias == nil {
		return except.NewWeaviateClientError(0, "an alias must be provided")
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
