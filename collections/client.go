package collections

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/request"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport) *Client {
	return &Client{t: t}
}

type Client struct {
	t internal.Transport
}

// WithConsistencyLevel default consistency level for all read / write requests made with this collection handle.
func WithConsistencyLevel(cl types.ConsistencyLevel) HandleOption {
	return func(rd *request.Defaults) {
		rd.ConsistencyLevel = cl
	}
}

// WithConsistencyLevel default tenant for all read / write requests made with this collection handle.
func WithTenant(tenant string) HandleOption {
	return func(rd *request.Defaults) {
		rd.Tenant = tenant
	}
}

func (c *Client) Use(collectionName string, options ...HandleOption) *Handle {
	rd := request.Defaults{CollectionName: collectionName}
	for _, opt := range options {
		opt(&rd)
	}
	return newHandle(c.t, rd)
}

func (c *Client) Create(ctx context.Context, collectionName string, options ...any) (*Handle, error) {
	if err := c.t.Do(ctx, nil, nil); err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	return c.Use(collectionName), nil
}

type Handle struct {
	transport internal.Transport
	defaults  request.Defaults

	Query *query.Client
	Data  *data.Client
}

func newHandle(t internal.Transport, rd request.Defaults) *Handle {
	return &Handle{
		transport: t,
		defaults:  rd,
		Query:     query.NewClient(t, rd),
		Data:      data.NewClient(t, rd),
	}
}

// HandleOption configures request defaults for collection handle.
type HandleOption func(*request.Defaults)

// WithOptions returns a new handle with different defaults.
func (h *Handle) WithOptions(options ...HandleOption) *Handle {
	defaults := h.defaults
	for _, opt := range options {
		opt(&defaults)
	}
	return newHandle(h.transport, defaults)
}
