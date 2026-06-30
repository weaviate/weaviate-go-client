package alias_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/alias"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		alias.NewClient(nil)
	}, "nil transport")
}

func TestClient_Create(t *testing.T) {
	for _, tt := range []struct {
		name  string
		alias collections.Alias
		stubs []testkit.Stub[api.CreateAliasRequest, any]
		err   testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			alias: collections.Alias{
				Collection: "GeorgeBarnes",
				Alias:      "MachineGunKelly",
			},
			stubs: []testkit.Stub[api.CreateAliasRequest, any]{
				{
					Request: &api.CreateAliasRequest{
						Alias: api.Alias{
							Collection: "GeorgeBarnes",
							Alias:      "MachineGunKelly",
						},
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.CreateAliasRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := alias.NewClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.Create(t.Context(), tt.alias)
			tt.err.Require(t, err, "create error")
		})
	}
}

func TestClient_Get(t *testing.T) {
	for _, tt := range []struct {
		name  string
		alias string
		stubs []testkit.Stub[any, api.Alias]
		want  *collections.Alias
		err   testkit.Error // Expected error.
	}{
		{
			name:  "successfully",
			alias: "MachineGunKelly",
			stubs: []testkit.Stub[any, api.Alias]{
				{
					Request: testkit.Ptr(api.GetAliasRequest("MachineGunKelly")),
					Response: api.Alias{
						Collection: "GeorgeBarnes",
						Alias:      "MachineGunKelly",
					},
				},
			},
			want: &collections.Alias{
				Collection: "GeorgeBarnes",
				Alias:      "MachineGunKelly",
			},
		},
		{
			name:  "with error",
			alias: "MachineGunKelly",
			stubs: []testkit.Stub[any, api.Alias]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := alias.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Get(t.Context(), tt.alias)
			tt.err.Require(t, err, "get error")
			require.EqualExportedValues(t, tt.want, got, "returned alias")
		})
	}
}

func TestClient_List(t *testing.T) {
	for _, tt := range []struct {
		name  string
		alias string
		stubs []testkit.Stub[any, api.ListAliasesResponse]
		want  []collections.Alias
		err   testkit.Error // Expected error.
	}{
		{
			name:  "successfully",
			alias: "MachineGunKelly",
			stubs: []testkit.Stub[any, api.ListAliasesResponse]{
				{
					Request: testkit.Ptr[any](api.ListAliasesRequest),
					Response: api.ListAliasesResponse{
						{Collection: "GeorgeBarnes", Alias: "MachineGunKelly"},
						{Collection: "LouisAttanasio", Alias: "LouieHaHa"},
					},
				},
			},
			want: []collections.Alias{
				{Collection: "GeorgeBarnes", Alias: "MachineGunKelly"},
				{Collection: "LouisAttanasio", Alias: "LouieHaHa"},
			},
		},
		{
			name:  "with error",
			alias: "MachineGunKelly",
			stubs: []testkit.Stub[any, api.ListAliasesResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := alias.NewClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.List(t.Context())
			tt.err.Require(t, err, "list error")
			require.EqualExportedValues(t, tt.want, got, "returned aliases")
		})
	}
}

func TestClient_Update(t *testing.T) {
	for _, tt := range []struct {
		name  string
		alias collections.Alias
		stubs []testkit.Stub[api.UpdateAliasRequest, any]
		err   testkit.Error // Expected error.
	}{
		{
			name: "successfully",
			alias: collections.Alias{
				Collection: "GeorgeBarnes",
				Alias:      "MachineGunKelly",
			},
			stubs: []testkit.Stub[api.UpdateAliasRequest, any]{
				{
					Request: &api.UpdateAliasRequest{
						Alias: api.Alias{
							Collection: "GeorgeBarnes",
							Alias:      "MachineGunKelly",
						},
					},
				},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.UpdateAliasRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := alias.NewClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.Update(t.Context(), tt.alias)
			tt.err.Require(t, err, "update error")
		})
	}
}

func TestClient_Delete(t *testing.T) {
	for _, tt := range []struct {
		name  string
		alias string
		stubs []testkit.Stub[any, any]
		err   testkit.Error // Expected error.
	}{
		{
			name:  "successfully",
			alias: "MachineGunKelly",
			stubs: []testkit.Stub[any, any]{
				{
					Request: testkit.Ptr(api.DeleteAliasRequest("MachineGunKelly")),
				},
			},
		},
		{
			name:  "with error",
			alias: "MachineGunKelly",
			stubs: []testkit.Stub[any, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := alias.NewClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.Delete(t.Context(), tt.alias)
			tt.err.Require(t, err, "delete error")
		})
	}
}
