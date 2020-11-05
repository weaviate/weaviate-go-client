package schema

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
	"github.com/semi-technologies/weaviate-go-client/weaviate/models"
	"net/http"
)

// ClassCreator builder object to create a schema class
type ClassCreator struct {
	connection   *connection.Connection
	class        *models.Class
	semanticKind semantics.Kind
}

// WithClass specifies the class that will be added to the schema
func (cc *ClassCreator) WithClass(class *models.Class) *ClassCreator {
	cc.class = class
	return cc
}

// WithKind specifies the semantic kind that is used for the class about to be created
// If not called the builder defaults to `things`
func (cc *ClassCreator) WithKind(semanticKind semantics.Kind) *ClassCreator {
	cc.semanticKind = semanticKind
	return cc
}

// Do create a class in the schema as specified in the builder
func (cc *ClassCreator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v", string(cc.semanticKind))
	responseData, err := cc.connection.RunREST(ctx, path, http.MethodPost, cc.class)
	return except.CheckResponnseDataErrorAndStatusCode(responseData, err, 200)
}
