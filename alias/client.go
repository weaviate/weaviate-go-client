package alias

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

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
