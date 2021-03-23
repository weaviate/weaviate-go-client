package data

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v2/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ObjectsGetter Builder to retrieve Things from weaviate
type ObjectsGetter struct {
	connection           *connection.Connection
	id                   string
	additionalProperties []string
	withLimit            bool
	limit                int
}

// WithID specifies the uuid of the Thing that should be retrieved
//  if omitted a set of objects matching the builder specifications will be retrieved
func (getter *ObjectsGetter) WithID(id string) *ObjectsGetter {
	getter.id = id
	return getter
}

// WithVector include the raw vector of the data object
func (getter *ObjectsGetter) WithVector() *ObjectsGetter {
	getter.additionalProperties = append(getter.additionalProperties, "vector")
	return getter
}

// WithAdditional parameters such as for example: classification, featureProjection
func (getter *ObjectsGetter) WithAdditional(additional string) *ObjectsGetter {
	getter.additionalProperties = append(getter.additionalProperties, additional)
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
func getAdditionalParams(additionalProperties []string) string {
	selectedProperties := make([]string, 0)

	for _, additional := range additionalProperties {
		selectedProperties = append(selectedProperties, additional)
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
