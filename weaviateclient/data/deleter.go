package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"net/http"
)

// Deleter builder to delete a data object
type Deleter struct {
	connection   *connection.Connection
	uuid         string
	semanticKind paragons.SemanticKind
}

// WithID specifies the uuid of the object about to be deleted
func (deleter *Deleter) WithID(uuid string) *Deleter {
	deleter.uuid = uuid
	return deleter
}

// WithKind specifies the semantic kind that is used for the data object
// If not called the builder defaults to `things`
func (deleter *Deleter) WithKind(semanticKind paragons.SemanticKind) *Deleter {
	deleter.semanticKind = semanticKind
	return deleter
}

// Do delete the specified data object from weaviate
func (deleter *Deleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v", deleter.semanticKind, deleter.uuid)
	responseData, err := deleter.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, err, 204)
}
