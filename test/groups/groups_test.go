package groups

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

func TestGroup_integration(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	rolesClient := client.Roles()
	groupsClient := client.Groups()

	const (
		group = "/with-special-characer"
		pizza = "Pizza"
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

	roleName := "TestingGroups"

	mustCreateRole(t, rbac.NewRole(roleName,
		rbac.BackupsPermission{Actions: []string{models.PermissionActionManageBackups}, Collection: pizza},
	))

	groups, err := groupsClient.OIDC().RolesGetter().WithGroupID(group).Do(ctx)
	require.NoError(t, err)
	require.Len(t, groups, 0)

	require.NoErrorf(t, groupsClient.OIDC().RolesAssigner().WithGroupId(group).WithRoles(roleName).Do(ctx), "assign %q role", roleName)

	groups, err = groupsClient.OIDC().RolesGetter().WithGroupID(group).Do(ctx)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, roleName, groups[0].Name)

	knownGroups, err := groupsClient.OIDC().GetKnownGroups().Do(ctx)
	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, roleName, knownGroups[0])

	require.NoErrorf(t, groupsClient.OIDC().RolesRevoker().WithGroupId(group).WithRoles(roleName).Do(ctx), "assign %q role", roleName)

	groups, err = groupsClient.OIDC().RolesGetter().WithGroupID(group).Do(ctx)
	require.NoError(t, err)
	require.Len(t, groups, 0)
}
