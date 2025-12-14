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

type createRequest struct {
	Path               *string           // Path to backup in the backend storage.
	Bucket             *string           // Dedicated bucket name.
	IncludeCollections []string          // Collections to be included in the backup.
	ExcludeCollections []string          // Collections to be excluded from the backup.
	MaxCPUPercentage   *int              // Maximum %CPU utilization.
	CompressionLevel   *CompressionLevel // Hint for selecting the optimal compression algorithm.
}

type (
	CreateOption     interface{ create(*createRequest) }
	CreateOptionFunc func(*createRequest)
)

func (f CreateOptionFunc) create(r *createRequest) { f(r) }

// Specify desired compression level.
// The server will select an appropriate compression algorithm based on this setting.
func WithCompressionLevel(cl CompressionLevel) CreateOption {
	return CreateOptionFunc(func(r *createRequest) {
		r.CompressionLevel = &cl
	})
}

func (c *Client) Create(ctx context.Context, id string, backend string, options ...CreateOption) (*Info, error) {
	var cfg createRequest
	for _, opt := range options {
		opt.create(&cfg)
	}

	req := &api.CreateBackupRequest{
		ID:                 id,
		Backend:            backend,
		IncludeCollections: cfg.IncludeCollections,
		ExcludeCollections: cfg.ExcludeCollections,
		Config: &api.CreateBackupConfig{
			Path:             cfg.Path,
			Bucket:           cfg.Bucket,
			MaxCPUPercentage: cfg.MaxCPUPercentage,
			CompressionLevel: (*string)(cfg.CompressionLevel),
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
	}, nil
}

type restoreRequest struct {
	Path               *string                // Path to backup in the backend storage.
	Bucket             *string                // Dedicated bucket name.
	IncludeCollections []string               // Collections to be included in the backup.
	ExcludeCollections []string               // Collections to be excluded from the backup.
	MaxCPUPercentage   *int                   // Maximum %CPU utilization.
	OverwriteAlias     bool                   // Allow overwriting aliases.
	RestoreUsers       *api.RBACRestoreOption // Select strategy for restoring RBAC users.
	RestoreRoles       *api.RBACRestoreOption // Select strategy for restoring RBAC roles.
}

type (
	RestoreOption     interface{ restore(*restoreRequest) }
	RestoreOptionFunc func(*restoreRequest)
)

func (f RestoreOptionFunc) restore(r *restoreRequest) { f(r) }

func WithOverwriteAlias(overwrite bool) RestoreOption {
	return RestoreOptionFunc(func(r *restoreRequest) {
		r.OverwriteAlias = overwrite
	})
}

func WithRestoreUsers(opt RBACRestore) RestoreOption {
	return RestoreOptionFunc(func(r *restoreRequest) {
		r.RestoreUsers = (*api.RBACRestoreOption)(&opt)
	})
}

func WithRestoreRoles(opt RBACRestore) RestoreOption {
	return RestoreOptionFunc(func(r *restoreRequest) {
		r.RestoreRoles = (*api.RBACRestoreOption)(&opt)
	})
}

func (c *Client) Restore(ctx context.Context, id string, backend string, options ...RestoreOption) (*Info, error) {
	var cfg restoreRequest
	for _, opt := range options {
		opt.restore(&cfg)
	}

	req := &api.RestoreBackupRequest{
		ID:                 id,
		Backend:            backend,
		IncludeCollections: cfg.IncludeCollections,
		ExcludeCollections: cfg.ExcludeCollections,
		OverwriteAlias:     cfg.OverwriteAlias,
		Config: &api.RestoreBackupConfig{
			Bucket:           cfg.Bucket,
			Path:             cfg.Path,
			MaxCPUPercentage: cfg.MaxCPUPercentage,
			RestoreUsers:     cfg.RestoreUsers,
			RestoreRoles:     cfg.RestoreRoles,
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
	}, nil
}

func (c *Client) GetCreateStatus(ctx context.Context, id string, backend string) (*Info, error) {
	return c.getStatus(ctx, id, backend, api.CreateBackup)
}

func (c *Client) GetRestoreStatus(ctx context.Context, id string, backend string) (*Info, error) {
	return c.getStatus(ctx, id, backend, api.RestoreBackup)
}

func (c *Client) getStatus(ctx context.Context, id string, backend string, operation api.BackupOperation) (*Info, error) {
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
	}, nil
}

type listRequest struct {
	StartingTimeAsc bool
}

type ListOption func(*listRequest)

// Order backups by their starting time in ascending order.
func WithStartingTimeAsc(asc bool) ListOption {
	return func(r *listRequest) {
		r.StartingTimeAsc = asc
	}
}

func (c *Client) List(ctx context.Context, id string, backend string, options ...ListOption) ([]Info, error) {
	var cfg listRequest
	for _, opt := range options {
		opt(&cfg)
	}

	req := &api.ListBackupsRequest{
		ID:              id,
		Backend:         backend,
		StartingTimeAsc: cfg.StartingTimeAsc,
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

func (c *Client) Cancel(ctx context.Context, id string, backend string) error {
	req := api.CancelBackupRequest{ID: id, Backend: backend}
	err := c.transport.Do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("cancel backup: %w", err)
	}
	return nil
}

// Customize backup's path within the bucket.
type WithPath string

var (
	_ CreateOption  = (*WithPath)(nil)
	_ RestoreOption = (*WithPath)(nil)
)

func (opt WithPath) create(r *createRequest)   { r.Path = (*string)(&opt) }
func (opt WithPath) restore(r *restoreRequest) { r.Path = (*string)(&opt) }

// Set dedicated bucket for this backup.
type WithBucket string

var (
	_ CreateOption  = (*WithBucket)(nil)
	_ RestoreOption = (*WithBucket)(nil)
)

func (opt WithBucket) create(r *createRequest)   { r.Bucket = (*string)(&opt) }
func (opt WithBucket) restore(r *restoreRequest) { r.Bucket = (*string)(&opt) }

// Include collections in the backup.
type WithIncludeCollections []string

var (
	_ CreateOption  = (*WithIncludeCollections)(nil)
	_ RestoreOption = (*WithIncludeCollections)(nil)
)

func (opt WithIncludeCollections) create(r *createRequest) {
	r.IncludeCollections = append(r.IncludeCollections, opt...)
}

func (opt WithIncludeCollections) restore(r *restoreRequest) {
	r.IncludeCollections = append(r.IncludeCollections, opt...)
}

// Exclude collections from the backup.
type WithExcludeCollections []string

var (
	_ CreateOption  = (*WithExcludeCollections)(nil)
	_ RestoreOption = (*WithExcludeCollections)(nil)
)

func (opt WithExcludeCollections) create(r *createRequest) {
	r.ExcludeCollections = append(r.ExcludeCollections, opt...)
}

func (opt WithExcludeCollections) restore(r *restoreRequest) {
	r.ExcludeCollections = append(r.ExcludeCollections, opt...)
}

// Limit CPU resources that will be allocated to the backup process.
type WithMaxCPUPercentage int

var (
	_ CreateOption  = (*WithMaxCPUPercentage)(nil)
	_ RestoreOption = (*WithMaxCPUPercentage)(nil)
)

func (opt WithMaxCPUPercentage) create(r *createRequest)   { r.MaxCPUPercentage = (*int)(&opt) }
func (opt WithMaxCPUPercentage) restore(r *restoreRequest) { r.MaxCPUPercentage = (*int)(&opt) }
