package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/except"
	"github.com/semi-technologies/weaviate/entities/models"
)

type BackupCreateStatusGetter struct {
	connection *connection.Connection
	backend    string
	backupID   string
}

// WithBackend specifies the backend backup is stored to
func (g *BackupCreateStatusGetter) WithBackend(backend string) *BackupCreateStatusGetter {
	g.backend = backend
	return g
}

// WithBackupID specifies unique id given to the backup
func (g *BackupCreateStatusGetter) WithBackupID(backupID string) *BackupCreateStatusGetter {
	g.backupID = backupID
	return g
}

func (g *BackupCreateStatusGetter) Do(ctx context.Context) (*models.BackupCreateStatusResponse, error) {
	responseData, err := g.connection.RunREST(ctx, g.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var obj models.BackupCreateStatusResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewDerivedWeaviateClientError(err)
}

func (g *BackupCreateStatusGetter) path() string {
	return fmt.Sprintf("/backups/%s/%s", g.backend, g.backupID)
}
