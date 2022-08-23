package snapshots

import (
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
)

const (
	STORAGE_PROVIDER_FILESYSTEM = "filesystem"
	STORAGE_PROVIDER_S3         = "s3"
	STORAGE_PROVIDER_GCS        = "gcs"
)

type API struct {
	connection *connection.Connection
}

func New(connection *connection.Connection) *API {
	return &API{connection}
}

// Creator creates snapshots creator builder
func (s *API) Creator() *SnapshotsCreator {
	return &SnapshotsCreator{
		helper: &snapshotsHelper{s.connection},
	}
}

// StatusCreateSnapshot creates the creator status getter builder
func (s *API) StatusCreateSnapshot() *SnapshotsCreateStatusGetter {
	return &SnapshotsCreateStatusGetter{
		helper: &snapshotsHelper{s.connection},
	}
}

// Restorer creates the restorer builder
func (s *API) Restorer() *SnapshotsRestorer {
	return &SnapshotsRestorer{
		helper: &snapshotsHelper{s.connection},
	}
}

// StatusRestoreSnapshot creates the restorer status getter builder
func (s *API) StatusRestoreSnapshot() *SnapshotsRestoreStatusGetter {
	return &SnapshotsRestoreStatusGetter{
		helper: &snapshotsHelper{s.connection},
	}
}
