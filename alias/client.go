package alias

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type Client struct {
	transport api.RESTTransport
}

func NewClient(t api.RESTTransport) *Client {
	return &Client{transport: t}
}

func (c *Client) Delete(ctx context.Context, alias string) error {
	req := api.DeleteAliasRequest(alias)
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("delete alias: %w", err)
	}
	return nil
}
