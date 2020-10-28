package batch

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/except"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// ThingsBatcher builder to add multiple things in one batch
type ThingsBatcher struct {
	connection *connection.Connection
	things     []*models.Thing
}

// WithObject add an object to the batch
func (tb *ThingsBatcher) WithObject(thing *models.Thing) *ThingsBatcher {
	tb.things = append(tb.things, thing)
	return tb
}

// Do add all the objects in the builder to weaviate
func (tb *ThingsBatcher) Do(ctx context.Context) ([]models.ThingsGetResponse, error) {
	body := paragons.ThingsBatchRequestBody{
		Fields: []string{"ALL"},
		Things: tb.things,
	}
	responseData, responseErr := tb.connection.RunREST(ctx, "/batching/things", http.MethodPost, body)
	batchErr := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if batchErr != nil {
		return nil, batchErr
	}

	var parsedResponse []models.ThingsGetResponse
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
