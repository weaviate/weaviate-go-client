package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type BackupCanceler struct {
	connection *connection.Connection
	backend    string
	backupID   string
}

func (bc *BackupCanceler) WithBackend(backend string) *BackupCanceler {
	bc.backend = backend
	return bc
}

func (bc *BackupCanceler) WithBackupID(id string) *BackupCanceler {
	bc.backupID = id
	return bc
}

func (bc *BackupCanceler) Do(ctx context.Context) error {
	res, err := bc.connection.RunREST(ctx, bc.path(), http.MethodDelete, nil)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusNoContent { // 204 - Successfully deleted
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (bc BackupCanceler) path() string {
	return fmt.Sprintf("/backups/%s/%s", bc.backend, bc.backupID)
}
