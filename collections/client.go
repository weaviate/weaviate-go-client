package collections

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/aggregate"
	"github.com/weaviate/weaviate-go-client/v6/batch"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport/stream"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/tenant"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

type streamingTransport interface {
	internal.Transport // For HTTP / gRPC requests.
	stream.Transport   // For server-side batching.
}

func NewClient(t streamingTransport) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{transport: t}
}

type Client struct {
	transport streamingTransport
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
	return newHandle(c.transport, rd)
}

type Handle struct {
	transport streamingTransport
	defaults  api.RequestDefaults

	Aggregate *aggregate.Client
	Data      *data.Client
	Query     *query.Client
	Tenants   *tenant.Client
}

func newHandle(t streamingTransport, rd api.RequestDefaults) *Handle {
	dev.AssertNotNil(t, "t")

	return &Handle{
		transport: t,
		defaults:  rd,

		Aggregate: aggregate.NewClient(t, rd),
		Data:      data.NewClient(t, rd),
		Query:     query.NewClient(t, rd),
		Tenants:   tenant.NewClient(t, rd.CollectionName),
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

// Objects creates a new iterator for objects in target collection.
// The context ctx will be used for all requests throughtout the
// iterator's lifecycle.
func (h *Handle) Objects(ctx context.Context) *query.ObjectIterator {
	return query.NewObjectIterator(ctx, h.Query)
}

// Count objects in the collection, respecting the tenant if provided.
func (h *Handle) Count(ctx context.Context) (int64, error) {
	req := api.CountObjectsRequest(h.defaults)
	var resp api.CountObjectsResponse
	if err := h.transport.Do(ctx, &req, &resp); err != nil {
		return 0, fmt.Errorf("count objects: %w", err)
	}
	return resp.Int64(), nil
}

// Batch opens a new batch stream. The context will be used throughout
// the whole streaming process and may be used to terminate it abruptly.
// In a normal course of operation, a batch should be closed explicitly.
func (h *Handle) Batch(ctx context.Context, options ...batch.Option) (*batch.Client, error) {
	return batch.NewClient(ctx, h.transport, h.defaults, options...), nil
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
	x, err := collectionToAPI(&collection)
	if err != nil {
		return nil, err
	}

	req := &api.CreateCollectionRequest{Collection: x}

	// No need to read the result of the request, we only need the name to create a handle.
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	return c.Use(collection.Name), nil
}

// GetConfig returns configuration for the collection.
// Returns nil with nil error if collections does not exist.
func (c *Client) GetConfig(ctx context.Context, collectionName string) (*Collection, error) {
	var resp api.Collection
	if err := c.transport.Do(ctx, api.GetCollectionRequest(collectionName), &resp); err != nil {
		return nil, fmt.Errorf("get collection config: %w", err)
	}
	collection, err := collectionFromAPI(&resp)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// List returns configurations for all collections defined in the schema.
func (c *Client) List(ctx context.Context) ([]Collection, error) {
	var resp api.ListCollectionsResponse
	if err := c.transport.Do(ctx, api.ListCollectionsRequest, &resp); err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}

	if len(resp) == 0 {
		return nil, nil
	}

	out := make([]Collection, len(resp))
	for i, c := range resp {
		x, err := collectionFromAPI(&c)
		if err != nil {
			return nil, err
		}
		out[i] = x
	}
	return out, nil
}

// Exists check if collection with this name exists. Always check the returned error,
// as Exists may return false with both nil (collection does not exist) and non-nil
// errors (request failed en route).
func (c *Client) Exists(ctx context.Context, collectionName string) (bool, error) {
	var exists api.ResourceExistsResponse
	if err := c.transport.Do(ctx, api.GetCollectionRequest(collectionName), &exists); err != nil {
		return false, fmt.Errorf("check collection exists: %w", err)
	}
	return exists.Bool(), nil
}

// Delete collection by name. Returns an error if no collection with this name exist.
func (c *Client) Delete(ctx context.Context, collectionName string) error {
	if err := c.transport.Do(ctx, api.DeleteCollectionRequest(collectionName), nil); err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}

// DeleteAll collections in the schema.
func (c *Client) DeleteAll(ctx context.Context) error {
	all, err := c.List(ctx)
	if err != nil {
		return fmt.Errorf("delete all collections: %w", err)
	}
	for _, collection := range all {
		if err := c.Delete(ctx, collection.Name); err != nil {
			return err
		}
	}
	return nil
}
