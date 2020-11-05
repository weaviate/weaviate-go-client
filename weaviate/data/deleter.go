package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
	"net/http"
)

// Deleter builder to delete a data object
type Deleter struct {
	connection   *connection.Connection
	uuid         string
	semanticKind semantics.Kind
}

// WithID specifies the uuid of the object about to be deleted
func (deleter *Deleter) WithID(uuid string) *Deleter {
	deleter.uuid = uuid
	return deleter
}

// WithKind specifies the semantic kind that is used for the data object
// If not called the builder defaults to `things`
func (deleter *Deleter) WithKind(semanticKind semantics.Kind) *Deleter {
	deleter.semanticKind = semanticKind
	return deleter
}

// Do delete the specified data object from weaviate
func (deleter *Deleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v", deleter.semanticKind, deleter.uuid)
	responseData, err := deleter.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return except.CheckResponnseDataErrorAndStatusCode(responseData, err, 204)
}
