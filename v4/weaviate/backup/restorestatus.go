package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type RestoreStatusRequest struct {
	request    models.BackupRestoreMeta
	connection *connection.Connection
}

// The ID of the backup. Must be URL-safe and work as a filesystem path, only lowercase, numbers, underscore, minus characters allowed.
func (e *RestoreStatusRequest) WithID(id string) *RestoreStatusRequest {
	e.request.ID = id
	return e
}

func (e *RestoreStatusRequest) reset() {
	e.request = models.BackupRestoreMeta{}
}

func (ob *RestoreStatusRequest) Do(ctx context.Context) ([]models.BackupRestoreMeta, error) {
	defer ob.reset()
	body := ob.request
	responseData, responseErr := ob.connection.RunREST(ctx,
		fmt.Sprintf("/backups/restore/%v", ob.request.ID), http.MethodGet, body)
	restoreErr := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if restoreErr != nil {
		return nil, restoreErr
	}

	var parsedResponse []models.BackupRestoreMeta
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
