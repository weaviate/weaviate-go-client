package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

func NewClient(t internal.Transport) *Client {
	return &Client{transport: t}
}

type Client struct {
	transport internal.Transport
}

type Info struct {
	Backend     string    // Backup storage backend
	ID          string    // Backup ID
	Path        string    // Path to backup in the backend storage
	Error       string    // Backup creation / restoration error.
	Status      Status    // Backup creation / restoration status.
	StartedAt   time.Time // Time at which the backup creation started.
	CompletedAt time.Time // Time at which the backup was completed, successfully or otherwise.

	// IncludesCollections is always empty for backups that are still being created.
	// This field will be populated for restored backups and already-created backups
	// returned from List, provided these backups included at least one collection.
	IncludesCollections []string // Collections included in the backup

	// Backup size in GiB. Similarly to IncludesCollections,
	// this value only exists for completed backups.
	SizeGiB *float32

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
	StatusSuccess      Status = Status(api.BackupStatusSuccess)
	StatusFailed       Status = Status(api.BackupStatusFailed)
	StatusCanceled     Status = Status(api.BackupStatusCanceled)

	RBACRestoreAll  RBACRestore = RBACRestore(api.RBACRestoreAll)
	RBACRestoreNone RBACRestore = RBACRestore(api.RBACRestoreNone)
)

type Create struct {
	Backend            string           // Required: backend storage.
	ID                 string           // Required: backup ID.
	Path               string           // Path to backup in the backend storage.
	Endpoint           string           // Name of the endpoint, e.g. s3.amazonaws.com
	Bucket             string           // Dedicated bucket name.
	IncludeCollections []string         // Collections to be included in the backup.
	ExcludeCollections []string         // Collections to be excluded from the backup.
	MaxCPUPercentage   int              // Maximum %CPU utilization.
	ChunkSize          int              // Target chunk size.
	CompressionLevel   CompressionLevel // Hint for selecting the optimal compression algorithm.
}

func (c *Client) Create(ctx context.Context, cfg *Create) (*Info, error) {
	cfg = internal.Optional(cfg)

	req := &api.CreateBackupRequest{
		Backend:            cfg.Backend,
		ID:                 cfg.ID,
		Bucket:             cfg.Bucket,
		BackupPath:         cfg.Path,
		Endpoint:           cfg.Endpoint,
		IncludeCollections: cfg.IncludeCollections,
		ExcludeCollections: cfg.ExcludeCollections,
		MaxCPUPercentage:   cfg.MaxCPUPercentage,
		CompressionLevel:   api.BackupCompressionLevel(cfg.CompressionLevel),
	}

	var resp api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("create backup: %w", err)
	}
	return newInfo(&resp, c, api.CreateBackup), nil
}

type Restore struct {
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

func (c *Client) Restore(ctx context.Context, cfg *Restore) (*Info, error) {
	cfg = internal.Optional(cfg)

	req := &api.RestoreBackupRequest{
		Backend:            cfg.Backend,
		ID:                 cfg.ID,
		Bucket:             cfg.Bucket,
		BackupPath:         cfg.Path,
		Endpoint:           cfg.Endpoint,
		IncludeCollections: cfg.IncludeCollections,
		ExcludeCollections: cfg.ExcludeCollections,
		MaxCPUPercentage:   cfg.MaxCPUPercentage,
		OverwriteAlias:     cfg.OverwriteAlias,
		RestoreUsers:       api.RBACRestoreOption(cfg.RestoreUsers),
		RestoreRoles:       api.RBACRestoreOption(cfg.RestoreRoles),
		NodeMapping:        cfg.NodeMapping,
	}

	var resp api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("restore backup: %w", err)
	}
	return newInfo(&resp, c, api.RestoreBackup), nil
}

type GetStatus struct {
	Backend string // Required: Backend storage.
	ID      string // Required: Backup ID.
}

func (c *Client) GetCreateStatus(ctx context.Context, cfg GetStatus) (*Info, error) {
	return c.getStatus(ctx, cfg, api.CreateBackup)
}

func (c *Client) GetRestoreStatus(ctx context.Context, cfg GetStatus) (*Info, error) {
	return c.getStatus(ctx, cfg, api.RestoreBackup)
}

func (c *Client) getStatus(ctx context.Context, cfg GetStatus, operation api.BackupOperation) (*Info, error) {
	req := &api.BackupStatusRequest{
		Backend:   cfg.Backend,
		ID:        cfg.ID,
		Operation: operation,
	}

	var resp api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("get create backup status: %w", err)
	}
	return newInfo(&resp, c, operation), nil
}

type List struct {
	Backend         string // Required: Backend storage.
	StartingTimeAsc bool   // Set to true to order backups by their StartedAt time in ascending order.
}

func (c *Client) List(ctx context.Context, cfg List) ([]*Info, error) {
	req := &api.ListBackupsRequest{
		Backend:         cfg.Backend,
		StartingTimeAsc: cfg.StartingTimeAsc,
	}

	var resp []*api.BackupInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	infos := make([]*Info, len(resp))
	for _, bak := range resp {
		infos = append(infos, newInfo(bak, c, completed))
	}
	return infos, nil
}

type Cancel struct {
	Backend string // Required: Backend storage.
	ID      string // Required: Backup ID.
}

// Cancel an in-progress backup.
func (c *Client) Cancel(ctx context.Context, cfg Cancel) error {
	req := api.CancelBackupRequest{
		Backend: cfg.Backend,
		ID:      cfg.ID,
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
const completed api.BackupOperation = api.BackupOperation(api.CreateBackup - 1)

func newInfo(bak *api.BackupInfo, c *Client, op api.BackupOperation) *Info {
	return &Info{
		ID:                  bak.ID,
		Path:                bak.Path,
		Backend:             bak.Backend,
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
