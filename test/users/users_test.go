package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestUsers_integration(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	rolesClient := client.Roles()
	usersClient := client.Users()

	const (
		adminRole  = "admin"
		viewerRole = "viewer"

		adminUser = "adam-the-admin"
		pizza     = "Pizza"
	)

	// mustCreateRole and register a t.Cleanup callback to delete it.
	mustCreateRole := func(tt *testing.T, role rbac.Role) {
		tt.Helper()

		tt.Cleanup(func() {
			err := rolesClient.Deleter().WithName(role.Name).Do(ctx)
			require.NoErrorf(tt, err, "delete role %q", role)

			exists, _ := rolesClient.Exists().WithName(role.Name).Do(ctx)
			require.Falsef(tt, exists, "role %q should not exist after deletion", role)
		})

		err := rolesClient.Creator().WithRole(role).Do(ctx)
		require.NoErrorf(tt, err, "create role %q", role.Name)
	}

	t.Run("get user roles", func(t *testing.T) {
		adminRoles, err := usersClient.UserRolesGetter().WithUser(adminUser).Do(ctx)
		require.NoErrorf(t, err, "fetch roles for %q user", adminUser)
		require.Lenf(t, adminRoles, 1, "wrong number of roles for %q user")

		userInfo, err := usersClient.MyUserGetter().Do(ctx)
		require.NoError(t, err, "fetch roles for current user")
		require.Lenf(t, userInfo.Roles, 1, "wrong number of roles for %q user")

		require.EqualExportedValues(t, userInfo.Roles, adminRoles,
			"expect same set of roles for both requests")
	})

	t.Run("assign and revoke a role", func(t *testing.T) {
		roleName := "AssignRevokeMe"

		mustCreateRole(t, rbac.NewRole(roleName,
			rbac.BackupsPermission{Actions: []string{models.PermissionActionManageBackups}, Collection: pizza},
		))

		// Act: assign
		err := usersClient.Assigner().WithUser(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "assign %q role", roleName)

		assignedUsers, _ := rolesClient.AssignedUsersGetter().WithRole(roleName).Do(ctx)
		require.Containsf(t, assignedUsers, adminUser, "should have %q role", roleName)

		// Act: revoke
		err = usersClient.Revoker().WithUser(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "revoke %q role", roleName)

		assignedUsers, _ = rolesClient.AssignedUsersGetter().WithRole(roleName).Do(ctx)
		require.NotContainsf(t, assignedUsers, adminUser, "should not have %q role", roleName)
	})
}
