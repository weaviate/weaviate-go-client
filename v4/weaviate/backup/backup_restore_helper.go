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
) (*models.SnapshotRestoreMeta, error) {
	return h.restoreByEndpoint(ctx, endpointRestore(className, storageName, backupID))
}

func (h *backupRestoreHelper) restoreByEndpoint(ctx context.Context, endpoint string) (*models.SnapshotRestoreMeta, error) {
	return h.runREST(ctx, endpoint, http.MethodPost)
}

func (h *backupRestoreHelper) statusRestore(ctx context.Context,
	className, storageName, backupID string,
) (*models.SnapshotRestoreMeta, error) {
	return h.statusRestoreByEndpoint(ctx, endpointRestore(className, storageName, backupID))
}

func (h *backupRestoreHelper) statusRestoreByEndpoint(ctx context.Context, endpoint string) (*models.SnapshotRestoreMeta, error) {
	return h.runREST(ctx, endpoint, http.MethodGet)
}

func (h *backupRestoreHelper) restoreAndWaitForCompletion(ctx context.Context,
	className, storageName, backupID string,
) (*models.SnapshotRestoreMeta, error) {
	endpoint := endpointRestore(className, storageName, backupID)
	if _, err := h.restoreByEndpoint(ctx, endpoint); err != nil {
		return nil, err
	}
	for {
		meta, err := h.statusRestoreByEndpoint(ctx, endpoint)
		if err != nil {
			return nil, err
		}
		switch *meta.Status {
		case models.SnapshotRestoreMetaStatusSUCCESS, models.SnapshotRestoreMetaStatusFAILED:
			return meta, nil
		default:
			time.Sleep(waitTimeoutRestore)
		}
	}
}

func (h *backupRestoreHelper) runREST(ctx context.Context, endpoint, httpMethod string) (*models.SnapshotRestoreMeta, error) {
	responseData, err := h.connection.RunREST(ctx, endpoint, httpMethod, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var obj models.SnapshotRestoreMeta
		decodeErr := responseData.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewDerivedWeaviateClientError(err)
}

func endpointRestore(className, storageName, backupID string) string {
	// TODO change snapshots to backups
	return fmt.Sprintf("/schema/%s/snapshots/%s/%s/restore", className, storageName, backupID)
}
