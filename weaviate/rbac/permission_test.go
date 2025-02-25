package rbac_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
)

func TestRole_UnmarshalJSON(t *testing.T) {
	var got *rbac.Role

	data := []byte(`{
	"name": "WeaviateRole",
	"permissions": [
		{"action": "manage_backups", "backups": {"collection": "Pizza"}},
		{"action": "manage_backups", "backups": {"collection": "Songs"}},
		{"action": "read_cluster"},
		{"action": "create_collections", "collections": {"collection": "Pizza"}},
		{"action": "read_collections", "collections": {"collection": "Pizza"}},
		{"action": "update_collections", "collections": {"collection": "Songs"}},
		{"action": "delete_collections", "collections": {"collection": "Songs"}},
		{"action": "create_data", "data": {"collection": "Pizza"}},
		{"action": "read_data", "data": {"collection": "Pizza"}},
		{"action": "update_data", "data": {"collection": "Songs"}},
		{"action": "delete_data", "data": {"collection": "Songs"}},
		{"action": "read_nodes", "nodes": {"collection": "Pizza", "verbosity": "minimal"}},
		{"action": "read_nodes", "nodes": {"collection": "Pizza", "verbosity": "verbose"}},
		{"action": "read_nodes", "nodes": {"collection": "Songs", "verbosity": "minimal"}},
		{"action": "create_roles", "roles": {"role": "CreatorReader", "scope": "all"}},
		{"action": "read_roles", "roles": {"role": "CreatorReader", "scope": "all"}},
		{"action": "update_roles", "roles": {"role": "UpdaterDeleter", "scope": "matching"}},
		{"action": "delete_roles", "roles": {"role": "UpdaterDeleter", "scope": "matching"}},
		{"action": "create_tenants"}, {"action": "read_tenants"},
		{"action": "read_users"}, {"action": "assign_and_revoke_users"}
	]
}`)

	want := rbac.NewRole("WeaviateRole",
		rbac.BackupsPermission{
			Actions:    []string{"manage_backups"},
			Collection: "Pizza",
		},
		rbac.BackupsPermission{
			Actions:    []string{"manage_backups"},
			Collection: "Songs",
		},
		rbac.ClusterPermission{Actions: []string{"read_cluster"}},
		rbac.CollectionsPermission{
			Actions:    []string{"create_collections", "read_collections"},
			Collection: "Pizza",
		},
		rbac.CollectionsPermission{
			Actions:    []string{"update_collections", "delete_collections"},
			Collection: "Songs",
		},
		rbac.DataPermission{
			Actions:    []string{"create_data", "read_data"},
			Collection: "Pizza",
		},
		rbac.DataPermission{
			Actions:    []string{"update_data", "delete_data"},
			Collection: "Songs",
		},
		rbac.NodesPermission{
			Actions:    []string{"read_nodes"},
			Collection: "Pizza",
			Verbosity:  "minimal",
		},
		rbac.NodesPermission{
			Actions:    []string{"read_nodes"},
			Collection: "Pizza",
			Verbosity:  "verbose",
		},
		rbac.NodesPermission{
			Actions:    []string{"read_nodes"},
			Collection: "Songs",
			Verbosity:  "minimal",
		},
		rbac.RolesPermission{
			Actions: []string{"create_roles", "read_roles"},
			Role:    "CreatorReader",
			Scope:   "all",
		},
		rbac.RolesPermission{
			Actions: []string{"update_roles", "delete_roles"},
			Role:    "UpdaterDeleter",
			Scope:   "matching",
		},
		rbac.TenantsPermission{Actions: []string{"create_tenants", "read_tenants"}},
		rbac.UsersPermission{Actions: []string{"read_users", "assign_and_revoke_users"}},
	)

	err := json.Unmarshal(data, &got)
	require.NoError(t, err, "unmarshal Weaviate role")

	require.Equal(t, want.Name, got.Name, "role name")
	require.ElementsMatch(t, want.Backups, got.Backups)
	require.ElementsMatch(t, want.Cluster, got.Cluster)
	require.ElementsMatch(t, want.Collections, got.Collections)
	require.ElementsMatch(t, want.Data, got.Data)
	require.ElementsMatch(t, want.Nodes, got.Nodes)
	require.ElementsMatch(t, want.Roles, got.Roles)
	require.ElementsMatch(t, want.Tenants, got.Tenants)
	require.ElementsMatch(t, want.Users, got.Users)
}
