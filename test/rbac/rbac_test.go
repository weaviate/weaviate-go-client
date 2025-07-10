package rbac

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
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
		rootRole   = "root"
		viewerRole = "viewer"

		rootUser = "adam-the-admin"
		pizza    = "Pizza"
	)

	// mustCreateRole and register a t.Cleanup callback to delete it.
	mustCreateRole := func(tt *testing.T, role rbac.Role) {
		tt.Helper()

		tt.Cleanup(func() {
			err := rolesClient.Deleter().WithName(role.Name).Do(ctx)
			require.NoErrorf(tt, err, "delete role %q", role)

			_, err = rolesClient.Getter().WithName(role.Name).Do(ctx)
			require.Errorf(tt, err, "role %q should not exist after deletion", role)
		})

		err := rolesClient.Creator().WithRole(role).Do(ctx)
		require.NoErrorf(tt, err, "create role %q", role)
	}

	hasPermissions := func(tt *testing.T, role string, permissions rbac.Permission) bool {
		tt.Helper()

		has, err := rolesClient.PermissionChecker().
			WithRole(role).
			WithPermission(permissions).
			Do(ctx)
		require.NoError(tt, err, "has-permissions failed")
		return has
	}

	t.Run("get all roles", func(t *testing.T) {
		all, err := rolesClient.AllGetter().Do(ctx)
		require.NoError(t, err, "fetch all roles")
		require.Lenf(t, all, 3, "wrong number of roles")
		require.Equal(t, all[0].Name, adminRole)
		require.Equal(t, all[1].Name, rootRole)
		require.Equal(t, all[2].Name, viewerRole)
	})

	t.Run("get assigned users", func(t *testing.T) {
		assigned, err := rolesClient.UserAssignmentGetter().WithRole(rootRole).Do(ctx)

		require.NoErrorf(t, err, "get users with role %q", rootRole)
		expected := []rbac.UserAssignment{
			{
				UserID:   rootUser,
				UserType: rbac.UserTypeDBEnv,
			},
		}
		require.ElementsMatchf(t, expected, assigned,
			"%q should be assigned to %q", rootRole, rootUser)
	})

	t.Run("get assigned users (Legacy API)", func(t *testing.T) {
		assigned, err := rolesClient.AssignedUsersGetter().WithRole(rootRole).Do(ctx)

		require.NoErrorf(t, err, "get users with role %q", rootRole)
		require.ElementsMatchf(t, []string{rootUser}, assigned,
			"%q should be assigned to %q", rootRole, rootUser)
	})

	t.Run("create role", func(t *testing.T) {
		roleName := "TestRole"

		mustCreateRole(t, rbac.NewRole(roleName,
			rbac.BackupsPermission{Actions: []string{models.PermissionActionManageBackups}, Collection: pizza},
		))

		_, err := rolesClient.Getter().WithName(roleName).Do(ctx)
		require.NoError(t, err, "role %q should exist after creation", roleName)

		testRole, err := rolesClient.Getter().WithName(roleName).Do(ctx)
		require.NoErrorf(t, err, "retrieve %q", roleName)

		require.Equal(t, testRole.Name, roleName)
		require.Len(t, testRole.Backups, 1)
	})

	t.Run("add permissions", func(t *testing.T) {
		roleName := "WantsMorePermissions"
		addPerm := rbac.TenantsPermission{Actions: []string{models.PermissionActionDeleteTenants}}

		mustCreateRole(t, rbac.NewRole(roleName,
			rbac.BackupsPermission{Actions: []string{models.PermissionActionManageBackups}, Collection: pizza},
		))

		err := rolesClient.PermissionAdder().
			WithRole(roleName).
			WithPermissions(addPerm).
			Do(ctx)
		require.NoErrorf(t, err, "add %q permission to %q",
			models.PermissionActionDeleteTenants, roleName)

		require.True(t, hasPermissions(t, roleName, addPerm),
			"%q role should have %q permission", roleName, models.PermissionActionDeleteTenants)
	})

	t.Run("remove permissions", func(t *testing.T) {
		roleName := "WantsLessPermissions"
		removePerm := rbac.TenantsPermission{Actions: []string{models.PermissionActionDeleteTenants}}

		mustCreateRole(t, rbac.NewRole(roleName,
			rbac.BackupsPermission{Actions: []string{models.PermissionActionManageBackups}, Collection: pizza},
		))

		err := rolesClient.PermissionRemover().
			WithRole(roleName).
			WithPermissions(removePerm).
			Do(ctx)
		require.NoErrorf(t, err, "remove %q permission from %q",
			models.PermissionActionDeleteTenants, roleName)

		require.Falsef(t, hasPermissions(t, roleName, removePerm),
			"%q role should not have %q permission", roleName, models.PermissionActionDeleteTenants)
	})
}
