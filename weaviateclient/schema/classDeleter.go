package schema

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"net/http"
)

// ClassDeleter builder to remove a class from weaviate
type ClassDeleter struct {
	connection   *connection.Connection
	semanticKind paragons.SemanticKind
	className    string
}

// WithClassName defines the name of the class that should be deleted
func (cd *ClassDeleter) WithClassName(className string) *ClassDeleter {
	cd.className = className
	return cd
}

// WithKind specifies the semantic kind that is used for the class about to be deleted
// If not called the builder defaults to `things`
func (cd *ClassDeleter) WithKind(semanticKind paragons.SemanticKind) *ClassDeleter {
	cd.semanticKind = semanticKind
	return cd
}

// Do delete the class from the weaviate schema
func (cd *ClassDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/%v", cd.semanticKind, cd.className)
	responseData, err := cd.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, err, 200)
}
