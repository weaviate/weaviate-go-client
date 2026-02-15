package collections

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport) *Client {
	dev.AssertNotNil(t, "t")
	return &Client{t: t}
}

type Client struct {
	t internal.Transport
}

// WithConsistencyLevel default consistency level for all read / write requests made with this collection handle.
func WithConsistencyLevel(cl types.ConsistencyLevel) HandleOption {
	return func(rd *api.RequestDefaults) {
		dev.AssertNotNil(rd, "rd")
		rd.ConsistencyLevel = api.ConsistencyLevel(cl)
	}
}

// WithConsistencyLevel default tenant for all read / write requests made with this collection handle.
func WithTenant(tenant string) HandleOption {
	return func(rd *api.RequestDefaults) {
		dev.AssertNotNil(rd, "rd")
		rd.Tenant = tenant
	}
}

func (c *Client) Use(collectionName string, options ...HandleOption) *Handle {
	rd := api.RequestDefaults{CollectionName: collectionName}
	for _, opt := range options {
		opt(&rd)
	}
	return newHandle(c.t, rd)
}

type Handle struct {
	transport internal.Transport
	defaults  api.RequestDefaults

	Data  *data.Client
	Query *query.Client
}

func newHandle(t internal.Transport, rd api.RequestDefaults) *Handle {
	dev.AssertNotNil(t, "t")

	return &Handle{
		transport: t,
		defaults:  rd,

		Data:  data.NewClient(t, rd),
		Query: query.NewClient(t, rd),
	}
}

func (h *Handle) CollectionName() string {
	return h.defaults.CollectionName
}

func (h *Handle) ConsistencyLevel() types.ConsistencyLevel {
	return types.ConsistencyLevel(h.defaults.ConsistencyLevel)
}

func (h *Handle) Tenant() string {
	return h.defaults.Tenant
}

// HandleOption configures request defaults for collection handle.
type HandleOption func(*api.RequestDefaults)

// WithOptions returns a new handle with different defaults.
func (h *Handle) WithOptions(options ...HandleOption) *Handle {
	defaults := h.defaults
	for _, opt := range options {
		opt(&defaults)
	}
	return newHandle(h.transport, defaults)
}

// Create new collection in the schema. A collection can be created with just the name.
// To configure the new collection, provide a single instance of CreateOptions as the options argument.
//
// Avoid passing multiple options arguments at once -- only the last one will be applied.
func (c *Client) Create(ctx context.Context, collection Collection) (*Handle, error) {
	req := &api.CreateCollectionRequest{Collection: collectionToAPI(&collection)}

	// No need to read the result of the request, we only need the name to create a handle.
	if err := c.t.Do(ctx, req, nil); err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	return c.Use(collection.Name), nil
}
