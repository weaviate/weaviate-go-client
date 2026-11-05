package export

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
	Backend             string
	ID                  string
	Path                string
	Status              Status
	PerShardStatus      map[string]map[string]ShardStatus
	IncludesCollections []string
	Error               string

	StartedAt   *time.Time
	CompletedAt *time.Time
	Took        time.Duration

	c *Client
}

type ShardStatus struct {
	Status          Status
	ObjectsExported int64
	SkipReason      string
	Error           string
}

type Status string

const (
	StatusStarted      = Status(api.ExportStatusStarted)
	StatusTransferring = Status(api.ExportStatusTransferring)
	StatusSuccess      = Status(api.ExportStatusSuccess)
	StatusFailed       = Status(api.ExportStatusFailed)
	StatusCanceled     = Status(api.ExportStatusCanceled)
)

type FileFormat string

const (
	FileFormatParquet = FileFormat(api.ExportFileFormatParquet)
)

type CreateOptions struct {
	Backend            string
	ID                 string
	FileFormat         FileFormat
	IncludeCollections []string
	ExcludeCollections []string
}

// Create export operation.
func (c *Client) Create(ctx context.Context, options CreateOptions) (*Info, error) {
	req := api.CreateExportRequest{
		Backend:            options.Backend,
		ID:                 options.ID,
		FileFormat:         api.ExportFileFormat(options.FileFormat),
		IncludeCollections: options.IncludeCollections,
		ExcludeCollections: options.ExcludeCollections,
	}

	var resp api.ExportInfo
	if err := c.transport.Do(ctx, &req, &resp); err != nil {
		return nil, fmt.Errorf("create export: %w", err)
	}
	return &Info{
		Backend:             resp.Backend,
		ID:                  resp.ID,
		Path:                resp.Path,
		IncludesCollections: resp.IncludesCollections,
		Status:              Status(resp.Status),
		StartedAt:           resp.StartedAt,
		c:                   c,
	}, nil
}

type GetStatusOptions struct {
	Backend string
	ID      string
}

// Get status of an export operation.
func (c *Client) GetCreateStatus(ctx context.Context, options GetStatusOptions) (*Info, error) {
	req := api.GetExportStatusRequest{
		Backend: options.Backend,
		ID:      options.ID,
	}

	var resp api.ExportInfo
	if err := c.transport.Do(ctx, &req, &resp); err != nil {
		return nil, fmt.Errorf("get export status: %w", err)
	}

	perShard := internal.MakeMap[string, map[string]ShardStatus](len(resp.PerShardStatus))
	for s, shard := range resp.PerShardStatus {
		nodes, ok := perShard[s]
		if !ok {
			nodes = internal.MakeMap[string, ShardStatus](len(shard))
			perShard[s] = nodes
		}
		for n, node := range nodes {
			nodes[n] = ShardStatus{
				Status:          Status(node.Status),
				ObjectsExported: node.ObjectsExported,
				SkipReason:      node.SkipReason,
				Error:           node.Error,
			}
		}
	}
	return &Info{
		Backend:             resp.Backend,
		ID:                  resp.ID,
		Path:                resp.Path,
		IncludesCollections: resp.IncludesCollections,
		Status:              Status(resp.Status),
		PerShardStatus:      perShard,
		StartedAt:           resp.StartedAt,
		CompletedAt:         resp.CompletedAt,
		Error:               resp.Error,
		Took:                resp.Took,
		c:                   c,
	}, nil
}

type CancelOptions struct {
	Backend string
	ID      string
}

// Cancel an in-progress export operation.
func (c *Client) CancelCreate(ctx context.Context, options CancelOptions) error {
	req := api.CancelExportRequest{
		Backend: options.Backend,
		ID:      options.ID,
	}

	if err := c.transport.Do(ctx, &req, nil); err != nil {
		return fmt.Errorf("cancel export: %w", err)
	}
	return nil
}

// IsCompleted returns true if the export operation has completed, successfully or otherwise.
func (i *Info) IsCompleted() bool {
	return i.Status == StatusSuccess ||
		i.Status == StatusFailed ||
		i.Status == StatusCanceled
}
