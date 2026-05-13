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
	req := &api.CreateRoleRequest{
		Role: api.Role{
			ID:          r.ID,
			Permissions: marshalPermissions(&r.Permissions),
		},
	}
	if err := rc.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

// Exists returns true and nil error if role with the given roleID exists.
func (rc *RolesClient) Exists(ctx context.Context, roleID string) (bool, error) {
	var resp api.ResourceExistsResponse
	if err := rc.transport.Do(ctx, api.GetRoleRequest(roleID), &resp); err != nil {
		return false, fmt.Errorf("check role exists: %w", err)
	}
	return resp.Bool(), nil
}

// Get fetches a role with the given roleID.
func (rc *RolesClient) Get(ctx context.Context, roleID string) (*Role, error) {
	var resp api.Role
	if err := rc.transport.Do(ctx, api.GetRoleRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}
	return unmarshalRole(&resp), nil
}

// List fetches all roles defined in the cluster.
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

// Delete deletes a role with the given roleID.
func (rc *RolesClient) Delete(ctx context.Context, roleID string) error {
	if err := rc.transport.Do(ctx, api.DeleteRoleRequest(roleID), nil); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	return nil
}

type AddPermissions Role

func (rc *RolesClient) AddPermissions(ctx context.Context, options AddPermissions) error {
	req := &api.AddPermissionsRequest{
		RoleID:      options.ID,
		Permissions: marshalPermissions(&options.Permissions),
	}
	if err := rc.transport.Do(ctx, &req, nil); err != nil {
		return fmt.Errorf("add role permissions: %w", err)
	}
	return nil
}

type RemovePermissions Role

func (rc *RolesClient) RemovePermissions(ctx context.Context, options RemovePermissions) error {
	req := &api.RemovePermissionsRequest{
		RoleID:      options.ID,
		Permissions: marshalPermissions(&options.Permissions),
	}
	if err := rc.transport.Do(ctx, &req, nil); err != nil {
		return fmt.Errorf("remove role permissions: %w", err)
	}
	return nil
}

type HasPermission struct {
	RoleID string

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

// HasPermission checks if a role contains a permission.
// Only one permission can be checked in a single call.
func (rc *RolesClient) HasPermission(ctx context.Context, roleID string, options HasPermission) (bool, error) {
	req := &api.HasPermissionRequest{
		RoleID:      options.RoleID,
		Alias:       api.AliasPermission(options.Alias),
		Backups:     api.BackupsPermission(options.Backups),
		Cluster:     api.ClusterPermission(options.Cluster),
		Collections: api.CollectionPermission(options.Collections),
		Data:        api.DataPermission(options.Data),
		Groups:      api.GroupPermission(options.Groups),
		Nodes: api.NodesPermission{
			Collection: options.Nodes.Collection,
			Verbosity:  api.NodeVerbosity(options.Nodes.Verbosity),
		},
		Replication: api.ReplicationPermission(options.Replication),
		Roles: api.RolePermission{
			RoleID: options.Roles.RoleID,
			Scope:  api.RoleScope(options.Roles.Scope),
		},
		Tenants: api.TenantPermission(options.Tenants),
		Users:   api.UserPermission(options.Users),
	}
	var resp api.HasPermissionResponse
	if err := rc.transport.Do(ctx, req, &resp); err != nil {
		return false, fmt.Errorf("check role has permission: %w", err)
	}
	return bool(resp), nil
}

func (rc *RolesClient) AssignedUserIDs(ctx context.Context, roleID string) ([]string, error) {
	var resp api.GetAssignedUsersResponse
	if err := rc.transport.Do(ctx, api.GetAssignedUsersRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get assigned users: %w", err)
	}
	return []string(resp), nil
}

type UserAssignment api.UserAssignment

func (rc *RolesClient) UserAssignments(ctx context.Context, roleID string) ([]UserAssignment, error) {
	var resp []api.UserAssignment
	if err := rc.transport.Do(ctx, api.GetUserAssignmentsRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get user assignments: %w", err)
	}

	out := make([]UserAssignment, len(resp))
	for i := range resp {
		out[i] = UserAssignment(resp[i])
	}
	return out, nil
}

type GroupAssignment api.GroupAssignment

func (rc *RolesClient) GroupAssignments(ctx context.Context, roleID string) ([]GroupAssignment, error) {
	var resp []api.GroupAssignment
	if err := rc.transport.Do(ctx, api.GetGroupAssignmentsRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get group assignments: %w", err)
	}

	out := make([]GroupAssignment, len(resp))
	for i := range resp {
		out[i] = GroupAssignment(resp[i])
	}
	return out, nil
}

func marshalPermissions(ps *Permissions) api.Permissions {
	var out api.Permissions
	for _, p := range ps.Aliases {
		out.Aliases = append(out.Aliases, api.AliasPermission(p))
	}
	for _, p := range ps.Backups {
		out.Backups = append(out.Backups, api.BackupsPermission(p))
	}
	for _, p := range ps.Cluster {
		out.Cluster = append(out.Cluster, api.ClusterPermission(p))
	}
	for _, p := range ps.Collections {
		out.Collections = append(out.Collections, api.CollectionPermission(p))
	}
	for _, p := range ps.Data {
		out.Data = append(out.Data, api.DataPermission(p))
	}
	for _, p := range ps.Groups {
		out.Groups = append(out.Groups, api.GroupPermission(p))
	}
	for _, p := range ps.Namespaces {
		out.Namespaces = append(out.Namespaces, api.NamespacePermission(p))
	}
	for _, p := range ps.Nodes {
		out.Nodes = append(out.Nodes, api.NodesPermission{
			Collection: p.Collection,
			Verbosity:  api.NodeVerbosity(p.Verbosity),
		})
	}
	for _, p := range ps.Replication {
		out.Replication = append(out.Replication, api.ReplicationPermission(p))
	}
	for _, p := range ps.Roles {
		out.Roles = append(out.Roles, api.RolePermission{
			RoleID: p.RoleID,
			Scope:  api.RoleScope(p.Scope),
		})
	}
	for _, p := range ps.Tenants {
		out.Tenants = append(out.Tenants, api.TenantPermission(p))
	}
	for _, p := range ps.Users {
		out.Users = append(out.Users, api.UserPermission(p))
	}
	return out
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
