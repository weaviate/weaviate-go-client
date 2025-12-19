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
	ID                  string    // Backup ID
	Path                string    // Path to backup in the backend storage
	Backend             string    // Backup storage backend
	IncludesCollections []string  // Collections included in the backup
	Status              Status    // Backup creation / restoration status.
	Error               string    // Backup creation / restoration error.
	StartedAt           time.Time // Time at which the backup creation started.
	CompletedAt         time.Time // Time at which the backup was completed, successfully or otherwise.
	SizeGiB             float32   // Backup size in GiB.

	c         *Client
	operation api.BackupOperation
}

// IsCompleted returns true if the backup operation has completed, successfully or otherwise.
func (i *Info) IsCompleted() bool {
	return i.Status == StatusSuccess || i.Status == StatusFailed || i.Status == StatusCanceled
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
	Path               string           // Path to backup in the backend storage.
	Bucket             string           // Dedicated bucket name.
	IncludeCollections []string         // Collections to be included in the backup.
	ExcludeCollections []string         // Collections to be excluded from the backup.
	MaxCPUPercentage   int              // Maximum %CPU utilization.
	CompressionLevel   CompressionLevel // Hint for selecting the optimal compression algorithm.
}

func (c *Client) Create(ctx context.Context, id, backend string, option ...Create) (*Info, error) {
	opt, _ := internal.Last(option...)

	req := &api.CreateBackupRequest{
		ID:                 id,
		Backend:            backend,
		IncludeCollections: opt.IncludeCollections,
		ExcludeCollections: opt.ExcludeCollections,
		Config: &api.CreateBackupConfig{
			Path:             opt.Path,
			Bucket:           opt.Bucket,
			MaxCPUPercentage: opt.MaxCPUPercentage,
			CompressionLevel: (*string)(opt.CompressionLevel),
		},
	}

	var resp api.Backup
	err := c.transport.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("create backup: %w", err)
	}

	return &Info{
		ID:                  resp.ID,
		Path:                resp.Path,
		Backend:             resp.Backend,
		IncludesCollections: resp.IncludesCollections,
		Status:              Status(resp.Status),
		Error:               resp.Error,
		StartedAt:           resp.StartedAt,
		CompletedAt:         resp.CompletedAt,
		SizeGiB:             resp.SizeGiB,

		operation: api.CreateBackup,
		c:         c,
	}, nil
}

type Restore struct {
	Path               string                // Path to backup in the backend storage.
	Bucket             string                // Dedicated bucket name.
	IncludeCollections []string              // Collections to be included in the backup.
	ExcludeCollections []string              // Collections to be excluded from the backup.
	MaxCPUPercentage   int                   // Maximum %CPU utilization.
	OverwriteAlias     bool                  // Allow overwriting aliases.
	RestoreUsers       api.RBACRestoreOption // Select strategy for restoring RBAC users.
	RestoreRoles       api.RBACRestoreOption // Select strategy for restoring RBAC roles.
}

func (c *Client) Restore(ctx context.Context, id, backend string, option ...Restore) (*Info, error) {
	opt, _ := internal.Last(option...)

	req := &api.RestoreBackupRequest{
		ID:                 id,
		Backend:            backend,
		IncludeCollections: opt.IncludeCollections,
		ExcludeCollections: opt.ExcludeCollections,
		OverwriteAlias:     opt.OverwriteAlias,
		Config: &api.RestoreBackupConfig{
			Bucket:           opt.Bucket,
			Path:             opt.Path,
			MaxCPUPercentage: opt.MaxCPUPercentage,
			RestoreUsers:     opt.RestoreUsers,
			RestoreRoles:     opt.RestoreRoles,
		},
	}

	var resp api.Backup
	err := c.transport.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("restore backup: %w", err)
	}

	return &Info{
		ID:                  resp.ID,
		Path:                resp.Path,
		Backend:             resp.Backend,
		IncludesCollections: resp.IncludesCollections,
		Status:              Status(resp.Status),
		Error:               resp.Error,
		StartedAt:           resp.StartedAt,
		CompletedAt:         resp.CompletedAt,
		SizeGiB:             resp.SizeGiB,

		operation: api.RestoreBackup,
		c:         c,
	}, nil
}

func (c *Client) GetCreateStatus(ctx context.Context, id, backend string) (*Info, error) {
	return c.getStatus(ctx, id, backend, api.CreateBackup)
}

func (c *Client) GetRestoreStatus(ctx context.Context, id, backend string) (*Info, error) {
	return c.getStatus(ctx, id, backend, api.RestoreBackup)
}

func (c *Client) getStatus(ctx context.Context, id, backend string, operation api.BackupOperation) (*Info, error) {
	req := &api.BackupStatusRequest{
		ID:        id,
		Backend:   backend,
		Operation: operation,
	}

	var resp api.Backup
	err := c.transport.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("get create backup status: %w", err)
	}
	return &Info{
		ID:                  resp.ID,
		Path:                resp.Path,
		Backend:             resp.Backend,
		IncludesCollections: resp.IncludesCollections,
		Status:              Status(resp.Status),
		Error:               resp.Error,
		StartedAt:           resp.StartedAt,
		CompletedAt:         resp.CompletedAt,
		SizeGiB:             resp.SizeGiB,

		operation: operation,
		c:         c,
	}, nil
}

type List struct {
	StartingTimeAsc bool
}

func (c *Client) List(ctx context.Context, id, backend string, option ...List) ([]Info, error) {
	opt, _ := internal.Last(option...)

	req := &api.ListBackupsRequest{
		ID:              id,
		Backend:         backend,
		StartingTimeAsc: opt.StartingTimeAsc,
	}

	var resp []api.Backup
	err := c.transport.Do(ctx, req, &resp)
	if err != nil {
		return nil, fmt.Errorf("list backups: %w", err)
	}

	backups := make([]Info, len(resp))
	for _, b := range resp {
		backups = append(backups, Info{
			ID:                  b.ID,
			Path:                b.Path,
			Backend:             b.Backend,
			IncludesCollections: b.IncludesCollections,
			Status:              Status(b.Status),
			Error:               b.Error,
			StartedAt:           b.StartedAt,
			CompletedAt:         b.CompletedAt,
			SizeGiB:             b.SizeGiB,
		})
	}
	return backups, nil
}

// Cancel an in-progress backup.
func (c *Client) Cancel(ctx context.Context, id, backend string) error {
	req := api.CancelBackupRequest{ID: id, Backend: backend}
	err := c.transport.Do(ctx, &req, nil)
	if err != nil {
		return fmt.Errorf("cancel backup: %w", err)
	}
	return nil
}
