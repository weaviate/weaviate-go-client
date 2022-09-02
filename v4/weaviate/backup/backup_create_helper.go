package backup

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

const waitTimeoutCreate = time.Second

type backupCreateHelper struct {
	connection *connection.Connection
}

func (h *backupCreateHelper) create(ctx context.Context, includeClasses, excludeClasses []string,
	storageName, backupID string,
) (*models.BackupCreateMeta, error) {
	return h.createByEndpoint(ctx, endpointCreate(storageName), backupID, includeClasses, excludeClasses)
}

func (h *backupCreateHelper) createByEndpoint(ctx context.Context, endpoint, backupID string, includeClasses, excludeClasses []string) (*models.BackupCreateMeta, error) {
	data := models.BackupCreateRequest{
		Include: includeClasses,
		Exclude: excludeClasses,
		ID:      backupID,
	}
	return h.runREST(ctx, endpoint, http.MethodPost, data)
}

func (h *backupCreateHelper) statusCreate(ctx context.Context,
	storageName, backupID string,
) (*models.BackupCreateMeta, error) {
	return h.statusCreateByEndpoint(ctx, endpointStatusCreate(storageName, backupID))
}

func (h *backupCreateHelper) statusCreateByEndpoint(ctx context.Context, endpoint string) (*models.BackupCreateMeta, error) {
	return h.runREST(ctx, endpoint, http.MethodGet, nil)
}

func (h *backupCreateHelper) createAndWaitForCompletion(ctx context.Context,
	storageName, backupID string, includeClasses, excludeClasses []string,
) (*models.BackupCreateMeta, error) {
	endpoint := endpointCreate(storageName)
	if _, err := h.createByEndpoint(ctx, endpoint, backupID, includeClasses, excludeClasses); err != nil {
		return nil, err
	}
	endpoint = endpointStatusCreate(storageName, backupID)
	for {
		meta, err := h.statusCreateByEndpoint(ctx, endpoint)
		if err != nil {
			return nil, err
		}
		switch *meta.Status {
		case models.BackupCreateMetaStatusSUCCESS, models.BackupCreateMetaStatusFAILED:
			return meta, nil
		default:
			time.Sleep(waitTimeoutCreate)
		}
	}
}

func (h *backupCreateHelper) runREST(ctx context.Context, endpoint, httpMethod string, data interface{}) (*models.BackupCreateMeta, error) {
	responseData, err := h.connection.RunREST(ctx, endpoint, httpMethod, data)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var obj models.BackupCreateMeta
		decodeErr := responseData.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewDerivedWeaviateClientError(err)
}

func endpointCreate(storageName string) string {
	return fmt.Sprintf("/backups/%s", storageName)
}

func endpointStatusCreate(storageName, ID string) string {
	return fmt.Sprintf("/backups/%s/%s", storageName, ID)
}
