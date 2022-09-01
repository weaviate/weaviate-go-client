package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type CreateRequest struct {
	request     models.BackupCreateRequest
	storageName string
	connection  *connection.Connection
}

// The ID of the backup. Must be URL-safe and work as a filesystem path, only lowercase, numbers, underscore, minus characters allowed.
func (e *CreateRequest) WithID(id string) *CreateRequest {
	e.request.ID = id
	return e
}

// The name of the storage to use.
func (e *CreateRequest) WithStorageName(storageName string) *CreateRequest {
	e.storageName = storageName
	return e
}

func (e *CreateRequest) WithInclude(include []string) *CreateRequest {
	e.request.Include = include
	return e
}

func (e *CreateRequest) WithExclude(exclude []string) *CreateRequest {
	e.request.Exclude = exclude
	return e
}

func (ob *CreateRequest) reset() {
	ob.request = models.BackupCreateRequest{}
}

// Do add all the objects in the builder to weaviate
func (ob *CreateRequest) Do(ctx context.Context) ([]models.BackupCreateMeta, error) {
	defer ob.reset()
	body := ob.request
	responseData, responseErr := ob.connection.RunREST(ctx, fmt.Sprintf("/backups/%s/create", ob.storageName), http.MethodPost, body)
	CreateErr := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if CreateErr != nil {
		return nil, CreateErr
	}

	var parsedResponse []models.BackupCreateMeta
	parseErr := responseData.DecodeBodyIntoTarget(&parsedResponse)
	return parsedResponse, parseErr
}
