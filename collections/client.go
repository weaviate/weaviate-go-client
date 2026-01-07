package collections

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t api.REST) *Client {
	return &Client{rest: t}
}

type Client struct {
	rest api.REST
	gRPC api.GRPC
}

// WithConsistencyLevel default consistency level for all read / write requests made with this collection handle.
func WithConsistencyLevel(cl types.ConsistencyLevel) HandleOption {
	return func(rd *api.RequestDefaults) {
		rd.ConsistencyLevel = api.ConsistencyLevel(cl)
	}
}

// WithConsistencyLevel default tenant for all read / write requests made with this collection handle.
func WithTenant(tenant string) HandleOption {
	return func(rd *api.RequestDefaults) {
		rd.Tenant = tenant
	}
}

func (c *Client) Use(collectionName string, options ...HandleOption) *Handle {
	rd := api.RequestDefaults{CollectionName: collectionName}
	for _, opt := range options {
		opt(&rd)
	}
	return newHandle(c.rest, c.gRPC, rd)
}

type Handle struct {
	rest     api.REST
	gRPC     api.GRPC
	defaults api.RequestDefaults

	Query *query.Client
	Data  *data.Client
}

func newHandle(r api.REST, g api.GRPC, rd api.RequestDefaults) *Handle {
	return &Handle{
		rest:     r,
		gRPC:     g,
		defaults: rd,
		Query:    query.NewClient(api.SearchTransport(g.Search), rd),
		Data:     data.NewClient(r, g, rd),
	}
}

// HandleOption configures request defaults for collection handle.
type HandleOption func(*api.RequestDefaults)

// WithOptions returns a new handle with different defaults.
func (h *Handle) WithOptions(options ...HandleOption) *Handle {
	defaults := h.defaults
	for _, opt := range options {
		opt(&defaults)
	}
	return newHandle(h.rest, h.gRPC, defaults)
}

// Create new collection in the schema. A collection can be created with just the name.
// To configure the new collection, provide a single instance of CreateOptions as the options argument.
//
// Avoid passing multiple options arguments at once -- only the last one will be applied.
func (c *Client) Create(ctx context.Context, collection Collection) (*Handle, error) {
	req := &api.CreateCollectionRequest{Collection: collection.toAPI()}

	// No need to read the result of the request, we only need the name to create a handle.
	if err := c.rest.Do(ctx, req, nil); err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	return c.Use(collection.Name), nil
}

// GetConfig returns configuration for the collection.
// Returns nil with nil error if collections does not exist.
func (c *Client) GetConfig(ctx context.Context, collectionName string) (*Collection, error) {
	var resp api.Collection
	if err := c.rest.Do(ctx, api.GetCollectionRequest(collectionName), &resp); err != nil {
		return nil, fmt.Errorf("get collection config: %w", err)
	}
	collection := fromAPI(&resp)
	return &collection, nil
}

// List returns configurations for all collections defined in the schema.
func (c *Client) List(ctx context.Context) ([]Collection, error) {
	var resp api.ListCollectionsResponse
	if err := c.rest.Do(ctx, api.ListCollectionsRequest, &resp); err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}

	out := make([]Collection, len(resp))
	for i, c := range resp {
		out[i] = fromAPI(&c)
	}
	return out, nil
}

// Exists check if collection with this name exists. Always check the returned error,
// as Exists may return false with both nil (collection does not exist) and non-nil
// errors (request failed en route).
func (c *Client) Exists(ctx context.Context, collectionName string) (bool, error) {
	var exists api.ResourceExistsResponse
	if err := c.rest.Do(ctx, api.GetCollectionRequest(collectionName), &exists); err != nil {
		return false, fmt.Errorf("check collection exists: %w", err)
	}
	return exists.Bool(), nil
}

// Delete collection by name.
// TODO(dyma): Should we return an error on "not found" or not?
func (c *Client) Delete(ctx context.Context, collectionName string) error {
	if err := c.rest.Do(ctx, api.DeleteCollectionRequest(collectionName), nil); err != nil {
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
