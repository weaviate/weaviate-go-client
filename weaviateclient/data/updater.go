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
	withMerge bool
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

func (updater *Updater) WithMerge() *Updater {
	updater.withMerge = true
	return updater
}

func (updater *Updater) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v", string(updater.semanticKind), updater.uuid)
	httpMethod := http.MethodPut
	expectedStatuscode := 200
	if updater.withMerge {
		httpMethod = http.MethodPatch
		expectedStatuscode = 204
	}
	responseData, responseErr := updater.runUpdate(ctx, path, httpMethod)
	if responseErr != nil {
		return responseErr
	}
	if responseData.StatusCode == expectedStatuscode {
		return nil
	}
	return clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}

func (updater *Updater) runUpdate(ctx context.Context, path string, httpMethod string) (*connection.ResponseData, error) {
	if updater.semanticKind == paragons.SemanticKindThings {
		thing := models.Thing {
			Class:              updater.className,
			ID:                 strfmt.UUID(updater.uuid),
			Schema:             updater.propertySchema,
		}
		return updater.connection.RunREST(ctx, path, httpMethod, thing)
	} else {
		action := models.Action {
			Class:              updater.className,
			ID:                 strfmt.UUID(updater.uuid),
			Schema:             updater.propertySchema,
		}
		return updater.connection.RunREST(ctx, path, httpMethod, action)
	}
}
