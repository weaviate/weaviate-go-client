package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/except"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
	"strings"
)

// underscoreProperties that have been set in the builder
type underscoreProperties struct {
	withUnderscoreInterpretation    bool
	withUnderscoreClassification    bool
	withUnderscoreNearestNeighbors  bool
	withUnderscoreFeatureProjection bool
	withUnderscoreVector            bool
}

// ThingsGetter Builder to retrieve Things from weaviate
type ThingsGetter struct {
	connection           *connection.Connection
	uuid                 string
	underscoreProperties *underscoreProperties
	withLimit bool
	limit int
}

// WithID specifies the uuid of the Thing that should be retrieved
//  if omitted a set of objects matching the builder specifications will be retrieved
func (getter *ThingsGetter) WithID(uuid string) *ThingsGetter {
	getter.uuid = uuid
	return getter
}

// WithUnderscoreInterpretation include a description on how the corpus of the data object is interpreted by weaviate
func (getter *ThingsGetter) WithUnderscoreInterpretation() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreInterpretation = true
	return getter
}

// WithUnderscoreClassification include information about the classifications
// may be nil if no classifications was executed on the object
func (getter *ThingsGetter) WithUnderscoreClassification() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreClassification = true
	return getter
}

// WithUnderscoreNearestNeighbors show the nearest neighbors of this data object
func (getter *ThingsGetter) WithUnderscoreNearestNeighbors() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreNearestNeighbors = true
	return getter
}

// WithUnderscoreFeatureProjection include a 2D projection of the objects for visualization
func (getter *ThingsGetter) WithUnderscoreFeatureProjection() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreFeatureProjection = true
	return getter
}

// WithUnderscoreVector include the raw vector of the data object
func (getter *ThingsGetter) WithUnderscoreVector() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreVector = true
	return getter
}

// WithLimit of results
func (getter *ThingsGetter) WithLimit(limit int) *ThingsGetter {
	getter.withLimit = true
	getter.limit = limit
	return getter
}

// Do get the data object
func (getter *ThingsGetter) Do(ctx context.Context) ([]*models.Thing, error) {
	param := getUnderscoreParams(getter.underscoreProperties)
	if getter.withLimit {
		param += fmt.Sprintf("?limit=%v", getter.limit)
	}
	responseData, err := getObjectList(ctx, "/things", getter.uuid, param, getter.connection)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 200 {
		return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
	}
	if getter.uuid == "" {
		var things models.ThingsListResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&things)
		return things.Things, decodeErr
	}
	var thing models.Thing
	decodeErr := responseData.DecodeBodyIntoTarget(&thing)
	return []*models.Thing{&thing}, decodeErr
}

// getUnderscoreParams build the query URL parameters for the requested underscore properties
func getUnderscoreParams(underscoreProperties *underscoreProperties) string {
	selectedProperties := make([]string, 0)

	if underscoreProperties.withUnderscoreInterpretation {
		selectedProperties = append(selectedProperties, "_interpretation")
	}
	if underscoreProperties.withUnderscoreClassification {
		selectedProperties = append(selectedProperties, "_classification")
	}
	if underscoreProperties.withUnderscoreVector {
		selectedProperties = append(selectedProperties, "_vector")
	}
	if underscoreProperties.withUnderscoreFeatureProjection {
		selectedProperties = append(selectedProperties, "_feature_projection")
	}
	if underscoreProperties.withUnderscoreNearestNeighbors {
		selectedProperties = append(selectedProperties, "_nearest_neighbors")
	}

	params := strings.Join(selectedProperties, ",")
	if len(params) > 0 {
		params = fmt.Sprintf("?include=%v", params)
	}

	return params
}

func getObjectList(ctx context.Context, basePath string, uuid string, urlParameters string, con *connection.Connection) (*connection.ResponseData, error) {
	path := basePath
	if uuid != "" {
		path += fmt.Sprintf("/%v/%v", uuid, urlParameters)
	} else {
		path += fmt.Sprintf("%v", urlParameters)
	}

	responseData, err := con.RunREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		return responseData, except.NewDerivedWeaviateClientError(err)
	}
	return responseData, err
}
