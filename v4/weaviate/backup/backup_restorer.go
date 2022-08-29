package backup

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupRestorer struct {
	helper            *backupRestoreHelper
	className         string
	storageName       string
	backupID          string
	waitForCompletion bool
}

// WithClassName specifies the class which should be restored
func (r *BackupRestorer) WithClassName(className string) *BackupRestorer {
	r.className = className
	return r
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

func (r *BackupRestorer) Do(ctx context.Context) (*models.SnapshotRestoreMeta, error) {
	if r.waitForCompletion {
		return r.helper.restoreAndWaitForCompletion(ctx, r.className, r.storageName, r.backupID)

	}
	return r.helper.restore(ctx, r.className, r.storageName, r.backupID)
}
