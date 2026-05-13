package rbac_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/rbac"
)

func TestNewRolesClient(t *testing.T) {
	require.Panics(t, func() {
		rbac.NewRolesClient(nil)
	}, "nil transport")
}

func TestRolesClient_Create(t *testing.T) {
	for _, tt := range []struct {
		name  string
		role  rbac.Role
		stubs []testkit.Stub[api.CreateRoleRequest, any]
		err   testkit.Error
	}{
		{
			name: "ok",
			role: rbac.Role{
				ID: "rock-n-role",
				Permissions: rbac.Permissions{
					Cluster: []rbac.ClusterPermission{
						{Read: true},
					},
				},
			},
			stubs: []testkit.Stub[api.CreateRoleRequest, any]{{
				Request: &api.CreateRoleRequest{
					Role: api.Role{
						ID: "rock-n-role",
						Permissions: api.Permissions{
							Cluster: []api.ClusterPermission{
								{Read: true},
							},
						},
					},
				},
			}},
		},
		{
			name: "with error",
			role: rbac.Role{ID: "empty-role"},
			stubs: []testkit.Stub[api.CreateRoleRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.Create(t.Context(), tt.role)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestRolesClient_Exists(t *testing.T) {
	for _, tt := range []struct {
		name   string
		roleID string
		stubs  []testkit.Stub[any, api.ResourceExistsResponse]
		want   bool
		err    testkit.Error
	}{
		{
			name:   "role exists",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, api.ResourceExistsResponse]{{
				Request:  testkit.Ptr(api.GetRoleRequest("rock-n-role")),
				Response: true,
			}},
			want: true,
		},
		{
			name:   "role not exists",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, api.ResourceExistsResponse]{{
				Request:  testkit.Ptr(api.GetRoleRequest("rock-n-role")),
				Response: false,
			}},
			want: false,
		},
		{
			name:   "with error",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, api.ResourceExistsResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Exists(t.Context(), tt.roleID)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "exists")
		})
	}
}

func TestRolesClient_Get(t *testing.T) {
	for _, tt := range []struct {
		name   string
		roleID string
		stubs  []testkit.Stub[any, api.Role]
		want   *rbac.Role
		err    testkit.Error
	}{
		{
			name:   "ok",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, api.Role]{{
				Request:  testkit.Ptr(api.GetRoleRequest("rock-n-role")),
				Response: api.Role{ID: "rock-n-role"},
			}},
			want: &rbac.Role{ID: "rock-n-role"},
		},
		{
			name:   "with error",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, api.Role]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.Get(t.Context(), tt.roleID)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestRolesClient_List(t *testing.T) {
	for _, tt := range []struct {
		name   string
		roleID string
		stubs  []testkit.Stub[any, []api.Role]
		want   []rbac.Role
		err    testkit.Error
	}{
		{
			name:   "ok",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, []api.Role]{{
				Request:  testkit.Ptr[any](api.ListRolesRequest),
				Response: []api.Role{{ID: "rock-n-role"}},
			}},
			want: []rbac.Role{{ID: "rock-n-role"}},
		},
		{
			name:   "with error",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, []api.Role]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.List(t.Context())
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestRolesClient_AddPermissions(t *testing.T) {
	for _, tt := range []struct {
		name  string
		add   rbac.AddPermissions
		stubs []testkit.Stub[api.AddPermissionsRequest, any]
		err   testkit.Error
	}{
		{
			name: "ok",
			add: rbac.AddPermissions{
				RoleID: "rock-n-role",
				Permissions: rbac.Permissions{
					Cluster: []rbac.ClusterPermission{
						{Read: true},
					},
				},
			},
			stubs: []testkit.Stub[api.AddPermissionsRequest, any]{{
				Request: &api.AddPermissionsRequest{
					RoleID: "rock-n-role",
					Permissions: api.Permissions{
						Cluster: []api.ClusterPermission{
							{Read: true},
						},
					},
				},
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.AddPermissionsRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.AddPermissions(t.Context(), tt.add)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestRolesClient_RemovePermissions(t *testing.T) {
	for _, tt := range []struct {
		name   string
		remove rbac.RemovePermissions
		stubs  []testkit.Stub[api.RemovePermissionsRequest, any]
		err    testkit.Error
	}{
		{
			name: "ok",
			remove: rbac.RemovePermissions{
				RoleID: "rock-n-role",
				Permissions: rbac.Permissions{
					Cluster: []rbac.ClusterPermission{
						{Read: true},
					},
				},
			},
			stubs: []testkit.Stub[api.RemovePermissionsRequest, any]{{
				Request: &api.RemovePermissionsRequest{
					RoleID: "rock-n-role",
					Permissions: api.Permissions{
						Cluster: []api.ClusterPermission{
							{Read: true},
						},
					},
				},
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.RemovePermissionsRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			err := c.RemovePermissions(t.Context(), tt.remove)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestRolesClient_HasPermission(t *testing.T) {
	for _, tt := range []struct {
		name       string
		permission rbac.HasPermission
		stubs      []testkit.Stub[api.HasPermissionRequest, api.HasPermissionResponse]
		want       bool
		err        testkit.Error
	}{
		{
			name: "has permission",
			permission: rbac.HasPermission{
				RoleID:  "rock-n-role",
				Cluster: rbac.ClusterPermission{Read: true},
			},
			stubs: []testkit.Stub[api.HasPermissionRequest, api.HasPermissionResponse]{{
				Request: &api.HasPermissionRequest{
					RoleID:  "rock-n-role",
					Cluster: api.ClusterPermission{Read: true},
				},
				Response: true,
			}},
			want: true,
		},
		{
			name: "not has permission",
			permission: rbac.HasPermission{
				RoleID:  "rock-n-role",
				Cluster: rbac.ClusterPermission{Read: true},
			},
			stubs: []testkit.Stub[api.HasPermissionRequest, api.HasPermissionResponse]{{
				Request: &api.HasPermissionRequest{
					RoleID:  "rock-n-role",
					Cluster: api.ClusterPermission{Read: true},
				},
				Response: false,
			}},
			want: false,
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.HasPermissionRequest, api.HasPermissionResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.HasPermission(t.Context(), tt.permission)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestRolesClient_AssignedUserIDs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		roleID string
		stubs  []testkit.Stub[any, api.GetAssignedUsersResponse]
		want   []string
		err    testkit.Error
	}{
		{
			name:   "has permission",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, api.GetAssignedUsersResponse]{{
				Request:  testkit.Ptr(api.GetAssignedUsersRequest("rock-n-role")),
				Response: []string{"john_doe", "jane_doe"},
			}},
			want: []string{"john_doe", "jane_doe"},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, api.GetAssignedUsersResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.AssignedUserIDs(t.Context(), tt.roleID)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestRolesClient_UserAssignments(t *testing.T) {
	for _, tt := range []struct {
		name   string
		roleID string
		stubs  []testkit.Stub[any, []api.UserAssignment]
		want   []rbac.UserAssignment
		err    testkit.Error
	}{
		{
			name:   "has permission",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, []api.UserAssignment]{{
				Request: testkit.Ptr(api.GetUserAssignmentsRequest("rock-n-role")),
				Response: []api.UserAssignment{
					{ID: "john_doe", Type: "db_env"},
				},
			}},
			want: []rbac.UserAssignment{
				{ID: "john_doe", Type: "db_env"},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, []api.UserAssignment]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.UserAssignments(t.Context(), tt.roleID)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestRolesClient_GroupAssignments(t *testing.T) {
	for _, tt := range []struct {
		name   string
		roleID string
		stubs  []testkit.Stub[any, []api.GroupAssignment]
		want   []rbac.GroupAssignment
		err    testkit.Error
	}{
		{
			name:   "has permission",
			roleID: "rock-n-role",
			stubs: []testkit.Stub[any, []api.GroupAssignment]{{
				Request: testkit.Ptr(api.GetGroupAssignmentsRequest("rock-n-role")),
				Response: []api.GroupAssignment{
					{ID: "external", Type: "oidc"},
				},
			}},
			want: []rbac.GroupAssignment{
				{ID: "external", Type: "oidc"},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, []api.GroupAssignment]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewRolesClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.GroupAssignments(t.Context(), tt.roleID)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}
