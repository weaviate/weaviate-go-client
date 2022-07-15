package batch

import (
	"context"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

// ObjectsBatchRequestBody wrapping objects to a batch
type ObjectsBatchRequestBody struct {
	Fields  []string         `json:"fields"`
	Objects []*models.Object `json:"objects"`
}

// ObjectsBatcher builder to add multiple objects in one batch
type ObjectsBatcher struct {
	connection *connection.Connection
	objects    []*models.Object
}

// WithObject adds an object to the batch
func (ob *ObjectsBatcher) WithObject(object *models.Object) *ObjectsBatcher {
	ob.objects = append(ob.objects, object)
	return ob
}

// WithObjects adds objects to the batch
func (ob *ObjectsBatcher) WithObjects(objects ...*models.Object) *ObjectsBatcher {
	ob.objects = append(ob.objects, objects...)
	return ob
}

func (ob *ObjectsBatcher) resetObjects() {
	ob.objects = []*models.Object{}
}

// Do add all the objects in the builder to weaviate
func (ob *ObjectsBatcher) Do(ctx context.Context) ([]models.ObjectsGetResponse, error) {
	defer ob.resetObjects()
	body := ObjectsBatchRequestBody{
		Fields:  []string{"ALL"},
		Objects: ob.objects,
	}
	responseData, responseErr := ob.connection.RunREST(ctx, "/batch/objects", http.MethodPost, body)
	batchErr := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}

	var parsedResponse []models.ObjectsGetResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
