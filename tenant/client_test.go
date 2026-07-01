package tenant_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/tenant"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		tenant.NewClient(nil, "Songs")
	}, "nil transport")
}

func TestClient_Create(t *testing.T) {
	for _, tt := range []struct {
		name    string
		tenants []tenant.Tenant
		stubs   []testkit.Stub[api.CreateTenantsRequest, any]
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			tenants: []tenant.Tenant{
				{Name: "john_doe", Status: tenant.Active},
				{Name: "jane_doe", Status: tenant.Frozen},
			},
			stubs: []testkit.Stub[api.CreateTenantsRequest, any]{
				{
					Request: &api.CreateTenantsRequest{
						Collection: "Songs",
						Tenants: []api.Tenant{
							{Name: "john_doe", Status: api.TenantStatusActive},
							{Name: "jane_doe", Status: api.TenantStatusFrozen},
						},
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.CreateTenantsRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := tenant.NewClient(transport, "Songs")
			require.NotNil(t, c, "nil client")

			err := c.Create(t.Context(), tt.tenants...)
			tt.err.Require(t, err, "create error")
		})
	}
}

func TestClient_Update(t *testing.T) {
	for _, tt := range []struct {
		name    string
		tenants []tenant.Tenant
		stubs   []testkit.Stub[api.UpdateTenantsRequest, any]
		err     testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			tenants: []tenant.Tenant{
				{Name: "john_doe", Status: tenant.Active},
				{Name: "jane_doe", Status: tenant.Frozen},
			},
			stubs: []testkit.Stub[api.UpdateTenantsRequest, any]{
				{
					Request: &api.UpdateTenantsRequest{
						Collection: "Songs",
						Tenants: []api.Tenant{
							{Name: "john_doe", Status: api.TenantStatusActive},
							{Name: "jane_doe", Status: api.TenantStatusFrozen},
						},
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.UpdateTenantsRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := tenant.NewClient(transport, "Songs")
			require.NotNil(t, c, "nil client")

			err := c.Update(t.Context(), tt.tenants...)
			tt.err.Require(t, err, "update error")
		})
	}
}

func TestClient_Get(t *testing.T) {
	for _, tt := range []struct {
		name    string
		tenants []string
		stubs   []testkit.Stub[api.GetTenantsRequest, api.GetTenantsResponse]
		want    []tenant.Tenant
		err     testkit.Error // Expected error.
	}{
		{
			name:    "successfully",
			tenants: []string{"john_doe", "jane_doe"},
			stubs: []testkit.Stub[api.GetTenantsRequest, api.GetTenantsResponse]{
				{
					Request: &api.GetTenantsRequest{
						Collection: "Songs",
						Tenants:    []string{"john_doe", "jane_doe"},
					},
					Response: api.GetTenantsResponse{
						{Name: "john_doe", Status: api.TenantStatusActive},
						{Name: "jane_doe", Status: api.TenantStatusFrozen},
					},
				},
			},
			want: []tenant.Tenant{
				{Name: "john_doe", Status: tenant.Active},
				{Name: "jane_doe", Status: tenant.Frozen},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.GetTenantsRequest, api.GetTenantsResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := tenant.NewClient(transport, "Songs")
			require.NotNil(t, c, "nil client")

			got, err := c.Get(t.Context(), tt.tenants...)
			tt.err.Require(t, err, "get-tenants error")
			require.EqualExportedValues(t, tt.want, got, "returned tenants")
		})
	}
}

func TestClient_Delete(t *testing.T) {
	for _, tt := range []struct {
		name    string
		tenants []string
		stubs   []testkit.Stub[api.DeleteTenantsRequest, any]
		err     testkit.Error // Expected error.
	}{
		{
			name:    "successfully",
			tenants: []string{"john_doe", "jane_doe"},
			stubs: []testkit.Stub[api.DeleteTenantsRequest, any]{
				{
					Request: &api.DeleteTenantsRequest{
						Collection: "Songs",
						Tenants:    []string{"john_doe", "jane_doe"},
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.DeleteTenantsRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := tenant.NewClient(transport, "Songs")
			require.NotNil(t, c, "nil client")

			err := c.Delete(t.Context(), tt.tenants...)
			tt.err.Require(t, err, "delete error")
		})
	}
}
