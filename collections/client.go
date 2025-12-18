package collections

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
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
	return newHandle(c.t, rd)
}

type Handle struct {
	transport internal.Transport
	defaults  api.RequestDefaults

	Query *query.Client
	Data  *data.Client
}

func newHandle(t internal.Transport, rd api.RequestDefaults) *Handle {
	return &Handle{
		transport: t,
		defaults:  rd,
		Query:     query.NewClient(t, rd),
		Data:      data.NewClient(t, rd),
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
	return newHandle(h.transport, defaults)
}

type CreateOptions struct {
	Description   string
	Properties    []Property
	References    []Reference
	Sharding      ShardingConfig
	Replication   ReplicationConfig
	InvertedIndex InvertedIndexConfig
	MultiTenancy  MultiTenancyConfig
}

func (c *Client) Create(ctx context.Context, collectionName string, options ...CreateOptions) (*Handle, error) {
	var collection Collection
	for _, opt := range options {
		collection.Description = opt.Description
		collection.Properties = opt.Properties
		collection.References = opt.References
		collection.Sharding = opt.Sharding
		collection.Replication = opt.Replication
		collection.InvertedIndex = opt.InvertedIndex
		collection.MultiTenancy = opt.MultiTenancy
	}
	collection.Name = collectionName

	req := &api.CreateCollectionRequest{Collection: collection.toAPI()}
	if err := c.t.Do(ctx, req, nil); err != nil {
		return nil, fmt.Errorf("create collection: %w", err)
	}
	return c.Use(collectionName), nil
}

func (c *Client) GetConfig(ctx context.Context, collectionName string) (*Collection, error) {
	if err := c.t.Do(ctx, nil, nil); err != nil {
		return nil, fmt.Errorf("get collection config: %w", err)
	}
	return nil, nil
}

func (c *Client) List(ctx context.Context) ([]Collection, error) {
	if err := c.t.Do(ctx, nil, nil); err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}
	return nil, nil
}

func (c *Client) Exists(ctx context.Context) (bool, error) {
	// TODO(dyma): send the same request as in GetConfig, but pass nil-dest to skip unmarshaling.
	if err := c.t.Do(ctx, nil, nil); err != nil {
		return false, fmt.Errorf("check collection exists: %w", err)
	}
	return true, nil
}

func (c *Client) Delete(ctx context.Context, collectionName string) error {
	if err := c.t.Do(ctx, api.DeleteCollectionRequest(collectionName), nil); err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}

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
