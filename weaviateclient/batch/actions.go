package batch

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type ActionsBatcher struct {
	connection *connection.Connection
	actions []*models.Action
}

func (ab *ActionsBatcher) WithObject(action *models.Action) *ActionsBatcher {
	ab.actions = append(ab.actions, action)
	return ab
}

func (ab *ActionsBatcher) Do(ctx context.Context) ([]models.ActionsGetResponse, error) {
	body := paragons.ActionsBatchRequestBody{
		Fields:  []string{"ALL"},
		Actions: ab.actions,
	}
	responseData, responseErr := ab.connection.RunREST(ctx, "/batching/actions", http.MethodPost, body)
	batchErr := clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}

	var parsedResponse []models.ActionsGetResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	if parseErr != nil {
		return nil, parseErr
	}
	return parsedResponse, nil
}