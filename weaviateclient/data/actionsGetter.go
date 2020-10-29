package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ActionsGetter Builder to retrieve Actions from weaviate
type ActionsGetter struct {
	connection           *connection.Connection
	uuid                 string
	underscoreProperties *underscoreProperties
	withLimit bool
	limit int
}

// WithID specifies the uuid of the Action that should be retrieved
//  if omitted a set of objects matching the builder specifications will be retrieved
func (getter *ActionsGetter) WithID(uuid string) *ActionsGetter {
	getter.uuid = uuid
	return getter
}

// WithUnderscoreInterpretation include a description on how the corpus of the data object is interpreted by weaviate
func (getter *ActionsGetter) WithUnderscoreInterpretation() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreInterpretation = true
	return getter
}

// WithUnderscoreClassification include information about the classifications
// may be nil if no classifications was executed on the object
func (getter *ActionsGetter) WithUnderscoreClassification() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreClassification = true
	return getter
}

// WithUnderscoreNearestNeighbors show the nearest neighbors of this data object
func (getter *ActionsGetter) WithUnderscoreNearestNeighbors() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreNearestNeighbors = true
	return getter
}

// WithUnderscoreFeatureProjection include a 2D projection of the objects for visualization
func (getter *ActionsGetter) WithUnderscoreFeatureProjection() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreFeatureProjection = true
	return getter
}

// WithUnderscoreVector include the raw vector of the data object
func (getter *ActionsGetter) WithUnderscoreVector() *ActionsGetter {
	getter.underscoreProperties.withUnderscoreVector = true
	return getter
}

// WithLimit of results
func (getter *ActionsGetter) WithLimit(limit int) *ActionsGetter {
	getter.withLimit = true
	getter.limit = limit
	return getter
}

// Do get the data object
func (getter *ActionsGetter) Do(ctx context.Context) ([]*models.Action, error) {
	param := getUnderscoreParams(getter.underscoreProperties)
	if getter.withLimit {
		param += fmt.Sprintf("?limit=%v", getter.limit)
	}
	responseData, err := getObjectList(ctx, "/actions", getter.uuid, param, getter.connection)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 200 {
		return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
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



