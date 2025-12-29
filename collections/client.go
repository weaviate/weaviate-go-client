package collections

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

func NewClient(t internal.Transport) *Client {
	return &Client{transport: t}
}

type Client struct {
	transport internal.Transport
}

func (c *Client) Use(collectionName string) *Handle {
	return newHandle(c.transport, collectionName)
}

type Handle struct {
	transport      internal.Transport
	collectionName string

	Query *query.Client
	Data  *data.Client
}

func newHandle(t internal.Transport, collectionName string) *Handle {
	return &Handle{
		transport:      t,
		collectionName: collectionName,
		Query:          query.NewClient(t, collectionName),
		Data:           data.NewClient(t, collectionName),
	}
}

type createCollectionRequest struct{ Collection }

type CreateCollectionOption func(*createCollectionRequest)

func WithProperties(properties ...Property) CreateCollectionOption {
	return func(r *createCollectionRequest) {
	}
}

func (c *Client) Create(ctx context.Context, collectionName string, options ...CreateCollectionOption) (*Handle, error) {
	return c.Use(collectionName), nil
}

func (c *Client) GetConfig(ctx context.Context, collectionName string) (*Collection, error) {
	return nil, nil
}

func (c *Client) List(ctx context.Context) ([]Collection, error) {
	return nil, nil
}

func (c *Client) Exists(ctx context.Context) (bool, error) {
	return true, nil
}

func (c *Client) Delete(ctx context.Context, collectionName string) error {
	return nil
}

func (c *Client) DeleteAll(ctx context.Context) error {
	return nil
}
