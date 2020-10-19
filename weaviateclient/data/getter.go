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

type ActionGetter struct {
	connection *connection.Connection
	uuid string
	underscoreProperties *underscoreProperties
}

type ThingGetter struct {
	connection *connection.Connection
	uuid string
	underscoreProperties *underscoreProperties
}

type underscoreProperties struct {
	withUnderscoreInterpretation bool
	withUnderscoreClassification bool
	withUnderscoreNearestNeighbors bool
	withUnderscoreFeatureProjection bool
	withUnderscoreVector bool
}

func (getter *ActionGetter) WithID(uuid string) *ActionGetter {
	getter.uuid = uuid
	return getter
}

func (getter *ThingGetter) WithID(uuid string) *ThingGetter {
	getter.uuid = uuid
	return getter
}


func (getter *ActionGetter) WithUnderscoreInterpretation() *ActionGetter {
	getter.underscoreProperties.withUnderscoreInterpretation = true
	return getter
}
func (getter *ThingGetter) WithUnderscoreInterpretation() *ThingGetter {
	getter.underscoreProperties.withUnderscoreInterpretation = true
	return getter
}

func (getter *ActionGetter) WithUnderscoreClassification() *ActionGetter {
	getter.underscoreProperties.withUnderscoreClassification = true
	return getter
}
func (getter *ThingGetter) WithUnderscoreClassification() *ThingGetter {
	getter.underscoreProperties.withUnderscoreClassification = true
	return getter
}

func (getter *ActionGetter) WithUnderscoreNearestNeighbors() *ActionGetter {
	getter.underscoreProperties.withUnderscoreNearestNeighbors = true
	return getter
}
func (getter *ThingGetter) WithUnderscoreNearestNeighbors() *ThingGetter {
	getter.underscoreProperties.withUnderscoreNearestNeighbors = true
	return getter
}

func (getter *ActionGetter) WithUnderscoreFeatureProjection() *ActionGetter {
	getter.underscoreProperties.withUnderscoreFeatureProjection = true
	return getter
}
func (getter *ThingGetter) WithUnderscoreFeatureProjection() *ThingGetter {
	getter.underscoreProperties.withUnderscoreFeatureProjection = true
	return getter
}

func (getter *ActionGetter) WithUnderscoreVector() *ActionGetter {
	getter.underscoreProperties.withUnderscoreVector = true
	return getter
}
func (getter *ThingGetter) WithUnderscoreVector() *ThingGetter {
	getter.underscoreProperties.withUnderscoreVector = true
	return getter
}


func (getter *ActionGetter) Do(ctx context.Context) ([]*models.Action, error) {
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

func (getter *ThingGetter) Do(ctx context.Context) ([]*models.Thing, error) {
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

