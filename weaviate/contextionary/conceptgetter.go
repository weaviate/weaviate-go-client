package contextionary

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// ConceptGetter builder to get weaviate concpets
type ConceptGetter struct {
	connection *connection.Connection
	concept    string
}

// WithConcept that should be retrieved
func (cg *ConceptGetter) WithConcept(concept string) *ConceptGetter {
	cg.concept = concept
	return cg
}

// Do get the concept
func (cg *ConceptGetter) Do(ctx context.Context) (*models.C11yWordsResponse, error) {
	path := fmt.Sprintf("/modules/text2vec-contextionary/concepts/%v", cg.concept)
	responseData, responseErr := cg.connection.RunREST(ctx, path, http.MethodGet, nil)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var concepts models.C11yWordsResponse
	parseErr := responseData.DecodeBodyIntoTarget(&concepts)
	if parseErr != nil {
		return nil, except.NewDerivedWeaviateClientError(parseErr)
	}
	return &concepts, nil
}
