package data

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/util"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ObjectsGetter Builder to retrieve Things from weaviate
type ObjectsGetter struct {
	connection           *connection.Connection
	id                   string
	className            string
	additionalProperties []string
	withLimit            bool
	limit                int
	dbVersionSupport     *util.DBVersionSupport
}

// WithID specifies the uuid of the object that should be retrieved
// if omitted a set of objects matching the builder specifications will be retrieved
func (getter *ObjectsGetter) WithID(id string) *ObjectsGetter {
	getter.id = id
	return getter
}

// WithClassName specifies the class name of the object that should be retrieved
func (getter *ObjectsGetter) WithClassName(className string) *ObjectsGetter {
	getter.className = className
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
	responseData, err := getter.objectList(ctx)
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

func (getter *ObjectsGetter) objectList(ctx context.Context) (*connection.ResponseData, error) {
	responseData, err := getter.connection.RunREST(ctx, getter.buildPath(), http.MethodGet, nil)
	if err != nil {
		return responseData, except.NewDerivedWeaviateClientError(err)
	}
	return responseData, nil
}

func (getter *ObjectsGetter) buildPath() string {
	basePath := buildObjectsGetPath(getter.id, getter.className, getter.dbVersionSupport)

	params := buildAdditionalParams(getter.additionalProperties)
	if getter.withLimit && len(params) > 0 {
		params += fmt.Sprintf("&limit=%v", getter.limit)
	} else if getter.withLimit {
		params = fmt.Sprintf("?limit=%v", getter.limit)
	}

	return fmt.Sprintf("%s%v", basePath, params)
}

// buildAdditionalParams build the query URL parameters for the requested underscore properties
func buildAdditionalParams(additionalProperties []string) string {
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
