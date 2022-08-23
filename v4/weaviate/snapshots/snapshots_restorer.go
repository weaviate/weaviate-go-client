package snapshots

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type SnapshotsRestorer struct {
	helper                *snapshotsHelper
	className             string
	storageProvider       string
	snapshotID            string
	withWaitForCompletion bool
}

// WithClassName specifies the class which should be backed up
func (s *SnapshotsRestorer) WithClassName(className string) *SnapshotsRestorer {
	s.className = className
	return s
}

// WithStorageProvider specifies the class which should be backed up
func (s *SnapshotsRestorer) WithStorageProvider(storageProvider string) *SnapshotsRestorer {
	s.storageProvider = storageProvider
	return s
}

// WithSnapshotID specifies the class which should be backed up
func (s *SnapshotsRestorer) WithSnapshotID(snapshotID string) *SnapshotsRestorer {
	s.snapshotID = snapshotID
	return s
}

// WithWaitForCompletion block while snapshost is being created (until it succeeds or fails)
func (s *SnapshotsRestorer) WithWaitForCompletion() *SnapshotsRestorer {
	s.withWaitForCompletion = true
	return s
}

func (s *SnapshotsRestorer) Do(ctx context.Context) (*models.SnapshotRestoreMeta, error) {
	if !s.withWaitForCompletion {
		return s.helper.restoreSnapshot(ctx, s.className, s.storageProvider, s.snapshotID)
	}
	return s.helper.restoreAndWaitForCompletion(ctx, s.className, s.storageProvider, s.snapshotID)
}
