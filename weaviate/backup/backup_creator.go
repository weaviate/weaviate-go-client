package backup

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

const waitTimeoutCreate = time.Second

type BackupCreator struct {
	connection        *connection.Connection
	statusGetter      *BackupCreateStatusGetter
	includeClasses    []string
	excludeClasses    []string
	backend           string
	backupID          string
	waitForCompletion bool
	config            *models.BackupConfig
}

func (c *BackupCreator) WithIncludeClassNames(classNames ...string) *BackupCreator {
	c.includeClasses = classNames
	return c
}

func (c *BackupCreator) WithExcludeClassNames(classNames ...string) *BackupCreator {
	c.excludeClasses = classNames
	return c
}

// WithBackend specifies the backend backup should be stored to
func (c *BackupCreator) WithBackend(backend string) *BackupCreator {
	c.backend = backend
	return c
}

// WithBackupID specifies unique id given to the backup
func (c *BackupCreator) WithBackupID(backupID string) *BackupCreator {
	c.backupID = backupID
	return c
}

// WithWaitForCompletion block until backup is created (succeeds or fails)
func (c *BackupCreator) WithWaitForCompletion(waitForCompletion bool) *BackupCreator {
	c.waitForCompletion = waitForCompletion
	return c
}

// WithConfig sets the compression configuration for the backup
func (c *BackupCreator) WithConfig(cfg *models.BackupConfig) *BackupCreator {
	c.config = cfg
	return c
}

func (c *BackupCreator) Do(ctx context.Context) (*models.BackupCreateResponse, error) {
	payload := models.BackupCreateRequest{
		ID:      c.backupID,
		Include: c.includeClasses,
		Exclude: c.excludeClasses,
		Config:  c.config,
	}

	if c.waitForCompletion {
		return c.createAndWaitForCompletion(ctx, payload)
	}
	return c.create(ctx, payload)
}

func (c *BackupCreator) create(ctx context.Context, payload models.BackupCreateRequest,
) (*models.BackupCreateResponse, error) {
	response, err := c.connection.RunREST(ctx, c.path(), http.MethodPost, payload)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == http.StatusOK {
		var obj models.BackupCreateResponse
		decodeErr := response.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(response)
}

func (c *BackupCreator) createAndWaitForCompletion(ctx context.Context, payload models.BackupCreateRequest,
) (*models.BackupCreateResponse, error) {
	response, err := c.create(ctx, payload)
	if err != nil {
		return nil, err
	}

	c.statusGetter.WithBackupID(c.backupID).WithBackend(c.backend)
	for {
		statusResponse, err := c.statusGetter.Do(ctx)
		if err != nil {
			return nil, err
		}
		switch *statusResponse.Status {
		case models.BackupCreateResponseStatusSUCCESS, models.BackupCreateResponseStatusFAILED:
			return c.merge(response, statusResponse), nil
		default:
			time.Sleep(waitTimeoutCreate)
		}
	}
}

func (c *BackupCreator) path() string {
	return fmt.Sprintf("/backups/%s", c.backend)
}

func (c *BackupCreator) merge(response *models.BackupCreateResponse,
	statusResponse *models.BackupCreateStatusResponse,
) *models.BackupCreateResponse {
	return &models.BackupCreateResponse{
		ID:      statusResponse.ID,
		Backend: statusResponse.Backend,
		Classes: response.Classes,
		Path:    statusResponse.Path,
		Status:  statusResponse.Status,
		Error:   statusResponse.Error,
	}
}
