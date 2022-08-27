package snapshots

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type SnapshotsCreateStatusGetter struct {
	helper          *snapshotsHelper
	className       string
	storageProvider string
	snapshotID      string
}

// WithClassName specifies the class which should be backed up
func (s *SnapshotsCreateStatusGetter) WithClassName(className string) *SnapshotsCreateStatusGetter {
	s.className = className
	return s
}

// WithStorageProvider specifies the class which should be backed up
func (s *SnapshotsCreateStatusGetter) WithStorageProvider(storageProvider string) *SnapshotsCreateStatusGetter {
	s.storageProvider = storageProvider
	return s
}

// WithSnapshotID specifies the class which should be backed up
func (s *SnapshotsCreateStatusGetter) WithSnapshotID(snapshotID string) *SnapshotsCreateStatusGetter {
	s.snapshotID = snapshotID
	return s
}

func (s *SnapshotsCreateStatusGetter) Do(ctx context.Context) (*models.SnapshotMeta, error) {
	return s.helper.statusCreateSnapshot(ctx, s.className, s.storageProvider, s.snapshotID)
}
