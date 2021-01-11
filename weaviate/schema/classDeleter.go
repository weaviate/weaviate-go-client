package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
)

// ClassDeleter builder to remove a class from weaviate
type ClassDeleter struct {
	connection   *connection.Connection
	semanticKind semantics.Kind
	className    string
}

// WithClassName defines the name of the class that should be deleted
func (cd *ClassDeleter) WithClassName(className string) *ClassDeleter {
	cd.className = className
	return cd
}

// WithKind specifies the semantic kind that is used for the class about to be deleted
// If not called the builder defaults to `things`
func (cd *ClassDeleter) WithKind(semanticKind semantics.Kind) *ClassDeleter {
	cd.semanticKind = semanticKind
	return cd
}

// Do delete the class from the weaviate schema
func (cd *ClassDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v", cd.className)
	responseData, err := cd.connection.RunREST(ctx, path, http.MethodDelete, nil)
	return except.CheckResponnseDataErrorAndStatusCode(responseData, err, 200)
}
