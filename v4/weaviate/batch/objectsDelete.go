package batch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type ObjectsBatchDeleter struct {
	connection *connection.Connection
	filter     *models.BatchDelete
}

// WithFilter sets the filter for matching which objects to delete
func (ob *ObjectsBatchDeleter) WithFilter(filter *models.BatchDelete) *ObjectsBatchDeleter {
	ob.filter = filter
	return ob
}

// Do delete's all the objects which match the builder's filter
func (ob *ObjectsBatchDeleter) Do(ctx context.Context) (*models.BatchDeleteResponse, error) {
	if ob.filter == nil {
		return nil, fmt.Errorf("filter must be set prior to deletion, use WithFilter")
	}

	responseData, responseErr := ob.connection.RunREST(ctx, "/batch/objects", http.MethodDelete, ob.filter)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}

	var parsedResponse models.BatchDeleteResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return &parsedResponse, parseErr
}
