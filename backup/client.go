package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewClient(t internal.Transport) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{transport: t}
}

type Client struct {
	transport internal.Transport
}

type Info struct {
	Backend string // Backup storage backend.
	ID      string // Backup ID.
	Bucket  string // Dedicated bucket name.
	Path    string // Path to backup in the backend storage.
	Error   string // Backup creation / restoration error.
	Status  Status // Backup creation / restoration status.

	StartedAt   *time.Time // Time at which the backup creation started.
	CompletedAt *time.Time // Time at which the backup was completed, successfully or otherwise.

	// IncludesCollections is always empty for backups that are still being created.
	// This field will be populated for restored backups and already-created backups
	// returned from List, provided these backups included at least one collection.
	IncludesCollections []string // Collections included in the backup

	// Backup size in GiB. Similarly to IncludesCollections,
	// this value only exists for completed backups.
	SizeGiB *float32

	// BaseBackupID is the ID of the backup this incremental backup was built on.
	// This field is empty if the backup is not incremental.
	BaseBackupID string

	c         *Client
	operation api.BackupOperation
}

// IsCompleted returns true if the backup operation has completed, successfully or otherwise.
// All backups returned from List are completed by definition.
func (i *Info) IsCompleted() bool {
	return i.operation == completed ||
		i.Status == StatusSuccess ||
		i.Status == StatusFailed ||
		i.Status == StatusCanceled
}

type (
	CompressionLevel api.BackupCompressionLevel
	Status           api.BackupStatus
	RBACRestore      api.RBACRestoreOption
)

const (
	CompressionLevelDefault             CompressionLevel = CompressionLevel(api.BackupCompressionLevelDefault)
	CompressionLevelBestSpeed           CompressionLevel = CompressionLevel(api.BackupCompressionLevelBestSpeed)
	CompressionLevelBestCompression     CompressionLevel = CompressionLevel(api.BackupCompressionLevelBestCompression)
	CompressionLevelZstdDefault         CompressionLevel = CompressionLevel(api.BackupCompressionLevelZstdDefault)
	CompressionLevelZstdBestSpeed       CompressionLevel = CompressionLevel(api.BackupCompressionLevelZstdBestSpeed)
	CompressionLevelZstdBestCompression CompressionLevel = CompressionLevel(api.BackupCompressionLevelZstdBestCompression)
	CompressionLevelNone                CompressionLevel = CompressionLevel(api.BackupCompressionLevelNone)

	StatusStarted      Status = Status(api.BackupStatusStarted)
	StatusTransferring Status = Status(api.BackupStatusTransferring)
	StatusTransferred  Status = Status(api.BackupStatusTransferred)
	StatusFinalizing   Status = Status(api.BackupStatusFinalizing)
	StatusCanceling    Status = Status(api.BackupStatusCanceling)
	StatusCanceled     Status = Status(api.BackupStatusCanceled)
	StatusSuccess      Status = Status(api.BackupStatusSuccess)
	StatusFailed       Status = Status(api.BackupStatusFailed)

	RBACRestoreAll  RBACRestore = RBACRestore(api.RBACRestoreAll)
	RBACRestoreNone RBACRestore = RBACRestore(api.RBACRestoreNone)
)

type CreateOptions struct {
	Backend            string           // Required: backend storage.
	ID                 string           // Required: backup ID.
	Path               string           // Path to backup in the backend storage.
	Endpoint           string           // Name of the endpoint, e.g. s3.amazonaws.com
	Bucket             string           // Dedicated bucket name.
	IncludeCollections []string         // Collections to be included in the backup.
	ExcludeCollections []string         // Collections to be excluded from the backup.
	BaseBackupID       string           // BaseBackupID for an incremental backup.
	MaxCPUPercentage   int              // Maximum %CPU utilization.
	ChunkSizeMiB       int              // Target chunk size in MiB.
	CompressionLevel   CompressionLevel // Hint for selecting the optimal compression algorithm.
}

/** Create a new backup.*/
func (c *Client) Create(ctx context.Context, options CreateOptions) (*Info, error) {
	req := &api.CreateBackupRequest{
		Backend:            options.Backend,
		ID:                 options.ID,
		Bucket:             options.Bucket,
		BackupPath:         options.Path,
		Endpoint:           options.Endpoint,
		IncludeCollections: options.IncludeCollections,
		ExcludeCollections: options.ExcludeCollections,
		BaseBackupID:       options.BaseBackupID,
		MaxCPUPercentage:   options.MaxCPUPercentage,
		ChunkSizeMiB:       options.ChunkSizeMiB,
		CompressionLevel:   api.BackupCompressionLevel(options.CompressionLevel),
	}

	var resp api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("create backup: %w", err)
	}

	info := infoFromAPI(&resp, c, api.BackupOperationCreate)
	return &info, nil
}

type RestoreOptions struct {
	Backend            string            // Required: backend storage.
	ID                 string            // Required: backup ID.
	Path               string            // Path to backup in the backend storage.
	Endpoint           string            // Name of the endpoint, e.g. s3.amazonaws.com
	Bucket             string            // Dedicated bucket name.
	IncludeCollections []string          // Collections to be included in the backup.
	ExcludeCollections []string          // Collections to be excluded from the backup.
	MaxCPUPercentage   int               // Maximum %CPU utilization.
	OverwriteAlias     bool              // Allow overwriting aliases.
	RestoreUsers       RBACRestore       // Select strategy for restoring RBAC users.
	RestoreRoles       RBACRestore       // Select strategy for restoring RBAC roles.
	NodeMapping        map[string]string // Remap node names stored in the backup.
}

func (c *Client) Restore(ctx context.Context, options RestoreOptions) (*Info, error) {
	req := &api.RestoreBackupRequest{
		Backend:            options.Backend,
		ID:                 options.ID,
		Bucket:             options.Bucket,
		BackupPath:         options.Path,
		Endpoint:           options.Endpoint,
		IncludeCollections: options.IncludeCollections,
		ExcludeCollections: options.ExcludeCollections,
		MaxCPUPercentage:   options.MaxCPUPercentage,
		OverwriteAlias:     options.OverwriteAlias,
		RestoreUsers:       api.RBACRestoreOption(options.RestoreUsers),
		RestoreRoles:       api.RBACRestoreOption(options.RestoreRoles),
		NodeMapping:        options.NodeMapping,
	}

	var resp api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("restore backup: %w", err)
	}

	info := infoFromAPI(&resp, c, api.BackupOperationRestore)
	return &info, nil
}

type GetStatusOptions struct {
	Backend string // Required: Backend storage.
	ID      string // Required: Backup ID.
}

func (c *Client) GetCreateStatus(ctx context.Context, options GetStatusOptions) (*Info, error) {
	return c.getStatus(ctx, options, api.BackupOperationCreate)
}

func (c *Client) GetRestoreStatus(ctx context.Context, options GetStatusOptions) (*Info, error) {
	return c.getStatus(ctx, options, api.BackupOperationRestore)
}

func (c *Client) getStatus(ctx context.Context, options GetStatusOptions, operation api.BackupOperation) (*Info, error) {
	req := &api.BackupStatusRequest{
		Backend:   options.Backend,
		ID:        options.ID,
		Operation: operation,
	}

	var resp api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("get backup status: %w", err)
	}

	info := infoFromAPI(&resp, c, operation)
	return &info, nil
}

type ListOptions struct {
	Backend         string // Required: Backend storage.
	StartingTimeAsc bool   // Set to true to order backups by their StartedAt time in ascending order.
}

func (c *Client) List(ctx context.Context, options ListOptions) ([]Info, error) {
	req := &api.ListBackupsRequest{
		Backend:         options.Backend,
		StartingTimeAsc: options.StartingTimeAsc,
	}

	var resp []api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	infos := make([]Info, len(resp))
	for i, bak := range resp {
		infos[i] = infoFromAPI(&bak, c, completed)
	}
	return infos, nil
}

type CancelOptions struct {
	Backend string // Required: Backend storage.
	ID      string // Required: Backup ID.
}

// Cancel an in-progress backup creation.
func (c *Client) CancelCreate(ctx context.Context, options CancelOptions) error {
	return c.cancel(ctx, options, api.BackupOperationCreate)
}

// Cancel an in-progress backup restoration.
func (c *Client) CancelRestore(ctx context.Context, options CancelOptions) error {
	return c.cancel(ctx, options, api.BackupOperationRestore)
}

func (c *Client) cancel(ctx context.Context, options CancelOptions, op api.BackupOperation) error {
	req := api.CancelBackupRequest{
		Backend:   options.Backend,
		ID:        options.ID,
		Operation: op,
	}
	if err := c.transport.Do(ctx, &req, nil); err != nil {
		return fmt.Errorf("cancel backup: %w", err)
	}
	return nil
}

// completed is a special flag for operations returned from List.
// Those operations technically have api.CreateBackup origin, but
// Info.operation is used to await backup completion (or any other status),
// and, since they are already completed, we can give awaiters a hint.
const completed api.BackupOperation = api.BackupOperation(api.BackupOperationCreate - 1)

func infoFromAPI(bak *api.BackupInfo, c *Client, op api.BackupOperation) Info {
	return Info{
		Backend:             bak.Backend,
		ID:                  bak.ID,
		Path:                bak.Path,
		Bucket:              bak.Bucket,
		BaseBackupID:        bak.BaseBackupID,
		Status:              Status(bak.Status),
		Error:               bak.Error,
		StartedAt:           bak.StartedAt,
		CompletedAt:         bak.CompletedAt,
		IncludesCollections: bak.IncludesCollections,
		SizeGiB:             bak.SizeGiB,

		operation: op,
		c:         c,
	}
}
