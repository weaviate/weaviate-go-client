package alias

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewClient(t internal.Transport) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{transport: t}
}

type Client struct {
	transport internal.Transport
}

type Alias api.Alias

// Create a new alias for a collection.
func (c *Client) Create(ctx context.Context, a Alias) error {
	req := &api.CreateAliasRequest{
		Alias: api.Alias(a),
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("create alias: %w", err)
	}
	return nil
}

// Get the alias by name if one exists.
func (c *Client) Get(ctx context.Context, alias string) (*Alias, error) {
	var a api.Alias
	if err := c.transport.Do(ctx, api.GetAliasRequest(alias), &a); err != nil {
		return nil, fmt.Errorf("get alias: %w", err)
	}
	return (*Alias)(&a), nil
}

// List all defined aliases.
func (c *Client) List(ctx context.Context) ([]Alias, error) {
	var list api.ListAliasesResponse
	if err := c.transport.Do(ctx, api.ListAliasesRequest, &list); err != nil {
		return nil, fmt.Errorf("list aliases: %w", err)
	}

	aliases := make([]Alias, len(list))
	for i := range list {
		aliases[i] = Alias(list[i])
	}
	return aliases, nil
}

// Update re-assigns the alias to a different collection.
func (c *Client) Update(ctx context.Context, a Alias) error {
	req := &api.UpdateAliasRequest{
		Alias: api.Alias(a),
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("update alias: %w", err)
	}
	return nil
}

// Delete the alias by name.
func (c *Client) Delete(ctx context.Context, alias string) error {
	if err := c.transport.Do(ctx, api.DeleteAliasRequest(alias), nil); err != nil {
		return fmt.Errorf("delete alias: %w", err)
	}
	return nil
}
