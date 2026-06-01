package rbac

import (
	"context"
	"fmt"
	"slices"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type RolesClient struct {
	transport internal.Transport
}

func NewRolesClient(t internal.Transport) *RolesClient {
	dev.AssertNotNil(t, "transport")
	return &RolesClient{transport: t}
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
	MCP         []MCPPermission
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
	MCPPermission         api.MCPPermission
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

func (c *RolesClient) Create(ctx context.Context, r Role) error {
	req := &api.CreateRoleRequest{
		Role: api.Role{
			ID:          r.ID,
			Permissions: marshalPermissions(&r.Permissions),
		},
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

// Exists returns true and nil error if role with the given roleID exists.
func (c *RolesClient) Exists(ctx context.Context, roleID string) (bool, error) {
	var resp api.ResourceExistsResponse
	if err := c.transport.Do(ctx, api.GetRoleRequest(roleID), &resp); err != nil {
		return false, fmt.Errorf("check role exists: %w", err)
	}
	return resp.Bool(), nil
}

// Get fetches a role with the given roleID.
func (c *RolesClient) Get(ctx context.Context, roleID string) (*Role, error) {
	var resp api.Role
	if err := c.transport.Do(ctx, api.GetRoleRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}
	role := unmarshalRole(&resp)
	return &role, nil
}

// List fetches all roles defined in the cluster.
func (c *RolesClient) List(ctx context.Context) ([]Role, error) {
	var resp []api.Role
	if err := c.transport.Do(ctx, api.ListRolesRequest, &resp); err != nil {
		return nil, fmt.Errorf("list role: %w", err)
	}
	roles := slices.Grow([]Role(nil), len(resp))
	for i := range resp {
		roles = append(roles, unmarshalRole(&resp[i]))
	}
	return roles, nil
}

// Delete deletes a role with the given roleID.
func (c *RolesClient) Delete(ctx context.Context, roleID string) error {
	if err := c.transport.Do(ctx, api.DeleteRoleRequest(roleID), nil); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}
	return nil
}

type AddPermissions struct {
	RoleID      string
	Permissions Permissions
}

func (c *RolesClient) AddPermissions(ctx context.Context, options AddPermissions) error {
	req := &api.ManagePermissionsRequest{
		RoleID:      options.RoleID,
		Verb:        api.PermissionVerbAdd,
		Permissions: marshalPermissions(&options.Permissions),
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("add role permissions: %w", err)
	}
	return nil
}

type RemovePermissions struct {
	RoleID      string
	Permissions Permissions
}

func (c *RolesClient) RemovePermissions(ctx context.Context, options RemovePermissions) error {
	req := &api.ManagePermissionsRequest{
		RoleID:      options.RoleID,
		Verb:        api.PermissionVerbRemove,
		Permissions: marshalPermissions(&options.Permissions),
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("remove role permissions: %w", err)
	}
	return nil
}

type HasPermission struct {
	RoleID string

	Alias       AliasPermission
	Backups     BackupsPermission
	Cluster     ClusterPermission
	Collection  CollectionPermission
	Data        DataPermission
	Group       GroupPermission
	MCP         MCPPermission
	Namespaces  NamespacePermission
	Nodes       NodesPermission
	Replication ReplicationPermission
	Role        RolePermission
	Tenant      TenantPermission
	User        UserPermission
}

// HasPermission checks if a role contains a permission.
// Only one permission can be checked in a single call.
func (c *RolesClient) HasPermission(ctx context.Context, options HasPermission) (bool, error) {
	req, err := marshalHasPermission(options)
	if err != nil {
		return false, err
	}
	req.RoleID = options.RoleID

	var resp api.HasPermissionResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return false, fmt.Errorf("check role has permission: %w", err)
	}
	return bool(resp), nil
}

func (c *RolesClient) AssignedUserIDs(ctx context.Context, roleID string) ([]string, error) {
	var resp api.GetAssignedUsersResponse
	if err := c.transport.Do(ctx, api.GetAssignedUsersRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get assigned users: %w", err)
	}
	return []string(resp), nil
}

type UserInfo api.UserInfo

func (c *RolesClient) UserAssignments(ctx context.Context, roleID string) ([]UserInfo, error) {
	var resp []api.UserInfo
	if err := c.transport.Do(ctx, api.GetUserAssignmentsRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get user assignments: %w", err)
	}

	out := make([]UserInfo, len(resp))
	for i := range resp {
		out[i] = UserInfo(resp[i])
	}
	return out, nil
}

type GroupInfo api.GroupInfo

func (c *RolesClient) GroupAssignments(ctx context.Context, roleID string) ([]GroupInfo, error) {
	var resp []api.GroupInfo
	if err := c.transport.Do(ctx, api.GetGroupAssignmentsRequest(roleID), &resp); err != nil {
		return nil, fmt.Errorf("get group assignments: %w", err)
	}

	out := make([]GroupInfo, len(resp))
	for i := range resp {
		out[i] = GroupInfo(resp[i])
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
	for _, p := range ps.MCP {
		out.MCP = append(out.MCP, api.MCPPermission(p))
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

func unmarshalRole(r *api.Role) Role {
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
	for _, p := range r.MCP {
		role.MCP = append(role.MCP, MCPPermission(p))
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
	return role
}

// marshalHasPermission returns [api.HasPermissionRequest] with exacly 1 permission
// (action + data) pair set and returns an error if more permissions are set in has.
// This validation happens at this step, because internal/api expects the request
// to be valid and will only send a single permission, even if more have been set,
// so the server cannot perform this validation.
//
// NOTE(dyma): we might want to introduce some Validator interface in the transport
// layer, which every request object can optionally implement, so that all validation
// is done right before we marshal the request and dispatch to the actual transport.
func marshalHasPermission(has HasPermission) (*api.HasPermissionRequest, error) {
	var count int
	for _, check := range []bool{
		has.Alias.Create,
	} {
		if check {
			count++
		}
	}
	if has.Alias.Create {
		count++
	}
	if has.Alias.Read {
		count++
	}
	if has.Alias.Update {
		count++
	}
	if has.Alias.Delete {
		count++
	}
	if has.Backups.Manage {
		count++
	}
	if has.Cluster.Read {
		count++
	}
	if has.Collection.Create {
		count++
	}
	if has.Collection.Read {
		count++
	}
	if has.Collection.Update {
		count++
	}
	if has.Collection.Delete {
		count++
	}
	if has.Data.Create {
		count++
	}
	if has.Data.Read {
		count++
	}
	if has.Data.Update {
		count++
	}
	if has.Data.Delete {
		count++
	}
	if has.Group.Read {
		count++
	}
	if has.Group.AssignAndRevoke {
		count++
	}
	if has.Namespaces.Manage {
		count++
	}
	if has.Nodes.Read {
		count++
	}
	if has.MCP.Create {
		count++
	}
	if has.MCP.Read {
		count++
	}
	if has.MCP.Update {
		count++
	}
	if has.Replication.Create {
		count++
	}
	if has.Replication.Read {
		count++
	}
	if has.Replication.Update {
		count++
	}
	if has.Replication.Delete {
		count++
	}
	if has.Role.Create {
		count++
	}
	if has.Role.Read {
		count++
	}
	if has.Role.Update {
		count++
	}
	if has.Role.Delete {
		count++
	}
	if has.Tenant.Create {
		count++
	}
	if has.Tenant.Read {
		count++
	}
	if has.Tenant.Update {
		count++
	}
	if has.Tenant.Delete {
		count++
	}
	if has.User.Create {
		count++
	}
	if has.User.Read {
		count++
	}
	if has.User.Update {
		count++
	}
	if has.User.Delete {
		count++
	}
	if has.User.AssignAndRevoke {
		count++
	}

	if count == 1 {
		return &api.HasPermissionRequest{
			Alias:       api.AliasPermission(has.Alias),
			Backups:     api.BackupsPermission(has.Backups),
			Cluster:     api.ClusterPermission(has.Cluster),
			Collections: api.CollectionPermission(has.Collection),
			Data:        api.DataPermission(has.Data),
			Groups:      api.GroupPermission(has.Group),
			Nodes: api.NodesPermission{
				Collection: has.Nodes.Collection,
				Verbosity:  api.NodeVerbosity(has.Nodes.Verbosity),
			},
			Replication: api.ReplicationPermission(has.Replication),
			Roles: api.RolePermission{
				RoleID: has.Role.RoleID,
				Scope:  api.RoleScope(has.Role.Scope),
			},
			Tenants: api.TenantPermission(has.Tenant),
			Users:   api.UserPermission(has.User),
		}, nil
	}

	return nil, fmt.Errorf("has-permission accepts 1 permission but %d were provided", count)
}
