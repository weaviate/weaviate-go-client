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

// Creator builder to create a data object in weaviate
type Creator struct {
	connection *connection.Connection
	className string
	uuid string
	propertySchema models.PropertySchema
	semanticKind paragons.SemanticKind
}

// WithClassName indicates what class the data object is associated with
func (creator *Creator) WithClassName(name string) *Creator {
	creator.className = name
	return creator
}

// WithID if specified the object will be created under this uuid
// weaviate will generate a uuid if this was not called or an empty string is specified.
func (creator *Creator) WithID(uuid string) *Creator {
	creator.uuid = uuid
	return creator
}

// WithSchema property values of the data object
func (creator *Creator) WithSchema(propertySchema models.PropertySchema) *Creator {
	creator.propertySchema = propertySchema
	return creator
}

// WithKind specifies the semantic kind that is used for the data object
// If not called the builder defaults to `things`
func (creator *Creator) WithKind(semanticKind paragons.SemanticKind) *Creator {
	creator.semanticKind = semanticKind
	return creator
}

// Do create the data object as specified in the builder
func (creator *Creator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v", string(creator.semanticKind))

	var err error
	var responseData *connection.ResponseData
	if creator.semanticKind == paragons.SemanticKindActions {
		action := models.Action{
			Class:              creator.className,
			Schema:             creator.propertySchema,
		}
		if creator.uuid != "" {
			action.ID = strfmt.UUID(creator.uuid)
		}
		responseData, err = creator.connection.RunREST(ctx, path, http.MethodPost, action)
	} else {
		thing := models.Thing{
			Class: creator.className,
			Schema: creator.propertySchema,
		}
		if creator.uuid != "" {
			thing.ID = strfmt.UUID(creator.uuid)
		}
		responseData, err = creator.connection.RunREST(ctx, path, http.MethodPost, thing)
	}
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, err, 200)
}

