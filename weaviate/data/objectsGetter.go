package data

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// additionalProperties that have been set in the builder
type additionalProperties struct {
	withAdditionalInterpretation    bool
	withAdditionalClassification    bool
	withAdditionalNearestNeighbors  bool
	withAdditionalFeatureProjection bool
	withAdditionalVector            bool
}

// ObjectsGetter Builder to retrieve Things from weaviate
type ObjectsGetter struct {
	connection           *connection.Connection
	id                   string
	additionalProperties *additionalProperties
	withLimit            bool
	limit                int
}

// WithID specifies the uuid of the Thing that should be retrieved
//  if omitted a set of objects matching the builder specifications will be retrieved
func (getter *ObjectsGetter) WithID(id string) *ObjectsGetter {
	getter.id = id
	return getter
}

// WithAdditionalInterpretation include a description on how the corpus of the data object is interpreted by weaviate
func (getter *ObjectsGetter) WithAdditionalInterpretation() *ObjectsGetter {
	getter.additionalProperties.withAdditionalInterpretation = true
	return getter
}

// WithAdditionalClassification include information about the classifications
// may be nil if no classifications was executed on the object
func (getter *ObjectsGetter) WithAdditionalClassification() *ObjectsGetter {
	getter.additionalProperties.withAdditionalClassification = true
	return getter
}

// WithAdditionalNearestNeighbors show the nearest neighbors of this data object
func (getter *ObjectsGetter) WithAdditionalNearestNeighbors() *ObjectsGetter {
	getter.additionalProperties.withAdditionalNearestNeighbors = true
	return getter
}

// WithAdditionalFeatureProjection include a 2D projection of the objects for visualization
func (getter *ObjectsGetter) WithAdditionalFeatureProjection() *ObjectsGetter {
	getter.additionalProperties.withAdditionalFeatureProjection = true
	return getter
}

// WithAdditionalVector include the raw vector of the data object
func (getter *ObjectsGetter) WithAdditionalVector() *ObjectsGetter {
	getter.additionalProperties.withAdditionalVector = true
	return getter
}

// WithLimit of results
func (getter *ObjectsGetter) WithLimit(limit int) *ObjectsGetter {
	getter.withLimit = true
	getter.limit = limit
	return getter
}

// Do get the data object
func (getter *ObjectsGetter) Do(ctx context.Context) ([]*models.Object, error) {
	param := getAdditionalParams(getter.additionalProperties)
	if getter.withLimit {
		param += fmt.Sprintf("?limit=%v", getter.limit)
	}
	responseData, err := getObjectList(ctx, "/objects", getter.id, param, getter.connection)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 200 {
		return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
	}
	if getter.id == "" {
		var objects models.ObjectsListResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&objects)
		return objects.Objects, decodeErr
	}
	var object models.Object
	decodeErr := responseData.DecodeBodyIntoTarget(&object)
	return []*models.Object{&object}, decodeErr
}

// getAdditionalParams build the query URL parameters for the requested underscore properties
func getAdditionalParams(additionalProperties *additionalProperties) string {
	selectedProperties := make([]string, 0)

	if additionalProperties.withAdditionalInterpretation {
		selectedProperties = append(selectedProperties, "interpretation")
	}
	if additionalProperties.withAdditionalClassification {
		selectedProperties = append(selectedProperties, "classification")
	}
	if additionalProperties.withAdditionalVector {
		selectedProperties = append(selectedProperties, "vector")
	}
	if additionalProperties.withAdditionalFeatureProjection {
		selectedProperties = append(selectedProperties, "featureProjection")
	}
	if additionalProperties.withAdditionalNearestNeighbors {
		selectedProperties = append(selectedProperties, "nearestNeighbors")
	}

	params := strings.Join(selectedProperties, ",")
	if len(params) > 0 {
		params = fmt.Sprintf("?include=%v", params)
	}

	return params
}

func getObjectList(ctx context.Context, basePath string, id string, urlParameters string, con *connection.Connection) (*connection.ResponseData, error) {
	path := basePath
	if id != "" {
		path += fmt.Sprintf("/%v/%v", id, urlParameters)
	} else {
		path += fmt.Sprintf("%v", urlParameters)
	}

	responseData, err := con.RunREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		return responseData, except.NewDerivedWeaviateClientError(err)
	}
	return responseData, err
}
