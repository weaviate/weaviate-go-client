package batch

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// ReferencesBatcher builder to add multiple references in one batch request
type ReferencesBatcher struct {
	connection *connection.Connection
	references []*models.BatchReference
}

// WithReference adds a reference to the current batch
func (rb *ReferencesBatcher) WithReference(reference *models.BatchReference) *ReferencesBatcher {
	rb.references = append(rb.references, reference)
	return rb
}

// Do add all the references in the batch to weaviate
func (rb *ReferencesBatcher) Do(ctx context.Context) ([]models.BatchReferenceResponse, error) {
	responseData, responseErr := rb.connection.RunREST(ctx, "/batching/references", http.MethodPost, rb.references)
	batchErr := clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}
	var batchResponse []models.BatchReferenceResponse
	decodeErr := responseData.DecodeBodyIntoTarget(&batchResponse)
	return batchResponse, decodeErr
}
