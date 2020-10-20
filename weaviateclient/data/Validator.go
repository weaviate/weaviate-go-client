package data

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type Validator struct {
	connection *connection.Connection
	semanticKind paragons.SemanticKind
	uuid string
	className string
	propertySchema models.PropertySchema
}

func (validator *Validator) WithID(uuid string) *Validator{
	validator.uuid = uuid
	return validator
}

func (validator *Validator) WithClassName(className string) *Validator {
	validator.className = className
	return validator
}

func (validator *Validator) WithSchema(propertySchema models.PropertySchema) *Validator {
	validator.propertySchema = propertySchema
	return validator
}

func (validator *Validator) WithKind(semanticKind paragons.SemanticKind) *Validator {
	validator.semanticKind = semanticKind
	return validator
}

func (validator *Validator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/validate", string(validator.semanticKind))
	var responseData *connection.ResponseData
	var err error
	if validator.semanticKind == paragons.SemanticKindThings {
		thing := models.Thing {
			Class:              validator.className,
			ID:                 strfmt.UUID(validator.uuid),
			Schema:             validator.propertySchema,
		}
		responseData, err = validator.connection.RunREST(ctx, path, http.MethodPost, thing)
	} else {
		action := models.Action {
			Class:              validator.className,
			ID:                 strfmt.UUID(validator.uuid),
			Schema:             validator.propertySchema,
		}
		responseData, err = validator.connection.RunREST(ctx, path, http.MethodPost, action)
	}
	if err != nil {
		return err
	}
	if responseData.StatusCode == 200 {
		return nil
	}
	return clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}
