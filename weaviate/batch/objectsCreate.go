package batch

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/pathbuilder"
	"github.com/weaviate/weaviate/entities/models"
)

// ObjectsBatchRequestBody wrapping objects to a batch
type ObjectsBatchRequestBody struct {
	Fields  []string         `json:"fields"`
	Objects []*models.Object `json:"objects"`
}

// ObjectsBatcher builder to add multiple objects in one batch
type ObjectsBatcher struct {
	connection       *connection.Connection
	grpcClient       *connection.GrpcClient
	objects          []*models.Object
	consistencyLevel string
}

// WithObjects adds objects to the batch
func (ob *ObjectsBatcher) WithObjects(object ...*models.Object) *ObjectsBatcher {
	ob.objects = append(ob.objects, object...)
	return ob
}

// WithObject adds one object to the batch
//
// Deprecated: Use WithObjects with the same syntax
func (ob *ObjectsBatcher) WithObject(object *models.Object) *ObjectsBatcher {
	return ob.WithObjects(object)
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (ob *ObjectsBatcher) WithConsistencyLevel(cl string) *ObjectsBatcher {
	ob.consistencyLevel = cl
	return ob
}

func (ob *ObjectsBatcher) resetObjects() {
	ob.objects = []*models.Object{}
}

// Do add all the objects in the builder to weaviate
func (ob *ObjectsBatcher) Do(ctx context.Context) ([]models.ObjectsGetResponse, error) {
	defer ob.resetObjects()
	if ob.grpcClient != nil {
		return ob.runGRPC(ctx)
	}
	return ob.runREST(ctx)
}

func (ob *ObjectsBatcher) runREST(ctx context.Context) ([]models.ObjectsGetResponse, error) {
	body := ObjectsBatchRequestBody{
		Fields:  []string{"ALL"},
		Objects: ob.objects,
	}
	path := pathbuilder.BatchObjects(pathbuilder.Components{
		ConsistencyLevel: ob.consistencyLevel,
	})
	responseData, responseErr := ob.connection.RunREST(ctx, path, http.MethodPost, body)
	batchErr := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}

	var parsedResponse []models.ObjectsGetResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}

func (ob *ObjectsBatcher) runGRPC(ctx context.Context) ([]models.ObjectsGetResponse, error) {
	return ob.grpcClient.BatchObjects(ctx, ob.objects, ob.consistencyLevel)
}
