package alias

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type Alias api.Alias

type Client struct {
	transport internal.Transport
}

func NewClient(t internal.Transport) *Client {
	return &Client{transport: t}
}

func (c *Client) Delete(ctx context.Context, alias string) error {
	req := api.DeleteAliasRequest(alias)
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("delete alias: %w", err)
	}
	return nil
}

func (c *Client) Create(ctx context.Context, alias, collection string) (*Alias, error) {
	req := &api.CreateAliasRequest{Alias: alias, Collection: collection}
	var resp api.Alias
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("create alias: %w", err)
	}
	return (*Alias)(&resp), nil
}
