package tenant

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewClient(t internal.Transport, collection string) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{
		transport:  t,
		collection: collection,
	}
}

type Client struct {
	transport  internal.Transport
	collection string
}

type Tenant struct {
	Name   string
	Status Status
}

type Status api.TenantStatus

const (
	Active     = Status(api.TenantStatusActive)
	Cold       = Status(api.TenantStatusCold)
	Freezing   = Status(api.TenantStatusFreezing)
	Frozen     = Status(api.TenantStatusFrozen)
	Hot        = Status(api.TenantStatusHot)
	Inactive   = Status(api.TenantStatusInactive)
	Offloaded  = Status(api.TenantStatusOffloaded)
	Offloading = Status(api.TenantStatusOffloading)
	Onloading  = Status(api.TenantStatusOnloading)
	Unfreezing = Status(api.TenantStatusUnfreezing)
)

// Create new tenants in the collection.
func (c *Client) Create(ctx context.Context, tenants ...Tenant) error {
	req := &api.CreateTenantsRequest{
		Collection: c.collection,
		Tenants:    make([]api.Tenant, len(tenants)),
	}
	for i, t := range tenants {
		req.Tenants[i] = api.Tenant{
			Name:   t.Name,
			Status: api.TenantStatus(t.Status),
		}
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("create tenants: %w", err)
	}
	return nil
}

// TODO(dyma): do we want to hedge for the possibility of some cursor-based API
// and accept an options struct instead?
func (c *Client) Get(ctx context.Context, tenants ...string) ([]Tenant, error) {
	req := &api.GetTenantsRequest{
		Collection: c.collection,
		Tenants:    tenants,
	}
	var resp api.GetTenantsResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("get tenants: %w", err)
	}

	out := make([]Tenant, len(resp))
	for i, t := range resp {
		out[i] = Tenant{
			Name:   t.Name,
			Status: Status(t.Status),
		}
	}
	return out, nil
}

// Update statuses of existing tenants to either [Active], [Inactive], or [Offloaded].
// All other statuses should be considered "read-only".
func (c *Client) Update(ctx context.Context, tenants ...Tenant) error {
	req := &api.UpdateTenantsRequest{
		Collection: c.collection,
		Tenants:    make([]api.Tenant, len(tenants)),
	}
	for i, t := range tenants {
		req.Tenants[i] = api.Tenant{
			Name:   t.Name,
			Status: api.TenantStatus(t.Status),
		}
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("update tenants: %w", err)
	}
	return nil
}

// Activate tenants.
func (c *Client) Activate(ctx context.Context, tenants ...string) error {
	return c.setStatus(ctx, "activate", Active, tenants)
}

// Deactivate tenants.
func (c *Client) Deactivate(ctx context.Context, tenants ...string) error {
	return c.setStatus(ctx, "deactivate", Inactive, tenants)
}

// Offload tenants to the configured backend storage.
func (c *Client) Offload(ctx context.Context, tenants ...string) error {
	return c.setStatus(ctx, "offload", Offloaded, tenants)
}

func (c *Client) setStatus(ctx context.Context, verb string, s Status, tenants []string) error {
	req := &api.UpdateTenantsRequest{
		Collection: c.collection,
		Tenants:    make([]api.Tenant, len(tenants)),
	}
	for i := range tenants {
		req.Tenants[i] = api.Tenant{
			Name:   tenants[i],
			Status: api.TenantStatus(s),
		}
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("%s tenants: %w", verb, err)
	}
	return nil
}

// Delete tenants and their data.
func (c *Client) Delete(ctx context.Context, tenants ...string) error {
	req := &api.DeleteTenantsRequest{
		Collection: c.collection,
		Tenants:    tenants,
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("delete tenants: %w", err)
	}
	return nil
}
