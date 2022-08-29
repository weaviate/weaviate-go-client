package backup

import (
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
)

const (
	STORAGE_FILESYSTEM = "filesystem"
	STORAGE_S3         = "s3"
	STORAGE_GCS        = "gcs"
)

type API struct {
	connection *connection.Connection
}

func New(connection *connection.Connection) *API {
	return &API{connection}
}

// Creator creates backup creator builder
func (s *API) Creator() *BackupCreator {
	return &BackupCreator{
		helper: &backupCreateHelper{s.connection},
	}
}

// CreateStatusGetter creates create status getter builder
func (s *API) CreateStatusGetter() *BackupCreateStatusGetter {
	return &BackupCreateStatusGetter{
		helper: &backupCreateHelper{s.connection},
	}
}

// Restorer creates restorer builder
func (s *API) Restorer() *BackupRestorer {
	return &BackupRestorer{
		helper: &backupRestoreHelper{s.connection},
	}
}

// RestoreStatusGetter creates restore status getter builder
func (s *API) RestoreStatusGetter() *BackupRestoreStatusGetter {
	return &BackupRestoreStatusGetter{
		helper: &backupRestoreHelper{s.connection},
	}
}
