package batch

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type ThingsBatcher struct {
	connection *connection.Connection
	things []*models.Thing
}

func (tb *ThingsBatcher) WithObject(thing *models.Thing) *ThingsBatcher{
	tb.things = append(tb.things, thing)
	return tb
}

func (tb *ThingsBatcher) Do(ctx context.Context) ([]models.ThingsGetResponse, error) {
	body := paragons.ThingsBatchRequestBody{
		Fields: []string{"ALL"},
		Things: tb.things,
	}
	responseData, responseErr := tb.connection.RunREST(ctx, "/batching/things", http.MethodPost, body)
	batchErr := clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}

	var parsedResponse []models.ThingsGetResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	if parseErr != nil {
		return nil, parseErr
	}
	return parsedResponse, nil
}