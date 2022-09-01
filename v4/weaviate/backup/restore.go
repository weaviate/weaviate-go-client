package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type RestoreRequest struct {
	request     models.BackupRestoreRequest
	storageName string
	connection  *connection.Connection
}

// The ID of the backup. Must be URL-safe and work as a filesystem path, only lowercase, numbers, underscore, minus characters allowed.
func (e *RestoreRequest) WithID(id string) *RestoreRequest {
	e.request.ID = id
	return e
}

func (e *RestoreRequest) WithStorageName(storageName string) *RestoreRequest {
	e.storageName = storageName
	return e
}

func (e *RestoreRequest) WithInclude(include []string) *RestoreRequest {
	e.request.Include = include
	return e
}

func (e *RestoreRequest) WithExclude(exclude []string) *RestoreRequest {
	e.request.Exclude = exclude
	return e
}

func (ob *RestoreRequest) reset() {
	ob.request = models.BackupRestoreRequest{}
}

// Do add all the objects in the builder to weaviate
func (ob *RestoreRequest) Do(ctx context.Context) ([]models.BackupRestoreMeta, error) {
	defer ob.reset()
	body := ob
	responseData, responseErr := ob.connection.RunREST(ctx,
		fmt.Sprintf("/backups/%s/%s/restore", ob.storageName, ob.request.ID), http.MethodPost, body)
	restoreErr := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if restoreErr != nil {
		return nil, restoreErr
	}

	var parsedResponse []models.BackupRestoreMeta
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
