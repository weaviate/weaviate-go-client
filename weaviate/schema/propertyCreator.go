package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// PropertyCreator builder to create a property within a schema class
type PropertyCreator struct {
	connection *connection.Connection
	className  string
	property   *models.Property
}

// WithClassName defines the name of the class on which the property will be created
func (pc *PropertyCreator) WithClassName(className string) *PropertyCreator {
	pc.className = className
	return pc
}

// WithProperty defines the property object that will be added to the schema class
func (pc *PropertyCreator) WithProperty(property *models.Property) *PropertyCreator {
	pc.property = property
	return pc
}

// Do create the property on the class specified in the builder
func (pc *PropertyCreator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/properties", pc.className)
	responseData, err := pc.connection.RunREST(ctx, path, http.MethodPost, pc.property)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
