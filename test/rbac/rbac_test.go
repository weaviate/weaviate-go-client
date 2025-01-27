package rbac

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
)

func TestRBAC_integration(t *testing.T) {
	if err := testenv.SetupLocalWeaviate(); err != nil {
		t.Fatalf("failed to setup weaviate: %s", err)
	}
	defer func() {
		fmt.Printf("TestBackups_integration TEAR DOWN START\n")
		if err := testenv.TearDownLocalWeaviate(); err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
		fmt.Printf("TestBackups_integration TEAR DOWN STOP\n")
	}()
	client := testsuit.CreateTestClient(false)
	testsuit.CleanUpWeaviate(t, client)

	ctx := context.Background()
	rolesClient := client.Roles()

	const (
		adminRole  = "admin"
		viewerRole = "viewer"

		adminUser = "ms_2d0e007e7136de11d5f29fce7a53dae219a51458@existiert.net"
	)

	t.Run("get all roles", func(t *testing.T) {
		all, err := rolesClient.AllGetter().Do(ctx)
		require.NoError(t, err)
		require.Len(t, all, 2, "wrong number of roles")
		require.Equal(t, all[0].Name, adminRole)
		require.Equal(t, all[1].Name, viewerRole)
	})
	t.Run("get user roles", func(t *testing.T) {

	})
	t.Run("get assigned users", func(t *testing.T) {

	})
	t.Run("get all roles", func(t *testing.T) {

	})
	t.Run("create role", func(t *testing.T) {

	})
	t.Run("add permissions", func(t *testing.T) {

	})
	t.Run("remove permissions", func(t *testing.T) {

	})
	t.Run("assign and revoke a role", func(t *testing.T) {

	})
}
