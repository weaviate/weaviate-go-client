package backup

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupCreator struct {
	helper            *backupCreateHelper
	className         string
	storageName       string
	backupID          string
	waitForCompletion bool
}

// WithClassName specifies the class which should be backed up
func (c *BackupCreator) WithClassName(className string) *BackupCreator {
	c.className = className
	return c
}

// WithStorageName specifies the storage where backup should be saved
func (c *BackupCreator) WithStorageName(storageName string) *BackupCreator {
	c.storageName = storageName
	return c
}

// WithBackupID specifies unique id given to the backup
func (c *BackupCreator) WithBackupID(backupID string) *BackupCreator {
	c.backupID = backupID
	return c
}

// WithWaitForCompletion block until backup is created (succeeds or fails)
func (c *BackupCreator) WithWaitForCompletion(waitForCompletion bool) *BackupCreator {
	c.waitForCompletion = waitForCompletion
	return c
}

func (c *BackupCreator) Do(ctx context.Context) (*models.SnapshotMeta, error) {
	if c.waitForCompletion {
		return c.helper.createAndWaitForCompletion(ctx, c.className, c.storageName, c.backupID)

	}
	return c.helper.create(ctx, c.className, c.storageName, c.backupID)
}
