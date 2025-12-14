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
	Status              any       // Backup creation / restoration status.
	Error               string    // Backup creation / restoration error.
	StartedAt           time.Time // Time at which the backup creation started.
	CompletedAt         time.Time // Time at which the backup was completed, successfully or otherwise.
	SizeGiB             float32   // Backup size in GiB.
}

// Customize backup's path within the bucket.
func WithPath(path string) CreateOption {
	return func(r *createBackupRequest) {
		r.Path = &path
	}
}

// Select backend storage used for this backup.
func WithBackend(backend string) CreateOption {
	return func(r *createBackupRequest) {
		r.Backend = &backend
	}
}

// Set dedicated bucket for this backup.
func WithBucket(bucket string) CreateOption {
	return func(r *createBackupRequest) {
		r.Bucket = &bucket
	}
}

// Include collections in the backup.
func WithIncludeCollections(collections ...string) CreateOption {
	return func(r *createBackupRequest) {
		r.IncludeCollections = append(r.IncludeCollections, collections...)
	}
}

// Exclude collections from the backup.
func WithExcludeCollections(collections ...string) CreateOption {
	return func(r *createBackupRequest) {
		r.ExcludeCollections = append(r.ExcludeCollections, collections...)
	}
}

// Limit CPU resources that will be allocated to the backup process.
func WithMaxCPUPercentage(cpu int) CreateOption {
	return func(r *createBackupRequest) {
		r.MaxCPUPercentage = &cpu
	}
}

// Specify desired compression level.
// The server will select an appropriate compression algorithm based on this setting.
func WithCompressionLevel(cl CompressionLevel) CreateOption {
	return func(r *createBackupRequest) {
		r.CompressionLevel = &cl
	}
}

type CompressionLevel string

const (
	CompressionLevelDefault             CompressionLevel = "DefaultCompression"
	CompressionLevelBestSpeed           CompressionLevel = "BestSpeed"
	CompressionLevelBestCompression     CompressionLevel = "BestCompression"
	CompressionLevelZstdDefault         CompressionLevel = "ZstdDefaultCompression"
	CompressionLevelZstdBestSpeed       CompressionLevel = "ZstdBestSpeed"
	CompressionLevelZstdBestCompression CompressionLevel = "ZstdBestCompression"
	CompressionLevelNono                CompressionLevel = "NoCompression"
)

type createBackupRequest struct {
	ID                 string            // Backup ID.
	Path               *string           // Path to backup in the backend storage.
	Backend            *string           // Backup storage backend.
	Bucket             *string           // Dedicated bucket name.
	IncludeCollections []string          // Collections to be included in the backup.
	ExcludeCollections []string          // Collections to be excluded from the backup.
	MaxCPUPercentage   *int              // Maximum %CPU utilization.
	CompressionLevel   *CompressionLevel // Hint for selecting the optimal compression algorithm.
}

type CreateOption func(*createBackupRequest)

func (c *Client) Create(ctx context.Context, id string, backend string, options ...CreateOption) (*Info, error) {
	var cbr createBackupRequest
	for _, opt := range options {
		opt(&cbr)
	}

	req := &api.CreateBackupRequest{
		Backend:            backend,
		ID:                 cbr.ID,
		IncludeCollections: cbr.IncludeCollections,
		ExcludeCollections: cbr.ExcludeCollections,
		Config: &api.BackupConfig{
			Path:             cbr.Path,
			Bucket:           cbr.Bucket,
			MaxCPUPercentage: cbr.MaxCPUPercentage,
			CompressionLevel: (*string)(cbr.CompressionLevel),
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
		Status:              resp.Status,
		Error:               resp.Error,
		StartedAt:           resp.StartedAt,
		CompletedAt:         resp.CompletedAt,
		SizeGiB:             resp.SizeGiB,
	}, nil
}
