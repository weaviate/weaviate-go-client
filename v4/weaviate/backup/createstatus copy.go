package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type CreateStatusRequest struct {
	request    models.BackupRestoreMeta
	connection *connection.Connection
}

func (e *CreateStatusRequest) WithID(id string) *CreateStatusRequest {
	e.request.ID = id
	return e
}

func (e *CreateStatusRequest) reset() {
	e.request = models.BackupRestoreMeta{}
}

func (ob *CreateStatusRequest) Do(ctx context.Context) ([]models.BackupRestoreMeta, error) {
	defer ob.reset()
	body := ob.request
	responseData, responseErr := ob.connection.RunREST(ctx,
		fmt.Sprintf("/backups/create/%v", ob.request.ID), http.MethodGet, body)
	restoreErr := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if restoreErr != nil {
		return nil, restoreErr
	}

	var parsedResponse []models.BackupRestoreMeta
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
