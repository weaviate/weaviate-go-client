package batch

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type ReferencesBatcher struct {
	connection *connection.Connection
	references []*models.BatchReference
}

func (rb *ReferencesBatcher) WithReference(reference *models.BatchReference) *ReferencesBatcher {
	rb.references = append(rb.references, reference)
	return rb
}

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
