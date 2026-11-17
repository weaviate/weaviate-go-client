package rbac_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/rbac"
)

func TestNewGroupsClient(t *testing.T) {
	require.Panics(t, func() {
		rbac.NewGroupsClient(nil)
	}, "nil transport")
}

func TestGroupsClient_AssignedRoles(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.AssignedRolesOptions
		stubs   []testkit.Stub[api.GetAssignedRolesRequest, []api.Role]
		want    []rbac.Role
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.AssignedRolesOptions{
				ID:                 "external",
				IncludePermissions: false,
			},
			stubs: []testkit.Stub[api.GetAssignedRolesRequest, []api.Role]{{
				Request: &api.GetAssignedRolesRequest{
					Kind:               api.RBACKindOIDC,
					Entity:             api.GroupID("external"),
					IncludePermissions: false,
				},
				Response: []api.Role{
					{ID: "rock-n-role"},
					{ID: "sushi-role"},
				},
			}},
			want: []rbac.Role{
				{ID: "rock-n-role"},
				{ID: "sushi-role"},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.GetAssignedRolesRequest, []api.Role]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewGroupsClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.AssignedRoles(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestGroupsClient_AssignRoles(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.AssignRolesOptions
		stubs   []testkit.Stub[api.ManageRolesRequest, any]
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.AssignRolesOptions{
				ID:    "external",
				Roles: []string{"rock-n-role", "sushi-role"},
			},
			stubs: []testkit.Stub[api.ManageRolesRequest, any]{{
				Request: &api.ManageRolesRequest{
					Kind:   api.RBACKindOIDC,
					Entity: api.GroupID("external"),
					Verb:   api.RoleVerbAssign,
					Roles:  []string{"rock-n-role", "sushi-role"},
				},
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.ManageRolesRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewGroupsClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.AssignRoles(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestGroupsClient_RevokeRoles(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.RevokeRolesOptions
		stubs   []testkit.Stub[api.ManageRolesRequest, any]
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.RevokeRolesOptions{
				ID:    "external",
				Roles: []string{"rock-n-role", "sushi-role"},
			},
			stubs: []testkit.Stub[api.ManageRolesRequest, any]{{
				Request: &api.ManageRolesRequest{
					Kind:   api.RBACKindOIDC,
					Entity: api.GroupID("external"),
					Verb:   api.RoleVerbRevoke,
					Roles:  []string{"rock-n-role", "sushi-role"},
				},
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.ManageRolesRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewGroupsClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.RevokeRoles(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
		})
	}
}
