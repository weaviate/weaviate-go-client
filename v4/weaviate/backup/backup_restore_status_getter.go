package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupRestoreStatusGetter struct {
	connection *connection.Connection
	backend    string
	backupID   string
}

// WithBackend specifies the backend backup is restored from
func (g *BackupRestoreStatusGetter) WithBackend(backend string) *BackupRestoreStatusGetter {
	g.backend = backend
	return g
}

// WithBackupID specifies unique id given to the backup
func (g *BackupRestoreStatusGetter) WithBackupID(backupID string) *BackupRestoreStatusGetter {
	g.backupID = backupID
	return g
}

func (g *BackupRestoreStatusGetter) Do(ctx context.Context) (*models.BackupRestoreStatusResponse, error) {
	responseData, err := g.connection.RunREST(ctx, g.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var obj models.BackupRestoreStatusResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewDerivedWeaviateClientError(err)
}

func (g *BackupRestoreStatusGetter) path() string {
	return fmt.Sprintf("/backups/%s/%s/restore", g.backend, g.backupID)
}
