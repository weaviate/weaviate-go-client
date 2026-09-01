package rbac_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/rbac"
)

func TestNewUsersClient(t *testing.T) {
	require.Panics(t, func() {
		rbac.NewUsersClient(nil)
	}, "nil transport")

	t.Run("namespaces", func(t *testing.T) {
		c := rbac.NewUsersClient(testkit.NopTransport)

		assert.NotNil(t, c, "nil client")
		assert.NotNil(t, c.DB, "nil DB namespace")
		assert.NotNil(t, c.OIDC, "nil OIDC namespace")
	})
}

func TestUsersClient_MyUserInfo(t *testing.T) {
	for _, tt := range []struct {
		name  string
		stubs []testkit.Stub[any, api.GetOwnUserInfoResponse]
		want  *rbac.MyUserInfo
		err   testkit.Error
	}{
		{
			name: "ok",
			stubs: []testkit.Stub[any, api.GetOwnUserInfoResponse]{{
				Request: testkit.Ptr[any](api.GetOwnUserInfoRequest),
				Response: api.GetOwnUserInfoResponse{
					ID: "john-malkovich",
					Roles: []api.Role{
						{ID: "rock-n-role"},
						{ID: "sushi-role"},
					},
					Groups: []string{"external", "internal"},
				},
			}},
			want: &rbac.MyUserInfo{
				ID: "john-malkovich",
				Roles: []rbac.Role{
					{ID: "rock-n-role"},
					{ID: "sushi-role"},
				},
				Groups: []string{"external", "internal"},
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, api.GetOwnUserInfoResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.MyUserInfo(t.Context())
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestUsersClient_DB_AssignedRoles(t *testing.T) {
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
				ID:                 "john-malkovich",
				IncludePermissions: false,
			},
			stubs: []testkit.Stub[api.GetAssignedRolesRequest, []api.Role]{{
				Request: &api.GetAssignedRolesRequest{
					Kind:               api.RBACKindDB,
					Entity:             api.UserID("john-malkovich"),
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
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			got, err := c.DB.AssignedRoles(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestUsersClient_DB_AssignRoles(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.AssignRolesOptions
		stubs   []testkit.Stub[api.ManageRolesRequest, any]
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.AssignRolesOptions{
				ID:    "john-malkovich",
				Roles: []string{"rock-n-role", "sushi-role"},
			},
			stubs: []testkit.Stub[api.ManageRolesRequest, any]{{
				Request: &api.ManageRolesRequest{
					Kind:   api.RBACKindDB,
					Entity: api.UserID("john-malkovich"),
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
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			err := c.DB.AssignRoles(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestUsersClient_DB_RevokeRoles(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.RevokeRolesOptions
		stubs   []testkit.Stub[api.ManageRolesRequest, any]
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.RevokeRolesOptions{
				ID:    "john-malkovich",
				Roles: []string{"rock-n-role", "sushi-role"},
			},
			stubs: []testkit.Stub[api.ManageRolesRequest, any]{{
				Request: &api.ManageRolesRequest{
					Kind:   api.RBACKindDB,
					Entity: api.UserID("john-malkovich"),
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
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			err := c.DB.RevokeRoles(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestUsersClient_DB_Create(t *testing.T) {
	for _, tt := range []struct {
		name   string
		userID string
		stubs  []testkit.Stub[any, api.CreateUserResponse]
		want   string
		err    testkit.Error
	}{
		{
			name:   "ok",
			userID: "john-malkovich",
			stubs: []testkit.Stub[any, api.CreateUserResponse]{{
				Request:  testkit.Ptr(api.CreateUserRequest("john-malkovich")),
				Response: api.CreateUserResponse{APIKey: "abracadabra"},
			}},
			want: "abracadabra",
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, api.CreateUserResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			got, err := c.DB.Create(t.Context(), tt.userID)
			tt.err.Require(t, err, "request error")
			assert.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestUsersClient_DB_Delete(t *testing.T) {
	for _, tt := range []struct {
		name   string
		userID string
		stubs  []testkit.Stub[any, any]
		err    testkit.Error
	}{
		{
			name:   "ok",
			userID: "john-malkovich",
			stubs: []testkit.Stub[any, any]{{
				Request: testkit.Ptr(api.DeleteUserRequest("john-malkovich")),
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			err := c.DB.Delete(t.Context(), tt.userID)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestUsersClient_DB_Activate(t *testing.T) {
	for _, tt := range []struct {
		name   string
		userID string
		stubs  []testkit.Stub[any, any]
		err    testkit.Error
	}{
		{
			name:   "ok",
			userID: "john-malkovich",
			stubs: []testkit.Stub[any, any]{{
				Request: testkit.Ptr(api.ActivateUserRequest("john-malkovich")),
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			err := c.DB.Activate(t.Context(), tt.userID)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestUsersClient_DB_Deactivate(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.DeactivateUserOptions
		stubs   []testkit.Stub[api.DeactivateUserRequest, any]
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.DeactivateUserOptions{
				ID:        "john-malkovich",
				RevokeKey: true,
			},
			stubs: []testkit.Stub[api.DeactivateUserRequest, any]{{
				Request: &api.DeactivateUserRequest{
					ID:        "john-malkovich",
					RevokeKey: true,
				},
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.DeactivateUserRequest, any]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			err := c.DB.Deactivate(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
		})
	}
}

func TestUsersClient_DB_RotateKey(t *testing.T) {
	for _, tt := range []struct {
		name   string
		userID string
		stubs  []testkit.Stub[any, api.RotateAPIKeyResponse]
		want   string
		err    testkit.Error
	}{
		{
			name:   "ok",
			userID: "john-malkovich",
			stubs: []testkit.Stub[any, api.RotateAPIKeyResponse]{{
				Request:  testkit.Ptr(api.RotateAPIKeyRequest("john-malkovich")),
				Response: api.RotateAPIKeyResponse{APIKey: "abracadabra"},
			}},
			want: "abracadabra",
		},
		{
			name: "with error",
			stubs: []testkit.Stub[any, api.RotateAPIKeyResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")
			require.NotNil(t, c.DB, "nil DB namespace")

			got, err := c.DB.RotateKey(t.Context(), tt.userID)
			tt.err.Require(t, err, "request error")
			assert.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestUsersClient_DB_UserInfo(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.UserInfoOptions
		stubs   []testkit.Stub[api.GetUserInfoRequest, api.UserInfo]
		want    *rbac.UserInfo
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.UserInfoOptions{
				ID:                "john-malkovich",
				IncludeLastUsedAt: true,
			},
			stubs: []testkit.Stub[api.GetUserInfoRequest, api.UserInfo]{{
				Request: &api.GetUserInfoRequest{
					ID:                "john-malkovich",
					IncludeLastUsedAt: true,
				},
				Response: api.UserInfo{
					ID:                 "john-malkovich",
					Type:               api.UserTypeDBEnv,
					Roles:              []string{"rock-n-role", "sushi-role"},
					Active:             true,
					Namespace:          "sandbox",
					APIKeyFirstLetters: "abr",
					CreatedAt:          testkit.Ptr(testkit.Now),
					LastUsedAt:         testkit.Ptr(testkit.Now),
				},
			}},
			want: &rbac.UserInfo{
				ID:                 "john-malkovich",
				Type:               rbac.UserTypeDBEnv,
				Roles:              []string{"rock-n-role", "sushi-role"},
				Active:             true,
				Namespace:          "sandbox",
				APIKeyFirstLetters: "abr",
				CreatedAt:          testkit.Ptr(testkit.Now),
				LastUsedAt:         testkit.Ptr(testkit.Now),
			},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.GetUserInfoRequest, api.UserInfo]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.DB.UserInfo(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}

func TestUsersClient_DB_List(t *testing.T) {
	for _, tt := range []struct {
		name    string
		options rbac.ListUsersOptions
		stubs   []testkit.Stub[api.ListUsersRequest, []api.UserInfo]
		want    []rbac.UserInfo
		err     testkit.Error
	}{
		{
			name: "ok",
			options: rbac.ListUsersOptions{
				IncludeLastUsedAt: true,
			},
			stubs: []testkit.Stub[api.ListUsersRequest, []api.UserInfo]{{
				Request: &api.ListUsersRequest{
					IncludeLastUsedAt: true,
				},
				Response: []api.UserInfo{{
					ID:                 "john-malkovich",
					Type:               api.UserTypeDBEnv,
					Roles:              []string{"rock-n-role", "sushi-role"},
					Active:             true,
					Namespace:          "sandbox",
					APIKeyFirstLetters: "abr",
					CreatedAt:          testkit.Ptr(testkit.Now),
					LastUsedAt:         testkit.Ptr(testkit.Now),
				}},
			}},
			want: []rbac.UserInfo{{
				ID:                 "john-malkovich",
				Type:               rbac.UserTypeDBEnv,
				Roles:              []string{"rock-n-role", "sushi-role"},
				Active:             true,
				Namespace:          "sandbox",
				APIKeyFirstLetters: "abr",
				CreatedAt:          testkit.Ptr(testkit.Now),
				LastUsedAt:         testkit.Ptr(testkit.Now),
			}},
		},
		{
			name: "with error",
			stubs: []testkit.Stub[api.ListUsersRequest, []api.UserInfo]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := rbac.NewUsersClient(transport)
			require.NotNil(t, c, "nil client")

			got, err := c.DB.List(t.Context(), tt.options)
			tt.err.Require(t, err, "request error")
			require.Equal(t, tt.want, got, "bad result")
		})
	}
}
