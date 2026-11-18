package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type ExportInfo struct {
	Backend             string
	ID                  string
	Path                string
	Status              ExportStatus
	PerShardStatus      map[string]map[string]ShardExportStatus
	IncludesCollections []string
	Error               string

	StartedAt   *time.Time
	CompletedAt *time.Time
	Took        time.Duration
}

var _ json.Unmarshaler = (*ExportInfo)(nil)

type ShardExportStatus struct {
	Status          ExportStatus
	ObjectsExported int64
	SkipReason      string
	Error           string
}

var _ json.Unmarshaler = (*ShardExportStatus)(nil)

type ExportStatus string

const (
	ExportStatusStarted      = ExportStatus(rest.ExportStatusResponseStatusSTARTED)
	ExportStatusTransferring = ExportStatus(rest.ExportStatusResponseStatusTRANSFERRING)
	ExportStatusSuccess      = ExportStatus(rest.ExportStatusResponseStatusSUCCESS)
	ExportStatusFailed       = ExportStatus(rest.ExportStatusResponseStatusFAILED)
	ExportStatusCanceled     = ExportStatus(rest.ExportStatusResponseStatusCANCELED)
)

type ExportFileFormat string

const (
	ExportFileFormatParquet = ExportFileFormat(rest.Parquet)
)

type CreateExportRequest struct {
	transports.BaseEndpoint

	Backend            string
	ID                 string
	FileFormat         ExportFileFormat
	IncludeCollections []string
	ExcludeCollections []string
}

var _ transports.Endpoint = (*CreateExportRequest)(nil)

func (*CreateExportRequest) Method() string { return http.MethodPost }
func (r *CreateExportRequest) Path() string { return "/export/" + r.Backend }
func (r *CreateExportRequest) Body() any {
	return rest.ExportCreateRequest{
		Id:         r.ID,
		FileFormat: rest.ExportCreateRequestFileFormat(r.FileFormat),
		Include:    r.IncludeCollections,
		Exclude:    r.ExcludeCollections,
	}
}

type GetExportStatusRequest struct {
	transports.BaseEndpoint
	RequestDefaults

	Backend string
	ID      string
}

var _ transports.Endpoint = (*GetExportStatusRequest)(nil)

func (*GetExportStatusRequest) Method() string { return http.MethodGet }
func (r *GetExportStatusRequest) Path() string { return "/export/" + r.Backend + "/" + r.ID }

type CancelExportRequest struct {
	transports.BaseEndpoint

	Backend string
	ID      string
}

func (*CancelExportRequest) Method() string { return http.MethodDelete }
func (r *CancelExportRequest) Path() string { return "/export/" + r.Backend + "/" + r.ID }

func (s *ShardExportStatus) UnmarshalJSON(data []byte) error {
	var v rest.ShardProgress
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*s = ShardExportStatus{
		Status:          ExportStatus(v.Status),
		ObjectsExported: v.ObjectsExported,
		SkipReason:      v.SkipReason,
		Error:           v.Error,
	}
	return nil
}

func (ei *ExportInfo) UnmarshalJSON(data []byte) error {
	var resp struct {
		Backend             string                                  `json:"backend,omitempty"`
		ID                  string                                  `json:"id,omitempty"`
		Path                string                                  `json:"path,omitempty"`
		Status              string                                  `json:"status,omitempty"`
		IncludesCollections []string                                `json:"classes,omitempty"`
		StartedAt           *time.Time                              `json:"startedAt,omitempty"`
		CompletedAt         *time.Time                              `json:"completedAt,omitempty"`
		Error               string                                  `json:"error,omitempty"`
		PerShardStatus      map[string]map[string]ShardExportStatus `json:"shardStatus,omitempty"`
		TookMillis          int64                                   `json:"tookInMs,omitempty"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	*ei = ExportInfo{
		Backend:             resp.Backend,
		ID:                  resp.ID,
		Path:                resp.Path,
		Status:              ExportStatus(resp.Status),
		PerShardStatus:      resp.PerShardStatus,
		IncludesCollections: resp.IncludesCollections,
		StartedAt:           resp.StartedAt,
		CompletedAt:         resp.CompletedAt,
		Error:               resp.Error,
		Took:                time.Duration(resp.TookMillis) * time.Millisecond,
	}
	return nil
}
