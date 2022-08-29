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

func (h *backupCreateHelper) create(ctx context.Context, className,
	storageName, backupID string,
) (*models.SnapshotMeta, error) {
	return h.createByEndpoint(ctx, endpointCreate(className, storageName, backupID))
}

func (h *backupCreateHelper) createByEndpoint(ctx context.Context, endpoint string) (*models.SnapshotMeta, error) {
	return h.runREST(ctx, endpoint, http.MethodPost)
}

func (h *backupCreateHelper) statusCreate(ctx context.Context,
	className, storageName, backupID string,
) (*models.SnapshotMeta, error) {
	return h.statusCreateByEndpoint(ctx, endpointCreate(className, storageName, backupID))
}

func (h *backupCreateHelper) statusCreateByEndpoint(ctx context.Context, endpoint string) (*models.SnapshotMeta, error) {
	return h.runREST(ctx, endpoint, http.MethodGet)
}

func (h *backupCreateHelper) createAndWaitForCompletion(ctx context.Context,
	className, storageName, backupID string,
) (*models.SnapshotMeta, error) {
	endpoint := endpointCreate(className, storageName, backupID)
	if _, err := h.createByEndpoint(ctx, endpoint); err != nil {
		return nil, err
	}
	for {
		meta, err := h.statusCreateByEndpoint(ctx, endpoint)
		if err != nil {
			return nil, err
		}
		switch *meta.Status {
		case models.SnapshotMetaStatusSUCCESS, models.SnapshotMetaStatusFAILED:
			return meta, nil
		default:
			time.Sleep(waitTimeoutCreate)
		}
	}
}

func (h *backupCreateHelper) runREST(ctx context.Context, endpoint, httpMethod string) (*models.SnapshotMeta, error) {
	responseData, err := h.connection.RunREST(ctx, endpoint, httpMethod, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var obj models.SnapshotMeta
		decodeErr := responseData.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewDerivedWeaviateClientError(err)
}

func endpointCreate(className, storageName, backupID string) string {
	// TODO change snapshots to backups
	return fmt.Sprintf("/schema/%s/snapshots/%s/%s", className, storageName, backupID)
}
