package data

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport, rd api.RequestDefaults) *Client {
	return &Client{
		transport: t,
		defaults:  rd,
	}
}

type Client struct {
	transport internal.Transport
	defaults  api.RequestDefaults
}

func (c *Client) Insert(ctx context.Context, options ...InsertOption) (string, error) {
	var ir insertRequest
	InsertOptions(options).Apply(&ir)

	req := &api.InsertObjectRequest{
		UUID:       ir.UUID,
		Properties: ir.Properties,
		Vectors:    ir.Vectors,
	}

	var resp api.InsertObjectResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return "", fmt.Errorf("insert object: %w", err)
	}

	return resp.UUID, nil
}

type insertRequest struct{ types.Object[types.Properties] }

type InsertOption func(*insertRequest)

type InsertOptions []InsertOption

func (opts InsertOptions) Apply(*insertRequest) {}

func WithUUID(uuid string) InsertOption {
	return func(r *insertRequest) {
		r.Object.UUID = uuid
	}
}

func WithProperties(p types.Properties) InsertOption {
	return func(r *insertRequest) {
		r.Object.Properties = p
	}
}

func WithVector(vectors ...types.Vector) InsertOption {
	return func(r *insertRequest) {
		if r.Object.Vectors == nil {
			r.Object.Vectors = make(map[string]api.Vector, len(vectors))
		}
		for _, v := range vectors {
			r.Object.Vectors[v.Name] = api.Vector(v)
		}
	}
}
