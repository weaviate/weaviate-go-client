package data

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func NewClient(t internal.Transport, rd internal.RequestDefaults) *Client {
	return &Client{
		transport: t,
		defaults:  rd,
	}
}

type Client struct {
	transport internal.Transport
	defaults  internal.RequestDefaults
}

func (c *Client) Insert(context.Context, ...InsertOption) (string, error) {
	return "", nil
}

type insertRequest struct{ types.Object[types.Properties] }

type InsertOption func(*insertRequest)

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
			r.Object.Vectors = make(map[string]types.Vector, len(vectors))
		}
		for _, v := range vectors {
			r.Object.Vectors[v.Name] = v
		}
	}
}
