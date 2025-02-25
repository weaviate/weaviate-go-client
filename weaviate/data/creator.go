package data

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// ObjectWrapper wrapping the result of a creation for both actions and things
type ObjectWrapper struct {
	Object *models.Object
}

// Creator builder to create a data object in weaviate
type Creator struct {
	connection       *connection.Connection
	className        string
	uuid             string
	vector           []float32
	vectors          models.Vectors
	propertySchema   models.PropertySchema
	consistencyLevel string
	tenant           string
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

func (creator *Creator) WithVectors(vectors models.Vectors) *Creator {
	creator.vectors = vectors
	return creator
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (creator *Creator) WithConsistencyLevel(cl string) *Creator {
	creator.consistencyLevel = cl
	return creator
}

// WithTenant sets tenant, object should be created for
func (c *Creator) WithTenant(tenant string) *Creator {
	c.tenant = tenant
	return c
}

// Do create the data object as specified in the builder
func (creator *Creator) Do(ctx context.Context) (*ObjectWrapper, error) {
	var err error
	var responseData *connection.ResponseData
	object, _ := creator.PayloadObject()

	path := creator.buildPath()
	responseData, err = creator.connection.RunREST(ctx, path, http.MethodPost, object)
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

func (creator *Creator) buildPath() string {
	path := "/objects"
	pathParams := url.Values{}

	if creator.consistencyLevel != "" {
		pathParams.Set("consistency_level", creator.consistencyLevel)
	}

	if len(pathParams) > 0 {
		path = fmt.Sprintf("%s?%v", path, pathParams.Encode())
	}

	return path
}

// PayloadObject returns the data object payload which may be used in a batch request
func (creator *Creator) PayloadObject() (*models.Object, error) {
	object := models.Object{
		Class:      creator.className,
		Properties: creator.propertySchema,
		Vector:     creator.vector,
		Vectors:    creator.vectors,
		Tenant:     creator.tenant,
	}
	if creator.uuid != "" {
		object.ID = strfmt.UUID(creator.uuid)
	}
	return &object, nil
}
