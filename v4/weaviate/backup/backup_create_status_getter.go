package backup

import (
	"context"

	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupCreateStatusGetter struct {
	helper      *backupCreateHelper
	className   string
	storageName string
	backupID    string
}

// WithClassName specifies the class which should be backed up
func (g *BackupCreateStatusGetter) WithClassName(className string) *BackupCreateStatusGetter {
	g.className = className
	return g
}

// WithStorageName specifies the storage where backup should be saved
func (g *BackupCreateStatusGetter) WithStorageName(storageName string) *BackupCreateStatusGetter {
	g.storageName = storageName
	return g
}

// WithBackupID specifies unique id given to the backup
func (g *BackupCreateStatusGetter) WithBackupID(backupID string) *BackupCreateStatusGetter {
	g.backupID = backupID
	return g
}

func (g *BackupCreateStatusGetter) Do(ctx context.Context) (*models.SnapshotMeta, error) {
	return g.helper.statusCreate(ctx, g.className, g.storageName, g.backupID)
}
