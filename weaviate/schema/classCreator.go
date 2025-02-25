package schema

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// ClassCreator builder object to create a schema class
type ClassCreator struct {
	connection *connection.Connection
	class      *models.Class
}

// WithClass specifies the class that will be added to the schema
func (cc *ClassCreator) WithClass(class *models.Class) *ClassCreator {
	cc.class = class
	return cc
}

// Do create a class in the schema as specified in the builder
func (cc *ClassCreator) Do(ctx context.Context) error {
	responseData, err := cc.connection.RunREST(ctx, "/schema", http.MethodPost, cc.class)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
