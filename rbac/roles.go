package rbac

import (
	"context"
	"fmt"

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
	Aliases     []AliasPermission
	Backups     []BackupsPermission
	Cluster     []ClusterPermission
	Collections []CollectionPermission
	Data        []DataPermission
	Groups      []GroupPermission
	Namespaces  []NamespacePermission
	Nodes       []NodesPermission
	Replication []ReplicationPermission
	Roles       []RolePermission
	Tenants     []TenantPermission
	Users       []UserPermission
}

type (
	AliasPermission      api.AliasPermission
	BackupsPermission    api.BackupsPermission
	ClusterPermission    api.ClusterPermission
	CollectionPermission api.CollectionPermission
	DataPermission       api.DataPermission
	GroupPermission      api.GroupPermission
	NodeVerbosity        api.NodeVerbosity
	NodesPermission      struct {
		Collection string
		Verbosity  NodeVerbosity

		Read bool
	}
	NamespacePermission   api.NamespacePermission
	ReplicationPermission api.ReplicationPermission
	RoleScope             api.RoleScope
	RolePermission        struct {
		RoleID string
		Scope  RoleScope

		Create bool
		Read   bool
		Update bool
		Delete bool
	}
	TenantPermission api.TenantPermission
	UserPermission   api.UserPermission
)

const (
	NodeVerbosityMinimal = NodeVerbosity(api.NodeVerbosityMinimal)
	NodeVerbosityVerbose = NodeVerbosity(api.NodeVerbosityVerbose)

	RoleScopeAll   = RoleScope(api.RoleScopeAll)
	RoleScopeMatch = RoleScope(api.RoleScopeMatch)
)

func (rc *RolesClient) Create(ctx context.Context, r Role) error {
	role := api.Role{ID: r.ID}
	for _, p := range r.Aliases {
		role.Aliases = append(role.Aliases, api.AliasPermission(p))
	}
	for _, p := range r.Backups {
		role.Backups = append(role.Backups, api.BackupsPermission(p))
	}
	for _, p := range r.Cluster {
		role.Cluster = append(role.Cluster, api.ClusterPermission(p))
	}
	for _, p := range r.Collections {
		role.Collections = append(role.Collections, api.CollectionPermission(p))
	}
	for _, p := range r.Data {
		role.Data = append(role.Data, api.DataPermission(p))
	}
	for _, p := range r.Groups {
		role.Groups = append(role.Groups, api.GroupPermission(p))
	}
	for _, p := range r.Namespaces {
		role.Namespaces = append(role.Namespaces, api.NamespacePermission(p))
	}
	for _, p := range r.Nodes {
		role.Nodes = append(role.Nodes, api.NodesPermission{
			Collection: p.Collection,
			Verbosity:  api.NodeVerbosity(p.Verbosity),
		})
	}
	for _, p := range r.Replication {
		role.Replication = append(role.Replication, api.ReplicationPermission(p))
	}
	for _, p := range r.Roles {
		role.Roles = append(role.Roles, api.RolePermission{
			RoleID: p.RoleID,
			Scope:  api.RoleScope(p.Scope),
		})
	}
	for _, p := range r.Tenants {
		role.Tenants = append(role.Tenants, api.TenantPermission(p))
	}
	for _, p := range r.Users {
		role.Users = append(role.Users, api.UserPermission(p))
	}

	req := &api.CreateRoleRequest{Role: role}
	if err := rc.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

func (rc *RolesClient) Exists(ctx context.Context, roleID string) (bool, error) {
	var resp api.ResourceExistsResponse
	if err := rc.transport.Do(ctx, api.GetRoleRequest(roleID), &resp); err != nil {
		return false, fmt.Errorf("check role exists: %w", err)
	}
	return resp.Bool(), nil
}

func (rc *RolesClient) Get(ctx context.Context, roleID string) (*Role, error) {
	var resp api.Role
	if err := rc.transport.Do(ctx, api.GetRoleRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}
	return unmarshalRole(&resp), nil
}

func (rc *RolesClient) List(ctx context.Context) ([]Role, error) {
	var resp []api.Role
	if err := rc.transport.Do(ctx, api.ListRolesRequest, &resp); err != nil {
		return nil, fmt.Errorf("list role: %w", err)
	}
	var roles []Role
	for i := range resp {
		roles = append(roles, *unmarshalRole(&resp[i]))
	}
	return roles, nil
}

func (rc *RolesClient) Delete(ctx context.Context, roleID string) error {
	if err := rc.transport.Do(ctx, api.DeleteRoleRequest(roleID), nil); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	return nil
}

func (rc *RolesClient) AddPermissions(context.Context, Permissions) error    { return nil }
func (rc *RolesClient) RemovePermissions(context.Context, Permissions) error { return nil }

type Permission struct {
	Alias       AliasPermission
	Backups     BackupsPermission
	Cluster     ClusterPermission
	Collections CollectionPermission
	Data        DataPermission
	Groups      GroupPermission
	Nodes       NodesPermission
	Replication ReplicationPermission
	Roles       RolePermission
	Tenants     TenantPermission
	Users       UserPermission
}

func (rc *RolesClient) HasPermission(context.Context) (bool, error) { return false, nil }

func (rc *RolesClient) AssignedUserIDs(context.Context) ([]string, error) { return nil, nil }

func (rc *RolesClient) UserAssignments(context.Context) ([]map[string]string, error) { return nil, nil }

func (rc *RolesClient) GroupAssignments(context.Context) ([]map[string]string, error) {
	return nil, nil
}

func unmarshalRole(r *api.Role) *Role {
	role := Role{ID: r.ID}
	for _, p := range r.Aliases {
		role.Aliases = append(role.Aliases, AliasPermission(p))
	}
	for _, p := range r.Backups {
		role.Backups = append(role.Backups, BackupsPermission(p))
	}
	for _, p := range r.Cluster {
		role.Cluster = append(role.Cluster, ClusterPermission(p))
	}
	for _, p := range r.Collections {
		role.Collections = append(role.Collections, CollectionPermission(p))
	}
	for _, p := range r.Data {
		role.Data = append(role.Data, DataPermission(p))
	}
	for _, p := range r.Groups {
		role.Groups = append(role.Groups, GroupPermission(p))
	}
	for _, p := range r.Namespaces {
		role.Namespaces = append(role.Namespaces, NamespacePermission(p))
	}
	for _, p := range r.Nodes {
		role.Nodes = append(role.Nodes, NodesPermission{
			Collection: p.Collection,
			Verbosity:  NodeVerbosity(p.Verbosity),
		})
	}
	for _, p := range r.Replication {
		role.Replication = append(role.Replication, ReplicationPermission(p))
	}
	for _, p := range r.Roles {
		role.Roles = append(role.Roles, RolePermission{
			RoleID: p.RoleID,
			Scope:  RoleScope(p.Scope),
		})
	}
	for _, p := range r.Tenants {
		role.Tenants = append(role.Tenants, TenantPermission(p))
	}
	for _, p := range r.Users {
		role.Users = append(role.Users, UserPermission(p))
	}
	return &role
}
