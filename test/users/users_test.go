package users

import (
	"context"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
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

	t.Run("get user roles (legacy API)", func(t *testing.T) {
		adminRoles, err := usersClient.UserRolesGetter().WithUserID(adminUser).WithIncludeFullRoles(true).Do(ctx)
		require.NoErrorf(t, err, "fetch roles for %q user", adminUser)
		require.Lenf(t, adminRoles, 1, "wrong number of roles for %q user")

		userInfo, err := usersClient.MyUserGetter().Do(ctx)
		require.NoError(t, err, "fetch roles for current user")
		require.Equal(t, adminUser, userInfo.UserID)
		require.Lenf(t, userInfo.Roles, 1, "wrong number of roles for %q user")

		require.EqualExportedValues(t, userInfo.Roles, adminRoles,
			"expect same set of roles for both requests")
	})

	t.Run("get user roles", func(t *testing.T) {
		adminRoles, err := usersClient.DB().RolesGetter().WithUserID(adminUser).WithIncludeFullRoles(true).Do(ctx)
		require.NoErrorf(t, err, "fetch roles for %q user", adminUser)
		require.Lenf(t, adminRoles, 1, "wrong number of roles for %q user")

		userInfo, err := usersClient.MyUserGetter().Do(ctx)
		require.NoError(t, err, "fetch roles for current user")
		require.Equal(t, adminUser, userInfo.UserID)
		require.Lenf(t, userInfo.Roles, 1, "wrong number of roles for %q user")

		require.EqualExportedValues(t, userInfo.Roles, adminRoles,
			"expect same set of roles for both requests")
	})

	t.Run("assign and revoke a role (legacy API)", func(t *testing.T) {
		roleName := "AssignRevokeMe"

		mustCreateRole(t, rbac.NewRole(roleName,
			rbac.BackupsPermission{Actions: []string{models.PermissionActionManageBackups}, Collection: pizza},
		))

		// Act: assign
		err := usersClient.Assigner().WithUserID(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "assign %q role", roleName)

		assignedUsers, _ := rolesClient.AssignedUsersGetter().WithRole(roleName).Do(ctx)
		require.Truef(t, slices.Contains(assignedUsers, adminUser), "should have %q role", roleName)

		// Act: revoke
		err = usersClient.Revoker().WithUserID(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "revoke %q role", roleName)

		assignedUsers, _ = rolesClient.AssignedUsersGetter().WithRole(roleName).Do(ctx)
		require.Falsef(t, slices.Contains(assignedUsers, adminUser), "should not have %q role", roleName)
	})

	t.Run("assign and revoke a role", func(t *testing.T) {
		roleName := "AssignRevokeMe"

		mustCreateRole(t, rbac.NewRole(roleName,
			rbac.BackupsPermission{Actions: []string{models.PermissionActionManageBackups}, Collection: pizza},
		))

		// Act: assign
		err := usersClient.DB().RolesAssigner().WithUserID(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "assign %q role", roleName)

		assignedUsers, _ := rolesClient.UserAssignmentGetter().WithRole(roleName).Do(ctx)
		require.Truef(t, slices.ContainsFunc(assignedUsers,
			func(e rbac.UserAssignment) bool { return e.UserID == adminUser }), "should have %q role", roleName)

		// Act: revoke
		err = usersClient.DB().RolesRevoker().WithUserID(adminUser).WithRoles(roleName).Do(ctx)
		require.NoErrorf(t, err, "revoke %q role", roleName)

		assignedUsers, _ = rolesClient.UserAssignmentGetter().WithRole(roleName).Do(ctx)
		require.Falsef(t, slices.ContainsFunc(assignedUsers,
			func(e rbac.UserAssignment) bool { return e.UserID == adminUser }), "should not have %q role", roleName)
	})

	t.Run("create user, deactivate, check, reactivate, check", func(t *testing.T) {
		userId := "test-user-" + uuid.New().String()

		apikey, err := usersClient.DB().Creator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.Greater(t, len(apikey), 10, "apikey is too short")

		defer func() {
			_, errd := usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
			require.NoError(t, errd)
		}()

		success, err := usersClient.DB().Deactivator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, success)

		userInfo, err := usersClient.DB().Getter().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.False(t, userInfo.Active)

		success, err = usersClient.DB().Activator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, success)

		userInfo, err = usersClient.DB().Getter().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, userInfo.Active)
	})

	t.Run("create user, deactivate succeeds, deactivate conflicts, activate succeeds, activate conflicts", func(t *testing.T) {
		userId := "test-user-" + uuid.New().String()

		apikey, err := usersClient.DB().Creator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.Greater(t, len(apikey), 10, "apikey is too short")

		defer func() {
			_, errd := usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
			require.NoError(t, errd)
		}()

		success, err := usersClient.DB().Deactivator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, success)

		success, err = usersClient.DB().Deactivator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.False(t, success)

		success, err = usersClient.DB().Activator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, success)

		success, err = usersClient.DB().Activator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.False(t, success)
	})

	t.Run("create user, get user info, delete user 200, delete user return 400", func(t *testing.T) {
		userId := "test-user-" + uuid.New().String()
		_, err := usersClient.DB().Creator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)

		defer func() {
			_, errd := usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
			require.NoError(t, errd)
		}()

		userInfo, err := usersClient.DB().Getter().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, userInfo.Active)

		deleted, err := usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, deleted)

		deleted, err = usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.False(t, deleted)
	})

	t.Run("create users, list all, check new users are there", func(t *testing.T) {
		userIds := map[string]struct{}{
			"test-user-" + uuid.New().String(): {},
			"test-user-" + uuid.New().String(): {},
			"test-user-" + uuid.New().String(): {},
		}
		for userId := range userIds {
			_, err := usersClient.DB().Creator().WithUserID(userId).Do(ctx)
			require.NoError(t, err)

			defer func() {
				_, errd := usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
				require.NoError(t, errd)
			}()
		}
		usersInfo, err := usersClient.DB().Lister().Do(ctx)
		require.NoError(t, err)
		found := 0
		for _, userInfo := range usersInfo {
			if _, ok := userIds[userInfo.UserID]; ok {
				found++
				continue
			}
		}
		require.Equal(t, len(userIds), found)
	})

	t.Run("create user, get key, rotate key, check new key is different", func(t *testing.T) {
		userId := "test-user-" + uuid.New().String()
		apikey, err := usersClient.DB().Creator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, apikey)
		defer func() {
			_, errd := usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
			require.NoError(t, errd)
		}()

		apikeyRotated, err := usersClient.DB().KeyRotator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.NotEqual(t, apikey, apikeyRotated)
	})

	t.Run("create user, get key, deactivate user, activate user, login user fails, rotate key, login user succeeds", func(t *testing.T) {
		userId := "test-user-" + uuid.New().String()
		apikey, err := usersClient.DB().Creator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, apikey)
		defer func() {
			_, errd := usersClient.DB().Deleter().WithUserID(userId).Do(ctx)
			require.NoError(t, errd)
		}()

		success, err := usersClient.DB().Deactivator().WithRevokeKey(true).WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, success)
		success, err = usersClient.DB().Activator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)
		require.True(t, success)

		cfg := weaviate.Config{
			Host:   container.HTTPAddress(),
			Scheme: "http",
			AuthConfig: auth.ApiKey{
				Value: apikey,
			},
		}

		client, err := weaviate.NewClient(cfg)
		require.NoError(t, err, "create test client")

		_, err = client.Users().MyUserGetter().Do(ctx)
		require.Error(t, err)

		apikeyRotated, err := usersClient.DB().KeyRotator().WithUserID(userId).Do(ctx)
		require.NoError(t, err)

		cfg = weaviate.Config{
			Host:   container.HTTPAddress(),
			Scheme: "http",
			AuthConfig: auth.ApiKey{
				Value: apikeyRotated,
			},
		}

		client, err = weaviate.NewClient(cfg)
		require.NoError(t, err, "create second test client")

		myUser, err := client.Users().MyUserGetter().Do(ctx)
		require.NoError(t, err)
		require.Equal(t, myUser.UserID, userId)
	})
}
