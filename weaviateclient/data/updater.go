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

type Updater struct {
	connection *connection.Connection
	semanticKind paragons.SemanticKind
	uuid string
	className string
	propertySchema models.PropertySchema
}

func (updater *Updater) WithID(uuid string) *Updater{
	updater.uuid = uuid
	return updater
}

func (updater *Updater) WithClassName(className string) *Updater {
	updater.className = className
	return updater
}

func (updater *Updater) WithSchema(propertySchema models.PropertySchema) *Updater {
	updater.propertySchema = propertySchema
	return updater
}

func (updater *Updater) WithKind(semanticKind paragons.SemanticKind) *Updater {
	updater.semanticKind = semanticKind
	return updater
}

func (updater *Updater) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v", string(updater.semanticKind), updater.uuid)
	var responseData *connection.ResponseData
	var responseErr error
	if updater.semanticKind == paragons.SemanticKindThings {
		thing := models.Thing {
			Class:              updater.className,
			ID:                 strfmt.UUID(updater.uuid),
			Schema:             updater.propertySchema,
		}
		responseData, responseErr = updater.connection.RunREST(ctx, path, http.MethodPut, thing)
	} else {
		action := models.Action {
			Class:              updater.className,
			ID:                 strfmt.UUID(updater.uuid),
			Schema:             updater.propertySchema,
		}
		responseData, responseErr = updater.connection.RunREST(ctx, path, http.MethodPut, action)
	}
	if responseErr != nil {
		return responseErr
	}
	if responseData.StatusCode == 200 {
		return nil
	}
	return clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}