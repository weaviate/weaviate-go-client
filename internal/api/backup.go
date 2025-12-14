package api

import (
	"net/http"
	"net/url"
	"time"
)

type (
	CreateBackupConfig struct {
		Bucket           *string `json:"Bucket,omitempty"`
		Path             *string `json:"Path,omitempty"`
		MaxCPUPercentage *int    `json:"CPUPercentage,omitempty"`
		CompressionLevel *string `json:"CompressionLevel,omitempty"`
	}
	RestoreBackupConfig struct {
		Bucket           *string            `json:"Bucket,omitempty"`
		Path             *string            `json:"Path,omitempty"`
		MaxCPUPercentage *int               `json:"CPUPercentage,omitempty"`
		RestoreUsers     *RBACRestoreOption `json:"restoreUsers,omitempty"`
		RestoreRoles     *RBACRestoreOption `json:"restoreRoles,omitempty"`
	}
	Backup struct {
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
	endpoint
	Backend string `json:"-"`

	ID                 string              `json:"id"`
	IncludeCollections []string            `json:"include,omitempty"`
	ExcludeCollections []string            `json:"exclude,omitempty"`
	Config             *CreateBackupConfig `json:"config,omitempty"`
}

var _ Endpoint = (*CreateBackupRequest)(nil)

func (*CreateBackupRequest) Method() string { return http.MethodPost }
func (r *CreateBackupRequest) Path() string { return "/backups/" + r.Backend }
func (r *CreateBackupRequest) Body() any    { return r }

type RestoreBackupRequest struct {
	endpoint
	Backend string `json:"-"`
	ID      string `json:"-"`

	IncludeCollections []string             `json:"include,omitempty"`
	ExcludeCollections []string             `json:"exclude,omitempty"`
	OverwriteAlias     bool                 `json:"overwriteAlias,omitempty"`
	Config             *RestoreBackupConfig `json:"config,omitempty"`
}

var _ Endpoint = (*RestoreBackupRequest)(nil)

func (*RestoreBackupRequest) Method() string { return http.MethodPost }
func (r *RestoreBackupRequest) Path() string { return "/backups/" + r.Backend }
func (r *RestoreBackupRequest) Body() any    { return r }

type BackupStatusRequest struct {
	endpoint
	Backend   string
	ID        string
	Operation BackupOperation
}

var _ Endpoint = (*BackupStatusRequest)(nil)

func (r *BackupStatusRequest) Method() string { return http.MethodGet }
func (r *BackupStatusRequest) Path() string {
	path := "/backups/" + r.Backend + "/" + r.ID
	if r.Operation == RestoreBackup {
		path += "/restore"
	}
	return path
}

type ListBackupsRequest struct {
	endpoint
	Backend         string
	ID              string
	StartingTimeAsc bool
}

var _ Endpoint = (*ListBackupsRequest)(nil)

func (r *ListBackupsRequest) Method() string { return http.MethodGet }
func (r *ListBackupsRequest) Path() string   { return "/backups/" + r.Backend }
func (r *ListBackupsRequest) Query() url.Values {
	if !r.StartingTimeAsc {
		return nil
	}
	return map[string][]string{"order": {"asc"}}
}

type CancelBackupRequest struct {
	endpoint
	Backend string `json:"-"`
	ID      string `json:"-"`
}

var _ Endpoint = (*CancelBackupRequest)(nil)

func (r *CancelBackupRequest) Method() string { return http.MethodDelete }
func (r *CancelBackupRequest) Path() string   { return "/backups/" + r.Backend + "/" + r.ID }
