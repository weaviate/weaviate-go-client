package collections

import (
	"context"

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

func (c *Client) Create(ctx context.Context, collectionName string, options ...any) *Handle {
	return c.Use(collectionName)
}

type Handle struct {
	transport      internal.Transport
	collectionName string

	Query *query.Client
}

func newHandle(t internal.Transport, collectionName string) *Handle {
	return &Handle{
		transport:      t,
		collectionName: collectionName,
		Query:          query.NewClient(t, collectionName),
	}
}
