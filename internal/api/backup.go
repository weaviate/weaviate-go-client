package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

type BackupInfo struct {
	Backend     string
	ID          string
	Path        string
	Error       string
	Status      BackupStatus
	StartedAt   time.Time
	CompletedAt time.Time

	IncludesCollections []string
	SizeGiB             *float32
}

var _ json.Unmarshaler = (*BackupInfo)(nil)

type BackupCompressionLevel string

const (
	BackupCompressionLevelDefault             BackupCompressionLevel = BackupCompressionLevel(rest.DefaultCompression)
	BackupCompressionLevelBestSpeed           BackupCompressionLevel = BackupCompressionLevel(rest.BestSpeed)
	BackupCompressionLevelBestCompression     BackupCompressionLevel = BackupCompressionLevel(rest.BestCompression)
	BackupCompressionLevelZstdDefault         BackupCompressionLevel = "ZstdDefaultCompression"
	BackupCompressionLevelZstdBestSpeed       BackupCompressionLevel = "ZstdBestSpeed"
	BackupCompressionLevelZstdBestCompression BackupCompressionLevel = "ZstdBestCompression"
	BackupCompressionLevelNone                BackupCompressionLevel = "NoCompression"
)

type BackupStatus string

const (
	BackupStatusStarted      BackupStatus = BackupStatus(rest.BackupListResponseStatusSTARTED)
	BackupStatusTransferring BackupStatus = BackupStatus(rest.BackupListResponseStatusTRANSFERRING)
	BackupStatusTransferred  BackupStatus = BackupStatus(rest.BackupListResponseStatusTRANSFERRED)
	BackupStatusSuccess      BackupStatus = BackupStatus(rest.BackupListResponseStatusSUCCESS)
	BackupStatusFailed       BackupStatus = BackupStatus(rest.BackupListResponseStatusFAILED)
	BackupStatusCanceled     BackupStatus = BackupStatus(rest.BackupListResponseStatusCANCELED)
)

type RBACRestoreOption string

const (
	RBACRestoreAll  RBACRestoreOption = RBACRestoreOption(rest.All)
	RBACRestoreNone RBACRestoreOption = RBACRestoreOption(rest.NoRestore)
)

type BackupOperation int

const (
	CreateBackup BackupOperation = iota
	RestoreBackup
)

type CreateBackupRequest struct {
	transport.BaseEndpoint

	Backend            string // Required: backend storage.
	ID                 string // Required: backup ID.
	BackupPath         string
	Endpoint           string
	Bucket             string
	IncludeCollections []string
	ExcludeCollections []string
	MaxCPUPercentage   int
	ChunkSizeMiB       int
	CompressionLevel   BackupCompressionLevel
}

// Compile-time assertion that CreateBackupRequest implements [transport.Endpoint].
var (
	_ transport.Endpoint = (*CreateBackupRequest)(nil)
	_ json.Marshaler     = (*CreateBackupRequest)(nil)
)

func (*CreateBackupRequest) Method() string { return http.MethodPost }
func (r *CreateBackupRequest) Path() string { return "/backups/" + r.Backend }
func (r *CreateBackupRequest) Body() any    { return r }

type RestoreBackupRequest struct {
	transport.BaseEndpoint

	Backend            string
	ID                 string
	BackupPath         string
	Endpoint           string
	Bucket             string
	IncludeCollections []string
	ExcludeCollections []string
	MaxCPUPercentage   int
	OverwriteAlias     bool
	RestoreUsers       RBACRestoreOption
	RestoreRoles       RBACRestoreOption
	NodeMapping        map[string]string
}

// Compile-time assertion that RestoreBackupRequest implements [transport.Endpoint].
var (
	_ transport.Endpoint = (*RestoreBackupRequest)(nil)
	_ json.Marshaler     = (*RestoreBackupRequest)(nil)
)

func (*RestoreBackupRequest) Method() string { return http.MethodPost }
func (r *RestoreBackupRequest) Path() string {
	return "/backups/" + r.Backend + "/" + r.ID + "/restore"
}
func (r *RestoreBackupRequest) Body() any { return r }

type BackupStatusRequest struct {
	transport.BaseEndpoint

	Backend   string
	ID        string
	Operation BackupOperation
}

// Compile-time assertion that BackupStatusRequest implements [transport.Endpoint].
var _ transport.Endpoint = (*BackupStatusRequest)(nil)

func (*BackupStatusRequest) Method() string { return http.MethodGet }
func (r *BackupStatusRequest) Path() string {
	path := "/backups/" + r.Backend + "/" + r.ID
	if r.Operation == RestoreBackup {
		path += "/restore"
	}
	return path
}

// ListBackupsRequest fetches all requests in a backend storage.
type ListBackupsRequest struct {
	transport.BaseEndpoint

	Backend         string
	StartingTimeAsc bool
}

// Compile-time assertion that ListBackupsRequest implements [transport.Endpoint].
var _ transport.Endpoint = (*ListBackupsRequest)(nil)

func (*ListBackupsRequest) Method() string { return http.MethodGet }
func (r *ListBackupsRequest) Path() string { return "/backups/" + r.Backend }
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

// Compile-time assertion that CancelBackupRequest implements [transport.Endpoint].
var _ transport.Endpoint = (*CancelBackupRequest)(nil)

func (*CancelBackupRequest) Method() string { return http.MethodDelete }
func (r *CancelBackupRequest) Path() string { return "/backups/" + r.Backend + "/" + r.ID }

// MarshalJSON implements json.Marshaler via rest.BackupCreateRequest.
func (r *CreateBackupRequest) MarshalJSON() ([]byte, error) {
	req := &rest.BackupCreateRequest{
		Id:      r.ID,
		Include: r.IncludeCollections,
		Exclude: r.ExcludeCollections,
		Config: rest.BackupConfig{
			Path:             r.BackupPath,
			Bucket:           r.Bucket,
			Endpoint:         r.Endpoint,
			CPUPercentage:    r.MaxCPUPercentage,
			ChunkSize:        r.ChunkSizeMiB,
			CompressionLevel: rest.BackupConfigCompressionLevel(r.CompressionLevel),
		},
	}
	return json.Marshal(req)
}

// MarshalJSON implements json.Marshaler via rest.BackupRestoreRequest.
func (r *RestoreBackupRequest) MarshalJSON() ([]byte, error) {
	req := &rest.BackupRestoreRequest{
		Include:        r.IncludeCollections,
		Exclude:        r.ExcludeCollections,
		OverwriteAlias: r.OverwriteAlias,
		NodeMapping:    r.NodeMapping,
		Config: rest.RestoreConfig{
			Bucket:        r.Bucket,
			Path:          r.BackupPath,
			Endpoint:      r.Endpoint,
			CPUPercentage: r.MaxCPUPercentage,
			RolesOptions:  rest.RestoreConfigRolesOptions(r.RestoreRoles),
			UsersOptions:  rest.RestoreConfigUsersOptions(r.RestoreUsers),
		},
	}
	return json.Marshal(req)
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BackupInfo) UnmarshalJSON(data []byte) error {
	var bak struct {
		ID          string       `json:"id,omitempty"`
		Path        string       `json:"path,omitempty"`
		Backend     string       `json:"backend,omitempty"`
		Error       string       `json:"error,omitempty"`
		Status      BackupStatus `json:"status,omitempty"`
		StartedAt   time.Time    `json:"startedAt"`
		CompletedAt time.Time    `json:"completedAt"`

		IncludesCollections []string `json:"classes,omitempty"`
		SizeGiB             *float32 `json:"size,omitempty"`
	}

	if err := json.Unmarshal(data, &bak); err != nil {
		return err
	}

	*b = BackupInfo{
		Backend:             bak.Backend,
		ID:                  bak.ID,
		Path:                bak.Path,
		Error:               bak.Error,
		Status:              bak.Status,
		StartedAt:           bak.StartedAt,
		CompletedAt:         bak.CompletedAt,
		IncludesCollections: bak.IncludesCollections,
		SizeGiB:             bak.SizeGiB,
	}
	return nil
}
