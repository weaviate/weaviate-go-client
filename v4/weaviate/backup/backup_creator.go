package backup

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupCreator struct {
	helper            *backupCreateHelper
	includeClasses    []string
	excludeClasses    []string
	storageName       string
	backupID          string
	waitForCompletion bool
}

func (c *BackupCreator) WithIncludeClasses(classes ...string) *BackupCreator {
	c.includeClasses = classes
	return c
}

func (c *BackupCreator) WithExcludeClasses(classes ...string) *BackupCreator {
	c.excludeClasses = classes
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

func (c *BackupCreator) Do(ctx context.Context) (*models.BackupCreateMeta, error) {
	if c.waitForCompletion {
		return c.helper.createAndWaitForCompletion(ctx, c.storageName, c.backupID, c.includeClasses, c.excludeClasses)
	}
	return c.helper.create(ctx, c.includeClasses, c.excludeClasses, c.storageName, c.backupID)
}
