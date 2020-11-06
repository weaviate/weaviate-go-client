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

// ObjectWrapper wrapping the result of a creation for both actions and things
type ObjectWrapper struct {
	Thing *models.Thing
	Action *models.Action
}

// Creator builder to create a data object in weaviate
type Creator struct {
	connection     *connection.Connection
	className      string
	uuid           string
	propertySchema models.PropertySchema
	semanticKind   semantics.Kind
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
func (creator *Creator) WithKind(semanticKind semantics.Kind) *Creator {
	creator.semanticKind = semanticKind
	return creator
}

// Do create the data object as specified in the builder
func (creator *Creator) Do(ctx context.Context) (*ObjectWrapper, error) {
	path := fmt.Sprintf("/%v", string(creator.semanticKind))

	var err error
	var responseData *connection.ResponseData
	if creator.semanticKind == semantics.Actions {
		action, _ := creator.PayloadAction()
		responseData, err = creator.connection.RunREST(ctx, path, http.MethodPost, action)
	} else {
		thing, _ := creator.PayloadThing()
		responseData, err = creator.connection.RunREST(ctx, path, http.MethodPost, thing)
	}
	respErr := except.CheckResponnseDataErrorAndStatusCode(responseData, err, 200)
	if respErr != nil {
		return nil, respErr
	}

	if creator.semanticKind == semantics.Actions {
		var resultAction models.Action
		parseErr := responseData.DecodeBodyIntoTarget(&resultAction)
		return &ObjectWrapper{
			Thing:  nil,
			Action: &resultAction,
		}, parseErr
	}
	var resultThing models.Thing
	parseErr := responseData.DecodeBodyIntoTarget(&resultThing)
	return &ObjectWrapper{
		Thing:  &resultThing,
		Action: nil,
	}, parseErr
}

// PayloadThing returns the data object payload which may be used in a batch request
func (creator *Creator) PayloadThing() (*models.Thing, error) {
	if creator.semanticKind != semantics.Things {
		return nil, except.NewDerivedWeaviateClientError(fmt.Errorf("builder has semantic kind action configured; please set the correct semantic type"))
	}
	thing := models.Thing{
		Class:  creator.className,
		Schema: creator.propertySchema,
	}
	if creator.uuid != "" {
		thing.ID = strfmt.UUID(creator.uuid)
	}
	return &thing, nil
}

// PayloadAction returns the data object payload which may be used in a batch request
func (creator *Creator) PayloadAction() (*models.Action, error) {
	if creator.semanticKind != semantics.Actions {
		return nil, except.NewDerivedWeaviateClientError(fmt.Errorf("builder has semantic kind thing configured; Please set the correct semantic type"))
	}
	action := models.Action{
		Class:  creator.className,
		Schema: creator.propertySchema,
	}
	if creator.uuid != "" {
		action.ID = strfmt.UUID(creator.uuid)
	}
	return &action, nil
}
