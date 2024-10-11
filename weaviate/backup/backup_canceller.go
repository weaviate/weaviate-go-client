package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
)

type BackupCanceller struct {
	connection *connection.Connection
	backend    string
	backupID   string
}

func (bc *BackupCanceller) WithBackend(backend string) *BackupCanceller {
	bc.backend = backend
	return bc
}

func (bc *BackupCanceller) WithBackupID(id string) *BackupCanceller {
	bc.backupID = id
	return bc
}

func (bc *BackupCanceller) Do(ctx context.Context) error {
	res, err := bc.connection.RunREST(ctx, bc.path(), http.MethodDelete, nil)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusNoContent { // 204 - Successfully deleted
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (bc BackupCanceller) path() string {
	return fmt.Sprintf("/backups/%s/%s", bc.backend, bc.backupID)
}
