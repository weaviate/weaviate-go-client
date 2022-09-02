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

const waitTimeoutRestore = time.Second

type backupRestoreHelper struct {
	connection *connection.Connection
}

func (h *backupRestoreHelper) restore(ctx context.Context, className,
	storageName, backupID string,
) (*models.BackupRestoreMeta, error) {
	return h.restoreByEndpoint(ctx, endpointRestore(storageName))
}

func (h *backupRestoreHelper) restoreByEndpoint(ctx context.Context, endpoint string) (*models.BackupRestoreMeta, error) {
	return h.runREST(ctx, endpoint, http.MethodPost)
}

func (h *backupRestoreHelper) statusRestore(ctx context.Context,
	storageName, backupID string,
) (*models.BackupRestoreMeta, error) {
	return h.statusRestoreByEndpoint(ctx, endpointStatusRestore(storageName, backupID))
}

func (h *backupRestoreHelper) statusRestoreByEndpoint(ctx context.Context, endpoint string) (*models.BackupRestoreMeta, error) {
	return h.runREST(ctx, endpoint, http.MethodGet)
}

func (h *backupRestoreHelper) restoreAndWaitForCompletion(ctx context.Context,
	className, storageName, backupID string,
) (*models.BackupRestoreMeta, error) {
	endpoint := endpointRestore(storageName)
	if _, err := h.restoreByEndpoint(ctx, endpoint); err != nil {
		return nil, err
	}
	endpoint = endpointStatusRestore(storageName, backupID)
	for {
		meta, err := h.statusRestoreByEndpoint(ctx, endpoint)
		if err != nil {
			return nil, err
		}
		switch *meta.Status {
		case models.BackupRestoreMetaStatusSUCCESS, models.BackupRestoreMetaStatusFAILED:
			return meta, nil
		default:
			time.Sleep(waitTimeoutRestore)
		}
	}
}

func (h *backupRestoreHelper) runREST(ctx context.Context, endpoint, httpMethod string) (*models.BackupRestoreMeta, error) {
	responseData, err := h.connection.RunREST(ctx, endpoint, httpMethod, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var obj models.BackupRestoreMeta
		decodeErr := responseData.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewDerivedWeaviateClientError(err)
}

func endpointRestore(storageName string) string {
	// TODO change snapshots to backups
	return fmt.Sprintf("/backups/%s/restore", storageName)
}

func endpointStatusRestore(storageName, ID string) string {
	// TODO change snapshots to backups
	return fmt.Sprintf("/restore/%s/%s/restore", storageName, ID)
}
