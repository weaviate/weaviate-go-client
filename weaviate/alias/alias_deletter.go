package alias

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// AliasDeleter builder to remove a alias from weaviate
type AliasDeleter struct {
	connection *connection.Connection
	alias      string
}

// WithAliasName defines the name of the class that should be deleted
func (cd *AliasDeleter) WithAliasName(alias string) *AliasDeleter {
	cd.alias = alias
	return cd
}

// Do delete the alias from the weaviate schema
func (cd *AliasDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/aliases/%v", cd.alias)
	responseData, err := cd.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 204)
}
