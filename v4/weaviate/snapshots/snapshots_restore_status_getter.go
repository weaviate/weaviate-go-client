package snapshots

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type SnapshotsRestoreStatusGetter struct {
	helper          *snapshotsHelper
	className       string
	storageProvider string
	snapshotID      string
}

// WithClassName specifies the class which should be backed up
func (s *SnapshotsRestoreStatusGetter) WithClassName(className string) *SnapshotsRestoreStatusGetter {
	s.className = className
	return s
}

// WithStorageProvider specifies the class which should be backed up
func (s *SnapshotsRestoreStatusGetter) WithStorageProvider(storageProvider string) *SnapshotsRestoreStatusGetter {
	s.storageProvider = storageProvider
	return s
}

// WithSnapshotID specifies the class which should be backed up
func (s *SnapshotsRestoreStatusGetter) WithSnapshotID(snapshotID string) *SnapshotsRestoreStatusGetter {
	s.snapshotID = snapshotID
	return s
}

func (s *SnapshotsRestoreStatusGetter) Do(ctx context.Context) (*models.SnapshotRestoreMeta, error) {
	return s.helper.restoreGetSnapshot(ctx, s.className, s.storageProvider, s.snapshotID)
}
