package api

import (
	"net/http"
	"net/url"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

type (
	RestoreBackupConfig struct{}
	Backup              struct {
		ID                  string
		Path                string
		Backend             string
		IncludesCollections []string
		Status              string
		Error               string
		StartedAt           time.Time
		CompletedAt         time.Time
		SizeGiB             float32
	}
)

type BackupCompressionLevel string

const (
	BackupCompressionLevelDefault             BackupCompressionLevel = "DefaultCompression"
	BackupCompressionLevelBestSpeed           BackupCompressionLevel = "BestSpeed"
	BackupCompressionLevelBestCompression     BackupCompressionLevel = "BestCompression"
	BackupCompressionLevelZstdDefault         BackupCompressionLevel = "ZstdDefaultCompression"
	BackupCompressionLevelZstdBestSpeed       BackupCompressionLevel = "ZstdBestSpeed"
	BackupCompressionLevelZstdBestCompression BackupCompressionLevel = "ZstdBestCompression"
	BackupCompressionLevelNone                BackupCompressionLevel = "NoCompression"
)

type BackupStatus string

const (
	BackupStatusStarted      BackupStatus = "STARTED"
	BackupStatusTransferring BackupStatus = "TRANSFERRING"
	BackupStatusSuccess      BackupStatus = "SUCCESS"
	BackupStatusFailed       BackupStatus = "FAILED"
	BackupStatusCanceled     BackupStatus = "CANCELED"
)

type RBACRestoreOption string

const (
	RBACRestoreAll  RBACRestoreOption = "all"
	RBACRestoreNone RBACRestoreOption = "noRestore"
)

type BackupOperation int

const (
	CreateBackup BackupOperation = iota
	RestoreBackup
)

type CreateBackupRequest struct {
	transport.BaseEndpoint

	Backend            string
	ID                 string
	Bucket             string
	BackupPath         string
	Endpoint           string
	IncludeCollections []string
	ExcludeCollections []string
	MaxCPUPercentage   int
	ChunkSize          int
	CompressionLevel   string
}

// Compile-time assertion that CreateBackupRequest implements [tranport.Endpoint].
var _ transport.Endpoint = (*CreateBackupRequest)(nil)

func (*CreateBackupRequest) Method() string { return http.MethodPost }
func (r *CreateBackupRequest) Path() string { return "/backups/" + r.Backend }
func (r *CreateBackupRequest) Body() any {
	return &rest.BackupCreateRequest{
		Id:      r.ID,
		Include: r.IncludeCollections,
		Exclude: r.ExcludeCollections,
		Config: rest.BackupConfig{
			Bucket:           r.Bucket,
			Endpoint:         r.Endpoint,
			CPUPercentage:    r.MaxCPUPercentage,
			ChunkSize:        r.ChunkSize,
			CompressionLevel: rest.BackupConfigCompressionLevel(r.CompressionLevel),
		},
	}
}

type RestoreBackupRequest struct {
	transport.BaseEndpoint

	Backend            string
	ID                 string
	Bucket             string
	BackupPath         string
	Endpoint           string
	IncludeCollections []string
	ExcludeCollections []string
	OverwriteAlias     bool
	MaxCPUPercentage   int
	RestoreUsers       RBACRestoreOption
	RestoreRoles       RBACRestoreOption
}

// Compile-time assertion that RestoreBackupRequest implements [transport.Endpoint].
var _ transport.Endpoint = (*RestoreBackupRequest)(nil)

func (*RestoreBackupRequest) Method() string { return http.MethodPost }
func (r *RestoreBackupRequest) Path() string { return "/backups/" + r.Backend }
func (r *RestoreBackupRequest) Body() any {
	return &rest.BackupRestoreRequest{
		Include:        r.IncludeCollections,
		Exclude:        r.ExcludeCollections,
		OverwriteAlias: r.OverwriteAlias,
		// NodeMapping:    make(map[string]string), // TODO(dyma): what should this field be?
		Config: rest.RestoreConfig{
			Bucket:        r.Bucket,
			Path:          r.BackupPath,
			Endpoint:      r.Endpoint,
			CPUPercentage: r.MaxCPUPercentage,
			RolesOptions:  rest.RestoreConfigRolesOptions(r.RestoreRoles),
			UsersOptions:  rest.RestoreConfigUsersOptions(r.RestoreUsers),
		},
	}
}

type BackupStatusRequest struct {
	transport.BaseEndpoint

	Backend   string
	ID        string
	Operation BackupOperation
}

// Compile-time assertion that BackupStatusRequest implements [tranport.Endpoint].
var (
	_ transport.Endpoint = (*BackupStatusRequest)(nil)
)

func (r *BackupStatusRequest) Method() string { return http.MethodGet }
func (r *BackupStatusRequest) Path() string {
	path := "/backups/" + r.Backend + "/" + r.ID
	if r.Operation == RestoreBackup {
		path += "/restore"
	}
	return path
}

type ListBackupsRequest struct {
	transport.BaseEndpoint

	Backend         string
	ID              string
	StartingTimeAsc bool
}

// Compile-time assertion that ListBackupsRequest implements [tranport.Endpoint].
var _ transport.Endpoint = (*ListBackupsRequest)(nil)

func (r *ListBackupsRequest) Method() string { return http.MethodGet }
func (r *ListBackupsRequest) Path() string   { return "/backups/" + r.Backend }
func (r *ListBackupsRequest) Query() url.Values {
	if !r.StartingTimeAsc {
		return nil
	}
	return url.Values{"order": {"asc"}}
}

type CancelBackupRequest struct {
	transport.BaseEndpoint

	Backend string
	ID      string
}

// Compile-time assertion that CancelBackupRequest implements [tranport.Endpoint].
var _ transport.Endpoint = (*CancelBackupRequest)(nil)

func (r *CancelBackupRequest) Method() string { return http.MethodDelete }
func (r *CancelBackupRequest) Path() string   { return "/backups/" + r.Backend + "/" + r.ID }
