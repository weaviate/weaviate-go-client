package collections

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport) *Client {
	return &Client{transport: t}
}

type Client struct {
	transport internal.Transport
}

func WithConsistencyLevel(cl types.ConsistencyLevel) HandleOption {
	return func(rd *internal.RequestDefaults) {
		rd.ConsistencyLevel = cl
	}
}

func WithTenant(tenant string) HandleOption {
	return func(rd *internal.RequestDefaults) {
		rd.Tenant = tenant
	}
}

func (c *Client) Use(collectionName string, options ...HandleOption) *Handle {
	rd := internal.RequestDefaults{CollectionName: collectionName}
	for _, opt := range options {
		opt(&rd)
	}
	return newHandle(c.transport, rd)
}

func (c *Client) Create(ctx context.Context, collectionName string, options ...any) *Handle {
	return c.Use(collectionName)
}

type Handle struct {
	transport internal.Transport
	defaults  internal.RequestDefaults

	Query *query.Client
	Data  *data.Client
}

func newHandle(t internal.Transport, rd internal.RequestDefaults) *Handle {
	return &Handle{
		transport: t,
		defaults:  rd,
		Query:     query.NewClient(t, rd),
		Data:      data.NewClient(t, rd),
	}
}

type HandleOption func(*internal.RequestDefaults)

// WithOptions returns a new handle with different defaults.
func (h *Handle) WithOptions(options ...HandleOption) *Handle {
	defaults := h.defaults
	for _, opt := range options {
		opt(&defaults)
	}
	return newHandle(h.transport, defaults)
}
