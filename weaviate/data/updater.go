package data

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/semantics"
	"github.com/semi-technologies/weaviate-go-client/weaviate/models"
	"net/http"
)

// Updater builder to update property values in a data object
type Updater struct {
	connection     *connection.Connection
	semanticKind   semantics.Kind
	uuid           string
	className      string
	propertySchema models.PropertySchema
	withMerge      bool
}

// WithID specifies the uuid of the object about to be  updated
func (updater *Updater) WithID(uuid string) *Updater {
	updater.uuid = uuid
	return updater
}

// WithClassName specifies the class of the object about to be updated
func (updater *Updater) WithClassName(className string) *Updater {
	updater.className = className
	return updater
}

// WithSchema specifies the property schema of the class about to be updated
func (updater *Updater) WithSchema(propertySchema models.PropertySchema) *Updater {
	updater.propertySchema = propertySchema
	return updater
}

// WithKind specifies the semantic kind that is used for the data object
// If not called the builder defaults to `things`
func (updater *Updater) WithKind(semanticKind semantics.Kind) *Updater {
	updater.semanticKind = semanticKind
	return updater
}

// WithMerge indicates that the object should be merged with the existing object instead of replacing it
func (updater *Updater) WithMerge() *Updater {
	updater.withMerge = true
	return updater
}

// Do update the data object specified in the builder
func (updater *Updater) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v", string(updater.semanticKind), updater.uuid)
	httpMethod := http.MethodPut
	expectedStatuscode := 200
	if updater.withMerge {
		httpMethod = http.MethodPatch
		expectedStatuscode = 204
	}
	responseData, responseErr := updater.runUpdate(ctx, path, httpMethod)
	return except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, expectedStatuscode)
}

func (updater *Updater) runUpdate(ctx context.Context, path string, httpMethod string) (*connection.ResponseData, error) {
	if updater.semanticKind == semantics.Things {
		thing := models.Thing{
			Class:  updater.className,
			ID:     strfmt.UUID(updater.uuid),
			Schema: updater.propertySchema,
		}
		return updater.connection.RunREST(ctx, path, httpMethod, thing)
	}
	action := models.Action{
		Class:  updater.className,
		ID:     strfmt.UUID(updater.uuid),
		Schema: updater.propertySchema,
	}
	return updater.connection.RunREST(ctx, path, httpMethod, action)
}
