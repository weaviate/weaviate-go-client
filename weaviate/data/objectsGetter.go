package data

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/pathbuilder"
	"github.com/weaviate/weaviate/entities/models"
)

// ObjectsGetter Builder to retrieve Things from weaviate
type ObjectsGetter struct {
	connection           *connection.Connection
	id                   string
	className            string
	additionalProperties []string
	withLimit            bool
	limit                int
	consistencyLevel     string
	nodeName             string
	dbVersionSupport     *db.VersionSupport
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

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'. Note that WithConsistencyLevel and
// WithNodeName are mutually exclusive.
func (getter *ObjectsGetter) WithConsistencyLevel(cl string) *ObjectsGetter {
	getter.consistencyLevel = cl
	return getter
}

// WithNodeName specifies the name of the target node which should fulfill the request.
// Note that WithNodeName and WithConsistencyLevel are mutually exclusive.
func (getter *ObjectsGetter) WithNodeName(name string) *ObjectsGetter {
	getter.nodeName = name
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
	basePath := getter.getPath()
	pathParams := getter.buildPathParams()
	return fmt.Sprintf("%s%v", basePath, pathParams)
}

func (getter *ObjectsGetter) getPath() string {
	return pathbuilder.ObjectsGet(pathbuilder.Components{
		ID:        getter.id,
		Class:     getter.className,
		DBVersion: getter.dbVersionSupport,
	})
}

func (getter *ObjectsGetter) buildPathParams() string {
	pathParams := make([]string, 0)

	additionalParams := buildAdditionalParams(getter.additionalProperties)
	if len(additionalParams) > 0 {
		pathParams = append(pathParams, additionalParams)
	}
	if getter.withLimit {
		pathParams = append(pathParams, fmt.Sprintf("limit=%v", getter.limit))
	}
	if len(getter.id) == 0 && len(getter.className) > 0 {
		if getter.dbVersionSupport.SupportsClassNameNamespacedEndpoints() {
			pathParams = append(pathParams, fmt.Sprintf("class=%v", getter.className))
		} else {
			getter.dbVersionSupport.WarnNotSupportedClassParameterInEndpointsForObjects()
		}
	}

	if getter.consistencyLevel != "" {
		pathParams = append(pathParams, fmt.Sprintf("consistency_level=%v", getter.consistencyLevel))
	}
	if getter.nodeName != "" {
		pathParams = append(pathParams, fmt.Sprintf("node_name=%v", getter.nodeName))
	}

	if len(pathParams) > 0 {
		return fmt.Sprintf("?%v", strings.Join(pathParams, "&"))
	}
	return ""
}

// buildAdditionalParams build the query URL parameters for the requested underscore properties
func buildAdditionalParams(additionalProperties []string) string {
	selectedProperties := make([]string, 0)
	selectedProperties = append(selectedProperties, additionalProperties...)

	params := strings.Join(selectedProperties, ",")
	if len(params) > 0 {
		params = fmt.Sprintf("include=%v", params)
	}

	return params
}
