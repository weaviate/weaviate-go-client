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

// TODO(dyma): generic Properties
type Object struct {
	UUID       uuid.UUID
	Properties any
	Vectors    []types.Vector
}

func (c *Client) Insert(ctx context.Context, ir *Object) (*types.Object[types.Map], error) {
	var req api.InsertObjectRequest
	if ir != nil {
		req = api.InsertObjectRequest{
			UUID:       ir.UUID,
			Properties: ir.Properties,
			Vectors:    newVectors(ir.Vectors),
		}
	}

	var resp api.InsertObjectResponse
	if err := c.transport.Do(ctx, &req, &resp); err != nil {
		return nil, fmt.Errorf("insert object: %w", err)
	}

	return &types.Object[types.Map]{
		UUID:               resp.UUID,
		Properties:         resp.Properties,
		Vectors:            types.Vectors(resp.Vectors),
		CreationTimeUnix:   resp.CreationTimeUnix,
		LastUpdateTimeUnix: resp.LastUpdateTimeUnix,
	}, nil
}

func (c *Client) Replace(ctx context.Context, ir Object) (*types.Object[types.Map], error) {
	req := &api.ReplaceObjectRequest{
		UUID:       ir.UUID,
		Properties: ir.Properties,
		Vectors:    newVectors(ir.Vectors),
	}

	var resp api.ReplaceObjectResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("insert object: %w", err)
	}

	return &types.Object[types.Map]{
		UUID:               resp.UUID,
		Properties:         resp.Properties,
		Vectors:            types.Vectors(resp.Vectors),
		CreationTimeUnix:   resp.CreationTimeUnix,
		LastUpdateTimeUnix: resp.LastUpdateTimeUnix,
	}, nil
}

func newVectors(vectors []types.Vector) api.Vectors {
	vs := make(api.Vectors, len(vectors))
	for _, v := range vectors {
		vs[v.Name] = api.Vector(v)
	}
	return vs
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
