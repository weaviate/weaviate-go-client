package backup

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

const (
	BACKEND_FILESYSTEM = "filesystem"
	BACKEND_S3         = "s3"
	BACKEND_GCS        = "gcs"
	BACKEND_AZURE      = "azure"
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
		connection:   s.connection,
		statusGetter: s.CreateStatusGetter(),
	}
}

// CreateStatusGetter creates create status getter builder
func (s *API) CreateStatusGetter() *BackupCreateStatusGetter {
	return &BackupCreateStatusGetter{
		connection: s.connection,
	}
}

// Restorer creates restorer builder
func (s *API) Restorer() *BackupRestorer {
	return &BackupRestorer{
		connection:   s.connection,
		statusGetter: s.RestoreStatusGetter(),
	}
}

// RestoreStatusGetter creates restore status getter builder
func (s *API) RestoreStatusGetter() *BackupRestoreStatusGetter {
	return &BackupRestoreStatusGetter{
		connection: s.connection,
	}
}

// Canceler creates a builder for "cancel backup" request.
func (s *API) Canceler() *BackupCanceler {
	return &BackupCanceler{
		connection: s.connection,
	}
}
