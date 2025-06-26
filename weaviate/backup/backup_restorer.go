package backup

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

const waitTimeoutRestore = time.Second

type BackupRestorer struct {
	connection        *connection.Connection
	statusGetter      *BackupRestoreStatusGetter
	includeClasses    []string
	excludeClasses    []string
	backend           string
	backupID          string
	waitForCompletion bool
	config            *models.RestoreConfig
	rbacConfig        *rbac.RBACConfig
}

func (c *BackupRestorer) WithIncludeClassNames(classNames ...string) *BackupRestorer {
	c.includeClasses = classNames
	return c
}

func (c *BackupRestorer) WithExcludeClassNames(classNames ...string) *BackupRestorer {
	c.excludeClasses = classNames
	return c
}

// WithBackend specifies the backend backup should be restored from
func (r *BackupRestorer) WithBackend(backend string) *BackupRestorer {
	r.backend = backend
	return r
}

// WithBackupID specifies unique id given to the backup
func (r *BackupRestorer) WithBackupID(backupID string) *BackupRestorer {
	r.backupID = backupID
	return r
}

// WithWaitForCompletion block until backup is restored (succeeds or fails)
func (r *BackupRestorer) WithWaitForCompletion(waitForCompletion bool) *BackupRestorer {
	r.waitForCompletion = waitForCompletion
	return r
}

func (r *BackupRestorer) WithConfig(cfg *models.RestoreConfig) *BackupRestorer {
	r.config = cfg
	return r
}

// WithRBAC sets the RBAC configuration for the restore
func (r *BackupRestorer) WithRBAC(rbacConfig *rbac.RBACConfig) *BackupRestorer {
	r.rbacConfig = rbacConfig
	return r
}

// WithRBACScope sets the RBAC scope for the restore (convenience method)
func (r *BackupRestorer) WithRBACScope(scope rbac.RBACScope) *BackupRestorer {
	if r.rbacConfig == nil {
		r.rbacConfig = &rbac.RBACConfig{}
	}
	r.rbacConfig.Scope = scope
	return r
}

// WithRoles sets which roles to include in the restore
func (r *BackupRestorer) WithRoles(selection rbac.RoleSelection, specificRoles ...string) *BackupRestorer {
	if r.rbacConfig == nil {
		r.rbacConfig = &rbac.RBACConfig{}
	}
	r.rbacConfig.RoleSelection = selection
	if selection == rbac.RoleSelectionSpecific {
		r.rbacConfig.SpecificRoles = specificRoles
	}
	return r
}

func (r *BackupRestorer) Do(ctx context.Context) (*models.BackupRestoreResponse, error) {
	payload := models.BackupRestoreRequest{
		Include: r.includeClasses,
		Exclude: r.excludeClasses,
		Config:  r.config,
	}
	
	// Note: RBAC configuration is passed as a custom extension
	// This would be handled via query parameters or custom headers
	// depending on the Weaviate server implementation

	if r.waitForCompletion {
		return r.restoreAndWaitForCompletion(ctx, payload)
	}
	return r.restore(ctx, payload)
}

func (r *BackupRestorer) restore(ctx context.Context, payload models.BackupRestoreRequest,
) (*models.BackupRestoreResponse, error) {
	response, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, payload)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == http.StatusOK {
		var obj models.BackupRestoreResponse
		decodeErr := response.DecodeBodyIntoTarget(&obj)
		return &obj, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(response)
}

func (r *BackupRestorer) restoreAndWaitForCompletion(ctx context.Context, payload models.BackupRestoreRequest,
) (*models.BackupRestoreResponse, error) {
	response, err := r.restore(ctx, payload)
	if err != nil {
		return nil, err
	}

	r.statusGetter.WithBackupID(r.backupID).WithBackend(r.backend)
	for {
		statusResponse, err := r.statusGetter.Do(ctx)
		if err != nil {
			return nil, err
		}
		switch *statusResponse.Status {
		case models.BackupRestoreResponseStatusSUCCESS, models.BackupRestoreResponseStatusFAILED:
			return r.merge(response, statusResponse), nil
		default:
			time.Sleep(waitTimeoutRestore)
		}
	}
}

func (r *BackupRestorer) path() string {
	return fmt.Sprintf("/backups/%s/%s/restore", r.backend, r.backupID)
}

func (r *BackupRestorer) merge(response *models.BackupRestoreResponse,
	statusResponse *models.BackupRestoreStatusResponse,
) *models.BackupRestoreResponse {
	return &models.BackupRestoreResponse{
		ID:      statusResponse.ID,
		Backend: statusResponse.Backend,
		Classes: response.Classes,
		Path:    statusResponse.Path,
		Status:  statusResponse.Status,
		Error:   statusResponse.Error,
	}
}
