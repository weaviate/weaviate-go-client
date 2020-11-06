package batch

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// ActionsBatchRequestBody wrapping actions to a batch
type ActionsBatchRequestBody struct {
	Fields  []string         `json:"fields"`
	Actions []*models.Action `json:"actions"`
}

// ActionsBatcher builder to add multiple actions in one batch
type ActionsBatcher struct {
	connection *connection.Connection
	actions    []*models.Action
}

// WithObject add an object to the batch
func (ab *ActionsBatcher) WithObject(action *models.Action) *ActionsBatcher {
	ab.actions = append(ab.actions, action)
	return ab
}

// Do add all the objects in the builder to weaviate
func (ab *ActionsBatcher) Do(ctx context.Context) ([]models.ActionsGetResponse, error) {
	body := ActionsBatchRequestBody{
		Fields:  []string{"ALL"},
		Actions: ab.actions,
	}
	responseData, responseErr := ab.connection.RunREST(ctx, "/batching/actions", http.MethodPost, body)
	batchErr := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}

	var parsedResponse []models.ActionsGetResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
