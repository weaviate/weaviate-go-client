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

// Validator builder object to validate a class
type Validator struct {
	connection     *connection.Connection
	semanticKind   semantics.Kind
	uuid           string
	className      string
	propertySchema models.PropertySchema
}

// WithID specifies the uuid of the object about to be validated
func (validator *Validator) WithID(uuid string) *Validator {
	validator.uuid = uuid
	return validator
}

// WithClassName specifies the class of the object about to be validated
func (validator *Validator) WithClassName(className string) *Validator {
	validator.className = className
	return validator
}

// WithSchema specifies the property schema of the class about to be validated
func (validator *Validator) WithSchema(propertySchema models.PropertySchema) *Validator {
	validator.propertySchema = propertySchema
	return validator
}

// WithKind specifies the semantic kind that is used for the data object
// If not called the builder defaults to `things`
func (validator *Validator) WithKind(semanticKind semantics.Kind) *Validator {
	validator.semanticKind = semanticKind
	return validator
}

// Do validate the data object specified in the builder
// Will return an error if the object is not valid or if there is a different error
func (validator *Validator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/validate", string(validator.semanticKind))
	var responseData *connection.ResponseData
	var err error
	if validator.semanticKind == semantics.Things {
		thing := models.Thing{
			Class:  validator.className,
			ID:     strfmt.UUID(validator.uuid),
			Schema: validator.propertySchema,
		}
		responseData, err = validator.connection.RunREST(ctx, path, http.MethodPost, thing)
	} else {
		action := models.Action{
			Class:  validator.className,
			ID:     strfmt.UUID(validator.uuid),
			Schema: validator.propertySchema,
		}
		responseData, err = validator.connection.RunREST(ctx, path, http.MethodPost, action)
	}
	return except.CheckResponnseDataErrorAndStatusCode(responseData, err, 200)
}
