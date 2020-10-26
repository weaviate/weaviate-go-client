package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
	"strings"
)

// ActionsGetter Builder to retrieve Actions from weaviate
type ActionsGetter struct {
	connection *connection.Connection
	uuid string
	underscoreProperties *underscoreProperties
}

// ThingsGetter Builder to retrieve Things from weaviate
type ThingsGetter struct {
	connection *connection.Connection
	uuid string
	underscoreProperties *underscoreProperties
}

// underscoreProperties that have been set in the builder
type underscoreProperties struct {
	withUnderscoreInterpretation bool
	withUnderscoreClassification bool
	withUnderscoreNearestNeighbors bool
	withUnderscoreFeatureProjection bool
	withUnderscoreVector bool
}

// WithID specifies the uuid of the Action that should be retrieved
//  if omitted a set of objects matching the builder specifications will be retrieved
func (getter *ActionsGetter) WithID(uuid string) *ActionsGetter {
	getter.uuid = uuid
	return getter
}

// WithID specifies the uuid of the Thing that should be retrieved
//  if omitted a set of objects matching the builder specifications will be retrieved
func (getter *ThingsGetter) WithID(uuid string) *ThingsGetter {
	getter.uuid = uuid
	return getter
}

// WithUnderscoreInterpretation include a description on how the corpus of the data object is interpreted by weaviate
func (getter *ActionsGetter) WithUnderscoreInterpretation() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreInterpretation = true
	return getter
}
// WithUnderscoreInterpretation include a description on how the corpus of the data object is interpreted by weaviate
func (getter *ThingsGetter) WithUnderscoreInterpretation() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreInterpretation = true
	return getter
}

// WithUnderscoreClassification include information about the classifications
// may be nil if no classifications was executed on the object
func (getter *ActionsGetter) WithUnderscoreClassification() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreClassification = true
	return getter
}
// WithUnderscoreClassification include information about the classifications
// may be nil if no classifications was executed on the object
func (getter *ThingsGetter) WithUnderscoreClassification() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreClassification = true
	return getter
}

// WithUnderscoreNearestNeighbors show the nearest neighbors of this data object
func (getter *ActionsGetter) WithUnderscoreNearestNeighbors() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreNearestNeighbors = true
	return getter
}
// WithUnderscoreNearestNeighbors show the nearest neighbors of this data object
func (getter *ThingsGetter) WithUnderscoreNearestNeighbors() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreNearestNeighbors = true
	return getter
}

// WithUnderscoreFeatureProjection include a 2D projection of the objects for visualization
func (getter *ActionsGetter) WithUnderscoreFeatureProjection() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreFeatureProjection = true
	return getter
}
// WithUnderscoreFeatureProjection include a 2D projection of the objects for visualization
func (getter *ThingsGetter) WithUnderscoreFeatureProjection() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreFeatureProjection = true
	return getter
}

// WithUnderscoreVector include the raw vector of the data object
func (getter *ActionsGetter) WithUnderscoreVector() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreVector = true
	return getter
}
// WithUnderscoreVector include the raw vector of the data object
func (getter *ThingsGetter) WithUnderscoreVector() *ThingsGetter {
	getter.underscoreProperties.withUnderscoreVector = true
	return getter
}

// Do get the data object
func (getter *ActionsGetter) Do(ctx context.Context) ([]*models.Action, error) {
	responseData, err := getObjectList(ctx, "/actions", getter.uuid, getParams(getter.underscoreProperties), getter.connection)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 200 {
		return nil, clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
	}
	if getter.uuid == "" {
		var actions models.ActionsListResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&actions)
		return actions.Actions, decodeErr
	}
	var action models.Action
	decodeErr := responseData.DecodeBodyIntoTarget(&action)
	return []*models.Action{&action}, decodeErr
}

// Do get the data object
func (getter *ThingsGetter) Do(ctx context.Context) ([]*models.Thing, error) {
	responseData, err := getObjectList(ctx, "/things", getter.uuid, getParams(getter.underscoreProperties), getter.connection)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 200 {
		return nil, clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
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

// getParams build the query URL parameters for the requested underscore properties
func getParams(underscoreProperties *underscoreProperties) string {
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
	}

	responseData, err := con.RunREST(ctx, path, http.MethodGet, nil)
	return responseData, err
}

