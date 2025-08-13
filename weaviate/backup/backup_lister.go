package backup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type BackupLister struct {
	connection *connection.Connection
	backend    string
}

func (bc *BackupLister) WithBackend(backend string) *BackupLister {
	bc.backend = backend
	return bc
}

func (bc *BackupLister) Do(ctx context.Context) (models.BackupListResponse, error) {
	response, err := bc.connection.RunREST(ctx, bc.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == http.StatusOK {
		var obj models.BackupListResponse
		decodeErr := response.DecodeBodyIntoTarget(&obj)
		return obj, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(response)
}

func (bc *BackupLister) path() string {
	return fmt.Sprintf("/backups/%s", bc.backend)
}
