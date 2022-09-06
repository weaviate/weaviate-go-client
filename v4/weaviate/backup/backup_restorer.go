package backup

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupRestorer struct {
	helper            *backupRestoreHelper
	includeClasses    []string
	excludeClasses    []string
	storageName       string
	backupID          string
	waitForCompletion bool
}

func (c *BackupRestorer) WithIncludeClassNames(classes ...string) *BackupRestorer {
	c.includeClasses = classes
	return c
}

func (c *BackupRestorer) WithExcludeClassNames(classes ...string) *BackupRestorer {
	c.excludeClasses = classes
	return c
}

// WithStorageName specifies the storage from backup should be restored
func (r *BackupRestorer) WithStorageName(storageName string) *BackupRestorer {
	r.storageName = storageName
	return r
}

// WithBackupID specifies unique id given to the backup
func (r *BackupRestorer) WithBackupID(backupID string) *BackupRestorer {
	r.backupID = backupID
	return r
}

// WithWaitForCompletion block until backup is restored (succeeds or fails)
func (r *BackupRestorer) WithWaitForCompletion(waitForCompletion bool) *BackupRestorer {
	r.waitForCompletion = waitForCompletion
	return r
}

func (r *BackupRestorer) Do(ctx context.Context) (*models.BackupRestoreMeta, error) {
	if r.waitForCompletion {
		return r.helper.restoreAndWaitForCompletion(ctx, r.includeClasses, r.excludeClasses, r.storageName, r.backupID)
	}
	return r.helper.restore(ctx, r.storageName, r.backupID, r.includeClasses, r.excludeClasses)
}
