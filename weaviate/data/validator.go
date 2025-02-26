package data

import (
	"context"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// Validator builder object to validate a class
type Validator struct {
	connection     *connection.Connection
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

// WithProperties specifies the property schema of the class about to be validated
func (validator *Validator) WithProperties(propertySchema models.PropertySchema) *Validator {
	validator.propertySchema = propertySchema
	return validator
}

// Do validate the data object specified in the builder
// Will return an error if the object is not valid or if there is a different error
func (validator *Validator) Do(ctx context.Context) error {
	path := "/objects/validate"
	object := models.Object{
		Class:      validator.className,
		ID:         strfmt.UUID(validator.uuid),
		Properties: validator.propertySchema,
	}
	responseData, err := validator.connection.RunREST(ctx, path, http.MethodPost, object)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
