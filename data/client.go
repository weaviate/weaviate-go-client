package data

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

func (c *Client) Insert(ctx context.Context, options ...InsertOption) (*types.Object[types.Map], error) {
	var ir insertRequest
	InsertOptions(options).Apply(&ir)

	req := &api.InsertObjectRequest{
		UUID:       ir.UUID,
		Properties: ir.Properties,
		References: nil, // TODO(dyma)
		Vectors:    api.Vectors(ir.Vectors),
	}

	var resp api.InsertObjectResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("insert object: %w", err)
	}

	return &types.Object[types.Map]{
		UUID:       resp.UUID,
		Properties: resp.Properties,
		// References: resp.References,
		Vectors:            types.Vectors(resp.Vectors),
		CreationTimeUnix:   resp.CreationTimeUnix,
		LastUpdateTimeUnix: resp.LastUpdateTimeUnix,
	}, nil
}

func (c *Client) Delete(ctx context.Context, id uuid.UUID) error {
	req := api.DeleteObjectRequest{
		RequestDefaults: c.defaults,
		UUID:            id,
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("delete alias: %w", err)
	}
	return nil
}

type insertRequest struct{ types.Object[types.Properties] }

type InsertOption func(*insertRequest)

type InsertOptions []InsertOption

func (opts InsertOptions) Apply(*insertRequest) {}

func WithUUID(uuid uuid.UUID) InsertOption {
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
