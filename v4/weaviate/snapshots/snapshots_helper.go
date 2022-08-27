package snapshots

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type snapshotsHelper struct {
	connection *connection.Connection
}

func (s *snapshotsHelper) createSnapshot(ctx context.Context,
	className, storageProvider, snapshotID string,
) (*models.SnapshotMeta, error) {
	url := fmt.Sprintf("/schema/%s/snapshots/%s/%s", className, storageProvider, snapshotID)
	resp, err := s.runREST(ctx, url, http.MethodPost)
	if err != nil {
		return nil, err
	}
	return s.decodeCreateResponse(resp)
}

func (s *snapshotsHelper) createAndWaitForCompletion(ctx context.Context,
	className, storageProvider, snapshotID string,
) (*models.SnapshotMeta, error) {
	meta, err := s.createSnapshot(ctx, className, storageProvider, snapshotID)
	if err != nil {
		return nil, err
	}
	for {
		meta, err = s.statusCreateSnapshot(ctx, className, storageProvider, snapshotID)
		if err != nil {
			return nil, err
		}
		switch *meta.Status {
		case models.SnapshotMetaStatusSUCCESS, models.SnapshotMetaStatusFAILED:
			return meta, nil
		default:
			time.Sleep(2.0 * time.Second)
		}
	}
}

func (s *snapshotsHelper) statusCreateSnapshot(ctx context.Context,
	className, storageProvider, snapshotID string,
) (*models.SnapshotMeta, error) {
	url := fmt.Sprintf("/schema/%s/snapshots/%s/%s", className, storageProvider, snapshotID)
	resp, err := s.runREST(ctx, url, http.MethodGet)
	if err != nil {
		return nil, err
	}
	return s.decodeCreateResponse(resp)
}

func (s *snapshotsHelper) decodeCreateResponse(
	responseData *connection.ResponseData,
) (*models.SnapshotMeta, error) {
	var obj models.SnapshotMeta
	decodeErr := responseData.DecodeBodyIntoTarget(&obj)
	return &obj, decodeErr
}

func (s *snapshotsHelper) runREST(ctx context.Context, url, httpMethod string) (*connection.ResponseData, error) {
	responseData, err := s.connection.RunREST(ctx, url, httpMethod, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		return responseData, nil
	}
	return nil, except.NewDerivedWeaviateClientError(err)
}

func (s *snapshotsHelper) restoreSnapshot(ctx context.Context,
	className, storageProvider, snapshotID string,
) (*models.SnapshotRestoreMeta, error) {
	url := fmt.Sprintf("/schema/%s/snapshots/%s/%s/restore", className, storageProvider, snapshotID)
	resp, err := s.runREST(ctx, url, http.MethodPost)
	if err != nil {
		return nil, err
	}
	return s.decodeRestoreResponse(resp)
}

func (s *snapshotsHelper) statusRestoreSnapshot(ctx context.Context,
	className, storageProvider, snapshotID string,
) (*models.SnapshotRestoreMeta, error) {
	url := fmt.Sprintf("/schema/%s/snapshots/%s/%s/restore", className, storageProvider, snapshotID)
	resp, err := s.runREST(ctx, url, http.MethodGet)
	if err != nil {
		return nil, err
	}
	return s.decodeRestoreResponse(resp)
}

func (s *snapshotsHelper) restoreAndWaitForCompletion(ctx context.Context,
	className, storageProvider, snapshotID string,
) (*models.SnapshotRestoreMeta, error) {
	meta, err := s.restoreSnapshot(ctx, className, storageProvider, snapshotID)
	if err != nil {
		return nil, err
	}
	for {
		meta, err = s.statusRestoreSnapshot(ctx, className, storageProvider, snapshotID)
		if err != nil {
			return nil, err
		}
		switch *meta.Status {
		case models.SnapshotRestoreMetaStatusSUCCESS, models.SnapshotRestoreMetaStatusFAILED:
			return meta, nil
		default:
			time.Sleep(2.0 * time.Second)
		}
	}
}

func (s *snapshotsHelper) decodeRestoreResponse(
	responseData *connection.ResponseData,
) (*models.SnapshotRestoreMeta, error) {
	var obj models.SnapshotRestoreMeta
	decodeErr := responseData.DecodeBodyIntoTarget(&obj)
	return &obj, decodeErr
}
