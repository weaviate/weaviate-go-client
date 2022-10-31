package data

import (
	"context"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ObjectWrapper wrapping the result of a creation for both actions and things
type ObjectWrapper struct {
	Object *models.Object
}

// Creator builder to create a data object in weaviate
type Creator struct {
	connection     *connection.Connection
	className      string
	uuid           string
	vector         []float32
	propertySchema models.PropertySchema
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

// WithProperties property values of the data object
func (creator *Creator) WithProperties(propertySchema models.PropertySchema) *Creator {
	creator.propertySchema = propertySchema
	return creator
}

func (creator *Creator) WithVector(vector []float32) *Creator {
	creator.vector = vector
	return creator
}

// Do create the data object as specified in the builder
func (creator *Creator) Do(ctx context.Context) (*ObjectWrapper, error) {
	var err error
	var responseData *connection.ResponseData
	object, _ := creator.PayloadObject()
	responseData, err = creator.connection.RunREST(ctx, "/objects", http.MethodPost, object)
	respErr := except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
	if respErr != nil {
		return nil, respErr
	}

	var resultObject models.Object
	parseErr := responseData.DecodeBodyIntoTarget(&resultObject)
	return &ObjectWrapper{
		Object: &resultObject,
	}, parseErr
}

// PayloadObject returns the data object payload which may be used in a batch request
func (creator *Creator) PayloadObject() (*models.Object, error) {
	object := models.Object{
		Class:      creator.className,
		Properties: creator.propertySchema,
		Vector:     creator.vector,
	}
	if creator.uuid != "" {
		object.ID = strfmt.UUID(creator.uuid)
	}
	return &object, nil
}
