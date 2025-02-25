package batch

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/pathbuilder"
	"github.com/weaviate/weaviate/entities/models"
)

// ReferencesBatcher builder to add multiple references in one batch request
type ReferencesBatcher struct {
	connection       *connection.Connection
	references       []*models.BatchReference
	consistencyLevel string
}

// WithReferences adds references to the current batch
func (rb *ReferencesBatcher) WithReferences(reference ...*models.BatchReference) *ReferencesBatcher {
	rb.references = append(rb.references, reference...)
	return rb
}

// WithReference adds a reference to the current batch
//
// Deprecated: Use WithReferences with the same syntax
func (rb *ReferencesBatcher) WithReference(reference *models.BatchReference) *ReferencesBatcher {
	return rb.WithReferences(reference)
}

// WithConsistencyLevel determines how many replicas must acknowledge a request
// before it is considered successful. Mutually exclusive with node_name param.
// Can be one of 'ALL', 'ONE', or 'QUORUM'.
func (rb *ReferencesBatcher) WithConsistencyLevel(cl string) *ReferencesBatcher {
	rb.consistencyLevel = cl
	return rb
}

// Do add all the references in the batch to weaviate
func (rb *ReferencesBatcher) Do(ctx context.Context) ([]models.BatchReferenceResponse, error) {
	path := pathbuilder.BatchReferences(pathbuilder.Components{
		ConsistencyLevel: rb.consistencyLevel,
	})
	responseData, responseErr := rb.connection.RunREST(ctx, path, http.MethodPost, rb.references)
	batchErr := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}
	var batchResponse []models.BatchReferenceResponse
	decodeErr := responseData.DecodeBodyIntoTarget(&batchResponse)
	return batchResponse, decodeErr
}
