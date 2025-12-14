package api

import (
	"net/http"
	"time"
)

type (
	BackupConfig struct {
		Bucket           *string `json:"Bucket,omitempty"`
		Path             *string `json:"Path,omitempty"`
		MaxCPUPercentage *int    `json:"CPUPercentage,omitempty"`
		CompressionLevel *string `json:"CompressionLevel,omitempty"`
	}
	Backup struct {
		ID                  string
		Path                string
		Backend             string
		IncludesCollections []string
		Status              any
		Error               string
		StartedAt           time.Time
		CompletedAt         time.Time
		SizeGiB             float32
	}
)

type CreateBackupRequest struct {
	ID                 string        `json:"id"`
	Backend            string        `json:"-"`
	IncludeCollections []string      `json:"include,omitempty"`
	ExcludeCollections []string      `json:"exclude,omitempty"`
	Config             *BackupConfig `json:"config,omitempty"`
}

var _ Endpoint = (*CreateBackupRequest)(nil)

func (r *CreateBackupRequest) Method() string   { return http.MethodPost }
func (r *CreateBackupRequest) Endpoint() string { return "/backups/" + r.Backend }
