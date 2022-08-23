package snapshots

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type SnapshotsCreator struct {
	helper                *snapshotsHelper
	className             string
	storageProvider       string
	snapshotID            string
	withWaitForCompletion bool
}

// WithClassName specifies the class which should be backed up
func (s *SnapshotsCreator) WithClassName(className string) *SnapshotsCreator {
	s.className = className
	return s
}

// WithStorageProvider specifies the class which should be backed up
func (s *SnapshotsCreator) WithStorageProvider(storageProvider string) *SnapshotsCreator {
	s.storageProvider = storageProvider
	return s
}

// WithSnapshotID specifies the class which should be backed up
func (s *SnapshotsCreator) WithSnapshotID(snapshotID string) *SnapshotsCreator {
	s.snapshotID = snapshotID
	return s
}

// WithWaitForCompletion block while snapshost is being created (until it succeeds or fails)
func (s *SnapshotsCreator) WithWaitForCompletion() *SnapshotsCreator {
	s.withWaitForCompletion = true
	return s
}

func (s *SnapshotsCreator) Do(ctx context.Context) (*models.SnapshotMeta, error) {
	if !s.withWaitForCompletion {
		return s.helper.createSnapshot(ctx, s.className, s.storageProvider, s.snapshotID)
	}
	return s.helper.createAndWaitForCompletion(ctx, s.className, s.storageProvider, s.snapshotID)
}
