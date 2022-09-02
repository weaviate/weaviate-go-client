package backup

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupRestoreStatusGetter struct {
	helper      *backupRestoreHelper
	storageName string
	backupID    string
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

func (g *BackupRestoreStatusGetter) Do(ctx context.Context) (*models.BackupRestoreMeta, error) {
	return g.helper.statusRestore(ctx, g.storageName, g.backupID)
}
