package backup

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupRestoreStatusGetter struct {
	helper      *backupRestoreHelper
	className   string
	storageName string
	backupID    string
}

// WithClassName specifies the class which should be restored
func (g *BackupRestoreStatusGetter) WithClassName(className string) *BackupRestoreStatusGetter {
	g.className = className
	return g
}

// WithStorageName specifies the storage from backup should be restored
func (g *BackupRestoreStatusGetter) WithStorageName(storageName string) *BackupRestoreStatusGetter {
	g.storageName = storageName
	return g
}

// WithBackupID specifies unique id given to the backup
func (g *BackupRestoreStatusGetter) WithBackupID(backupID string) *BackupRestoreStatusGetter {
	g.backupID = backupID
	return g
}

func (g *BackupRestoreStatusGetter) Do(ctx context.Context) (*models.SnapshotRestoreMeta, error) {
	return g.helper.statusRestore(ctx, g.className, g.storageName, g.backupID)
}
