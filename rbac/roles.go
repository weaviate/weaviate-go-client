package rbac

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type RolesClient struct {
	transport internal.Transport
}

type Role struct {
	ID string
	Permissions
}

type Permissions struct {
	Alias       []AliasPermission
	Backups     []BackupsPermission
	Cluster     []ClusterPermission
	Collections []CollectionsPermission
	Data        []DataPermission
	Groups      []GroupsPermission
	Nodes       []NodesPermission
	Replication []ReplicationPermission
	Roles       []RolesPermission
	Tenants     []TenantsPermission
	Users       []UsersPermission
}

type (
	AliasPermission       api.AliasPermission
	BackupsPermission     api.BackupsPermission
	ClusterPermission     api.ClusterPermission
	CollectionsPermission api.CollectionsPermission
	DataPermission        api.DataPermission
	GroupsPermission      api.GroupsPermission
	NodesVerbosity        api.NodesVerbosity
	NodesPermission       api.NodesPermission
	ReplicationPermission api.ReplicationPermission
	RolesScope            api.RolesScope
	RolesPermission       api.RolesPermission
	TenantsPermission     api.TenantsPermission
	UsersPermission       api.UsersPermission
)

func (rc *RolesClient) Create(context.Context, Role) error                      { return nil }
func (rc *RolesClient) Exists(ctx context.Context, roleID string) (bool, error) { return false, nil }
func (rc *RolesClient) Get(ctx context.Context, roleID string) (*Role, error)   { return nil, nil }
func (rc *RolesClient) List(ctx context.Context) ([]Role, error)                { return nil, nil }
func (rc *RolesClient) Delete(ctx context.Context, roleID string) error         { return nil }

func (rc *RolesClient) AddPermissions(context.Context, Permissions) error    { return nil }
func (rc *RolesClient) RemovePermissions(context.Context, Permissions) error { return nil }

type Permission struct {
	Alias       AliasPermission
	Backups     BackupsPermission
	Cluster     ClusterPermission
	Collections CollectionsPermission
	Data        DataPermission
	Groups      GroupsPermission
	Nodes       NodesPermission
	Replication ReplicationPermission
	Roles       RolesPermission
	Tenants     TenantsPermission
	Users       UsersPermission
}

func (rc *RolesClient) HasPermission(context.Context) (bool, error) { return false, nil }

func (rc *RolesClient) AssignedUserIDs(context.Context) ([]string, error) { return nil, nil }

func (rc *RolesClient) UserAssignments(context.Context) ([]map[string]string, error) { return nil, nil }

func (rc *RolesClient) GroupAssignments(context.Context) ([]map[string]string, error) {
	return nil, nil
}
