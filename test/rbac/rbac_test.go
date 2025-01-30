package rbac

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestRBAC_integration(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	rolesClient := client.Roles()

	const (
		adminRole  = "admin"
		viewerRole = "viewer"

		adminUser = "adam-the-admin"
	)

	pizza := "Pizza"
	manageBackups := models.PermissionActionManageBackups
	deleteTenants := models.PermissionActionDeleteTenants

	t.Run("get all roles", func(t *testing.T) {
		all, err := rolesClient.AllGetter().Do(ctx)
		require.NoError(t, err, "fetch all roles")
		require.Lenf(t, all, 2, "wrong number of roles")
		require.Equal(t, *all[0].Name, adminRole)
		require.Equal(t, *all[1].Name, viewerRole)
	})
	t.Run("get user roles", func(t *testing.T) {
		adminRoles, err := rolesClient.UserRolesGetter().WithUser(adminUser).Do(ctx)
		require.NoErrorf(t, err, "fetch roles for %q user", adminUser)
		require.Lenf(t, adminRoles, 1, "wrong number of roles for %q user")

		ownRoles, err := rolesClient.UserRolesGetter().Do(ctx)
		require.NoError(t, err, "fetch roles for current user")
		require.Lenf(t, ownRoles, 1, "wrong number of roles for %q user")

		require.EqualExportedValues(t, ownRoles, adminRoles, "expect same set of roles for both requests")
	})
	t.Run("get assigned users", func(t *testing.T) {
		assigned, err := rolesClient.AssignedUsersGetter().WithRole(adminRole).Do(ctx)

		require.NoErrorf(t, err, "get users with role %q", adminRole)
		require.ElementsMatchf(t, []string{adminUser}, assigned, "only %q should be assigned to %q", adminUser, adminRole)
	})
	t.Run("create role", func(t *testing.T) {
		roleName := "TestRole"
		t.Cleanup(func() {
			err := rolesClient.Deleter().WithName(roleName).Do(ctx)
			require.NoErrorf(t, err, "delete role %q", roleName)

			exists, _ := rolesClient.Exists().WithName(roleName).Do(ctx)
			require.Falsef(t, exists, "role %q should not exist after deletion", roleName)
		})

		err := rolesClient.Creator().
			WithName(roleName).
			WithPermissions(&models.Permission{
				Action:  &manageBackups,
				Backups: &models.PermissionBackups{Collection: &pizza}}).
			Do(ctx)
		require.NoErrorf(t, err, "create role %q", roleName)

		exists, err := rolesClient.Exists().WithName(roleName).Do(ctx)
		require.NoError(t, err, "check if role exists")
		require.Truef(t, exists, "role %q should exist after creation", roleName)

		testRole, err := rolesClient.Getter().WithName(roleName).Do(ctx)
		require.NoErrorf(t, err, "retrieve %q", roleName)

		require.Equal(t, *testRole.Name, roleName)
		require.Len(t, testRole.Permissions, 1)
	})
	t.Run("add permissions", func(t *testing.T) {
		roleName := "WantsMorePermissions"
		addPerm := models.Permission{
			Action: &deleteTenants,
			Tenants: &models.PermissionTenants{
				Collection: &pizza,
			},
		}

		{
			err := rolesClient.Creator().
				WithName(roleName).
				WithPermissions(&models.Permission{
					Action:  &manageBackups,
					Backups: &models.PermissionBackups{Collection: &pizza}}).
				Do(ctx)
			require.NoErrorf(t, err, "create role %q", roleName)
		}

		err := rolesClient.PermissionAdder().
			WithRole(roleName).
			WithPermissions(&addPerm).
			Do(ctx)
		require.NoErrorf(t, err, "add %q permission to %q", deleteTenants, roleName)

		has, err := rolesClient.PermissionChecker().
			WithRole(roleName).
			WithPermissions(&addPerm).
			Do(ctx)
		require.NoError(t, err, "has-permissions failed")
		require.True(t, has, "%q role should have %q permission", roleName, deleteTenants)
	})
	t.Run("remove permissions", func(t *testing.T) {
		roleName := "WantsLessPermissions"
		removePerm := models.Permission{
			Action: &deleteTenants,
			Tenants: &models.PermissionTenants{
				Collection: &pizza,
			},
		}

		{
			err := rolesClient.Creator().
				WithName(roleName).
				// Create an extra permission so that the role would not be
				// deleted with its otherwise only permission is removed.
				WithPermissions(&removePerm, &models.Permission{
					Action:  &manageBackups,
					Backups: &models.PermissionBackups{Collection: &pizza}}).
				Do(ctx)
			require.NoErrorf(t, err, "create role %q", roleName)
		}

		err := rolesClient.PermissionRemover().
			WithRole(roleName).
			WithPermissions(&removePerm).
			Do(ctx)
		require.NoErrorf(t, err, "remove %q permission from %q", deleteTenants, roleName)

		has, err := rolesClient.PermissionChecker().
			WithRole(roleName).
			WithPermissions(&removePerm).
			Do(ctx)
		require.NoError(t, err, "has-permissions failed")
		require.Falsef(t, has, "%q role should not have %q permission", roleName, deleteTenants)

	})
	t.Run("assign and revoke a role", func(t *testing.T) {
		roleName := "AssignRevokeMe"

		{
			err := rolesClient.Creator().
				WithName(roleName).
				// Create an extra permission so that the role would not be
				// deleted with its otherwise only permission is removed.
				WithPermissions(&models.Permission{
					Action:  &manageBackups,
					Backups: &models.PermissionBackups{Collection: &pizza}}).
				Do(ctx)
			require.NoErrorf(t, err, "create role %q", roleName)
		}

		// Act: assign
		err := rolesClient.Assigner().WithUser(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "assign %q role", roleName)

		assignedUsers, _ := rolesClient.AssignedUsersGetter().WithRole(roleName).Do(ctx)
		require.Containsf(t, assignedUsers, adminUser, "should have %q role", roleName)

		// Act: revoke
		err = rolesClient.Revoker().WithUser(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "revoke %q role", roleName)

		assignedUsers, _ = rolesClient.AssignedUsersGetter().WithRole(roleName).Do(ctx)
		require.NotContainsf(t, assignedUsers, adminUser, "should not have %q role", roleName)
	})
}
