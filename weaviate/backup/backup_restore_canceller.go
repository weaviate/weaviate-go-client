package backup

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type BackupRestoreCanceler struct {
	connection *connection.Connection
	backend    string
	backupID   string
	bucket     string
	backupPath string
}

// WithBackend specifies the backup backend (e.g., s3, gcs, azure, filesystem)
func (rc *BackupRestoreCanceler) WithBackend(backend string) *BackupRestoreCanceler {
	rc.backend = backend
	return rc
}

// WithBackupID specifies the unique backup identifier
func (rc *BackupRestoreCanceler) WithBackupID(id string) *BackupRestoreCanceler {
	rc.backupID = id
	return rc
}

// WithBucket specifies the storage bucket (optional)
func (rc *BackupRestoreCanceler) WithBucket(bucket string) *BackupRestoreCanceler {
	rc.bucket = bucket
	return rc
}

// WithPath specifies the path within the bucket (optional)
func (rc *BackupRestoreCanceler) WithPath(path string) *BackupRestoreCanceler {
	rc.backupPath = path
	return rc
}

// Do cancels an ongoing backup restore operation
func (rc *BackupRestoreCanceler) Do(ctx context.Context) error {
	res, err := rc.connection.RunREST(ctx, rc.path(), http.MethodDelete, nil)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusNoContent { // 204 - Successfully cancelled
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (rc *BackupRestoreCanceler) path() string {
	basePath := fmt.Sprintf("/backups/%s/%s/restore", rc.backend, rc.backupID)

	params := url.Values{}
	if rc.bucket != "" {
		params.Set("bucket", rc.bucket)
	}
	if rc.backupPath != "" {
		params.Set("path", rc.backupPath)
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", basePath, params.Encode())
	}
	return basePath
}
