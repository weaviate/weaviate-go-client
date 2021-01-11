package batch

import (
	"context"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
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

// WithObject add an object to the batch
func (ob *ObjectsBatcher) WithObject(object *models.Object) *ObjectsBatcher {
	ob.objects = append(ob.objects, object)
	return ob
}

// Do add all the objects in the builder to weaviate
func (tb *ObjectsBatcher) Do(ctx context.Context) ([]models.ObjectsGetResponse, error) {
	body := ObjectsBatchRequestBody{
		Fields:  []string{"ALL"},
		Objects: tb.objects,
	}
	responseData, responseErr := tb.connection.RunREST(ctx, "/batch/objects", http.MethodPost, body)
	batchErr := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}

	var parsedResponse []models.ObjectsGetResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
