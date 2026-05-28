package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type Role struct {
	ID          string `json:"name"`
	Permissions `json:"permissions"`
}

var (
	_ json.Marshaler   = (*Role)(nil)
	_ json.Unmarshaler = (*Role)(nil)
)

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

var _ json.Marshaler = (*Permissions)(nil)

type (
	AliasPermission struct {
		Collection string `json:"collection,omitempty"`
		Alias      string `json:"alias,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	BackupsPermission struct {
		Collection string `json:"collection,omitempty"`

		Manage bool `json:"-"`
	}
	ClusterPermission struct {
		Read bool `json:"-"`
	}
	CollectionPermission struct {
		Collection string `json:"collection,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	DataPermission struct {
		Collection string `json:"collection,omitempty"`
		Tenant     string `json:"tenant,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	GroupPermission struct {
		GroupID string `json:"group,omitempty"`
		Type    string `json:"groupType,omitempty"`

		Read            bool `json:"-"`
		AssignAndRevoke bool `json:"-"`
	}
	NamespacePermission struct {
		Namespace string `json:"namespace,omitempty"`

		Manage bool `json:"-"`
	}
	NodeVerbosity   rest.PermissionNodesVerbosity
	NodesPermission struct {
		Collection string        `json:"collection,omitempty"`
		Verbosity  NodeVerbosity `json:"verbosity,omitempty"`

		Read bool `json:"-"`
	}
	MCPPermission struct {
		Create bool
		Read   bool
		Update bool
	}
	ReplicationPermission struct {
		Collection string `json:"collection,omitempty"`
		Shard      string `json:"shard,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	RoleScope      rest.PermissionRolesScope
	RolePermission struct {
		RoleID string    `json:"role,omitempty"`
		Scope  RoleScope `json:"scope,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	TenantPermission struct {
		Collection string `json:"collection,omitempty"`
		Tenant     string `json:"tenant,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	UserPermission struct {
		UserID string `json:"user,omitempty"`

		Create          bool `json:"-"`
		Read            bool `json:"-"`
		Update          bool `json:"-"`
		Delete          bool `json:"-"`
		AssignAndRevoke bool `json:"-"`
	}
)

const (
	NodeVerbosityMinimal = NodeVerbosity(rest.Minimal)
	NodeVerbosityVerbose = NodeVerbosity(rest.Verbose)

	RoleScopeAll   = RoleScope(rest.PermissionRolesScopeAll)
	RoleScopeMatch = RoleScope(rest.PermissionRolesScopeMatch)
)

type CreateRoleRequest struct {
	transports.BaseEndpoint
	Role Role
}

var _ transports.Endpoint = (*CreateRoleRequest)(nil)

func (r *CreateRoleRequest) Method() string { return http.MethodPost }
func (r *CreateRoleRequest) Path() string   { return "/authz/roles" }
func (r *CreateRoleRequest) Body() any      { return &r.Role }

// GetRoleRequest fetches the role by it's ID.
// Use with [ResourceExistsResponse] or [Role] response types.
var GetRoleRequest = transports.IdentityEndpoint[string](http.MethodGet, "/authz/roles/%s")

// ListRolesRequest fetches all roles defined in the cluster.
// Use with []Role response type.
var ListRolesRequest = transports.StaticEndpoint(http.MethodGet, "/authz/roles")

// DeleteRoleRequest deletes a role by it's ID.
var DeleteRoleRequest = transports.IdentityEndpoint[string](http.MethodDelete, "/authz/roles/%s")

// GetAssignedUsersRequest retrieves IDs of users this role is assigned to.
// Use with [GetAssignedUsersResponse].
var GetAssignedUsersRequest = transports.IdentityEndpoint[string](http.MethodGet, "/authz/roles/%s/users")

type GetAssignedUsersResponse []string

// GetUserAssignmentsRequest retrieves IDs and type of users this role is assigned to.
// Use with [UserAssignment].
var GetUserAssignmentsRequest = transports.IdentityEndpoint[string](http.MethodGet, "/authz/roles/%s/user-assignments")

// GetUserAssignmentsRequest retrieves IDs and type of groups this role is assigned to.
// Use with [GroupAssignment].
var GetGroupAssignmentsRequest = transports.IdentityEndpoint[string](http.MethodGet, "/authz/roles/%s/group-assignments")

// oapi-codegen does not generate inline response types.
// https://github.com/oapi-codegen/oapi-codegen/issues/513
type (
	UserAssignment struct {
		ID   string `json:"userId"`
		Type string `json:"userType"`
	}
	GroupAssignment struct {
		ID   string `json:"groupId"`
		Type string `json:"groupType"`
	}
)

// AddPermissionsRequest adds permissions to a role.
type AddPermissionsRequest struct {
	transports.BaseEndpoint
	RoleID      string      `json:"-"`
	Permissions Permissions `json:"permissions"`
}

var _ transports.Endpoint = (*AddPermissionsRequest)(nil)

func (*AddPermissionsRequest) Method() string { return http.MethodPost }
func (r *AddPermissionsRequest) Path() string {
	return "/authz/roles/" + url.PathEscape(r.RoleID) + "/add-permissions"
}
func (r *AddPermissionsRequest) Body() any { return &r }

// RemovePermissionsRequest removes permissions from a role.
type RemovePermissionsRequest struct {
	transports.BaseEndpoint
	RoleID      string      `json:"-"`
	Permissions Permissions `json:"permissions"`
}

var _ transports.Endpoint = (*RemovePermissionsRequest)(nil)

func (*RemovePermissionsRequest) Method() string { return http.MethodPost }
func (r *RemovePermissionsRequest) Path() string {
	return "/authz/roles/" + url.PathEscape(r.RoleID) + "/remove-permissions"
}
func (r *RemovePermissionsRequest) Body() any { return &r }

// HasPermissionRequest checks if a role contains a permission.
type HasPermissionRequest struct {
	transports.BaseEndpoint

	RoleID      string
	Alias       AliasPermission
	Backups     BackupsPermission
	Cluster     ClusterPermission
	Collections CollectionPermission
	Data        DataPermission
	Groups      GroupPermission
	Namespaces  NamespacePermission
	Nodes       NodesPermission
	MCP         MCPPermission
	Replication ReplicationPermission
	Roles       RolePermission
	Tenants     TenantPermission
	Users       UserPermission
}

var (
	_ transports.Endpoint = (*HasPermissionRequest)(nil)
	_ json.Marshaler      = (*HasPermissionRequest)(nil)
)

func (*HasPermissionRequest) Method() string { return http.MethodPost }
func (r *HasPermissionRequest) Path() string {
	return "/authz/roles/" + url.PathEscape(r.RoleID) + "/has-permission"
}
func (r *HasPermissionRequest) Body() any { return &r }

func (r *Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          string       `json:"name"`
		Permissions *Permissions `json:"permissions"`
	}{
		ID:          r.ID,
		Permissions: &r.Permissions,
	})
}

type HasPermissionResponse bool

func (r *HasPermissionRequest) MarshalJSON() ([]byte, error) {
	var action rest.PermissionAction
	var data any

	switch {
	case r.Alias.Create:
		action, data = rest.CreateAliases, r.Alias
	case r.Alias.Read:
		action, data = rest.ReadAliases, r.Alias
	case r.Alias.Update:
		action, data = rest.UpdateAliases, r.Alias
	case r.Alias.Delete:
		action, data = rest.DeleteAliases, r.Alias
	case r.Backups.Manage:
		action, data = rest.ManageBackups, r.Backups
	case r.Cluster.Read:
		action, data = rest.ReadCluster, nil
	case r.Collections.Create:
		action, data = rest.CreateCollections, r.Collections
	case r.Collections.Read:
		action, data = rest.ReadCollections, r.Collections
	case r.Collections.Update:
		action, data = rest.UpdateCollections, r.Collections
	case r.Collections.Delete:
		action, data = rest.DeleteCollections, r.Collections
	case r.Data.Create:
		action, data = rest.CreateData, r.Data
	case r.Data.Read:
		action, data = rest.ReadData, r.Data
	case r.Data.Update:
		action, data = rest.UpdateData, r.Data
	case r.Data.Delete:
		action, data = rest.DeleteData, r.Data
	case r.Groups.Read:
		action, data = rest.ReadGroups, r.Groups
	case r.Groups.AssignAndRevoke:
		action, data = rest.AssignAndRevokeGroups, r.Groups
	case r.Namespaces.Manage:
		action, data = rest.ManageNamespaces, r.Namespaces
	case r.Nodes.Read:
		action, data = rest.ReadNodes, r.Nodes
	case r.MCP.Create:
		action = rest.CreateMcp
	case r.MCP.Read:
		action = rest.ReadMcp
	case r.MCP.Update:
		action = rest.UpdateMcp
	case r.Replication.Create:
		action, data = rest.CreateReplicate, r.Replication
	case r.Replication.Read:
		action, data = rest.ReadReplicate, r.Replication
	case r.Replication.Update:
		action, data = rest.UpdateReplicate, r.Replication
	case r.Replication.Delete:
		action, data = rest.DeleteReplicate, r.Replication
	case r.Roles.Create:
		action, data = rest.CreateRoles, r.Roles
	case r.Roles.Read:
		action, data = rest.ReadRoles, r.Roles
	case r.Roles.Update:
		action, data = rest.UpdateRoles, r.Roles
	case r.Roles.Delete:
		action, data = rest.DeleteRoles, r.Roles
	case r.Tenants.Create:
		action, data = rest.CreateTenants, r.Tenants
	case r.Tenants.Read:
		action, data = rest.ReadTenants, r.Tenants
	case r.Tenants.Update:
		action, data = rest.UpdateTenants, r.Tenants
	case r.Tenants.Delete:
		action, data = rest.DeleteTenants, r.Tenants
	case r.Users.Create:
		action, data = rest.CreateUsers, r.Users
	case r.Users.Read:
		action, data = rest.ReadUsers, r.Users
	case r.Users.Update:
		action, data = rest.UpdateUsers, r.Users
	case r.Users.Delete:
		action, data = rest.DeleteUsers, r.Users
	case r.Users.AssignAndRevoke:
		action, data = rest.AssignAndRevokeUsers, r.Users
	}

	permission := map[string]any{"action": action}
	if data != nil {
		if kind := permissionKind(action); kind != "" {
			permission[kind] = data
		}
	}
	return json.Marshal(permission)
}

// MarshalJSON flattens permissions into an array,
// creating an individual entry for every action.
func (ps *Permissions) MarshalJSON() ([]byte, error) {
	permissions := make([]map[string]any, 0)

	// Add permission with given action and resources.
	// The last part of the action string is the resource kind,
	// e.g. "aliases" in "create_aliases" or "groups" in "assign_and_revoke_groups".
	add := func(action rest.PermissionAction, data any) {
		permission := map[string]any{"action": action}
		if data != nil {
			if kind := permissionKind(action); kind != "" {
				permission[kind] = data
			}
		}
		permissions = append(permissions, permission)
	}

	for _, p := range ps.Aliases {
		if p.Create {
			add(rest.CreateAliases, p)
		}
		if p.Read {
			add(rest.ReadAliases, p)
		}
		if p.Update {
			add(rest.UpdateAliases, p)
		}
		if p.Delete {
			add(rest.DeleteAliases, p)
		}
	}

	for _, p := range ps.Backups {
		if p.Manage {
			add(rest.ManageBackups, p)
		}
	}

	for _, p := range ps.Cluster {
		if p.Read {
			add(rest.ReadCluster, nil)
		}
	}

	for _, p := range ps.Collections {
		if p.Create {
			add(rest.CreateCollections, p)
		}
		if p.Read {
			add(rest.ReadCollections, p)
		}
		if p.Update {
			add(rest.UpdateCollections, p)
		}
		if p.Delete {
			add(rest.DeleteCollections, p)
		}
	}

	for _, p := range ps.Data {
		if p.Create {
			add(rest.CreateData, p)
		}
		if p.Read {
			add(rest.ReadData, p)
		}
		if p.Update {
			add(rest.UpdateData, p)
		}
		if p.Delete {
			add(rest.DeleteData, p)
		}
	}

	for _, p := range ps.Groups {
		if p.Read {
			add(rest.ReadGroups, p)
		}
		if p.AssignAndRevoke {
			add(rest.AssignAndRevokeGroups, p)
		}
	}

	for _, p := range ps.Namespaces {
		if p.Manage {
			add(rest.ManageNamespaces, p)
		}
	}

	for _, p := range ps.Nodes {
		if p.Read {
			add(rest.ReadNodes, p)
		}
	}

	for _, p := range ps.MCP {
		if p.Create {
			add(rest.CreateMcp, nil)
		}
		if p.Read {
			add(rest.ReadMcp, nil)
		}
		if p.Update {
			add(rest.UpdateMcp, nil)
		}
	}

	for _, p := range ps.Replication {
		if p.Create {
			add(rest.CreateReplicate, p)
		}
		if p.Read {
			add(rest.ReadReplicate, p)
		}
		if p.Update {
			add(rest.UpdateReplicate, p)
		}
		if p.Delete {
			add(rest.DeleteReplicate, p)
		}
	}

	for _, p := range ps.Roles {
		if p.Create {
			add(rest.CreateRoles, p)
		}
		if p.Read {
			add(rest.ReadRoles, p)
		}
		if p.Update {
			add(rest.UpdateRoles, p)
		}
		if p.Delete {
			add(rest.DeleteRoles, p)
		}
	}

	for _, p := range ps.Tenants {
		if p.Create {
			add(rest.CreateTenants, p)
		}
		if p.Read {
			add(rest.ReadTenants, p)
		}
		if p.Update {
			add(rest.UpdateTenants, p)
		}
		if p.Delete {
			add(rest.DeleteTenants, p)
		}
	}

	for _, p := range ps.Users {
		if p.Create {
			add(rest.CreateUsers, p)
		}
		if p.Read {
			add(rest.ReadUsers, p)
		}
		if p.Update {
			add(rest.UpdateUsers, p)
		}
		if p.Delete {
			add(rest.DeleteUsers, p)
		}
		if p.AssignAndRevoke {
			add(rest.AssignAndRevokeUsers, p)
		}
	}

	return json.Marshal(permissions)
}

// UnmarshalJSON groups role's permissions by resources.
func (r *Role) UnmarshalJSON(data []byte) error {
	// [rest.Role] does not easily lend itself to merging,
	// we get more control by using our own data type and
	// unmarshaling concrete permissions lazily.
	var role struct {
		ID          string `json:"name"`
		Permissions []struct {
			Action      rest.PermissionAction `json:"action"`
			Alias       json.RawMessage       `json:"aliases"`
			Backup      json.RawMessage       `json:"backups"`
			Cluster     json.RawMessage       `json:"cluster"`
			Collection  json.RawMessage       `json:"collections"`
			Data        json.RawMessage       `json:"data"`
			Group       json.RawMessage       `json:"groups"`
			Namespace   json.RawMessage       `json:"namespaces"`
			Node        json.RawMessage       `json:"nodes"`
			MCP         json.RawMessage       `json:"mcp"`
			Replication json.RawMessage       `json:"replicate"`
			Role        json.RawMessage       `json:"roles"`
			Tenant      json.RawMessage       `json:"tenants"`
			User        json.RawMessage       `json:"users"`
		} `json:"permissions"`
	}

	if err := json.Unmarshal(data, &role); err != nil {
		return err
	}

	*r = Role{ID: role.ID}

	// Set of permissions with unique resource identifiers.
	// The key is a concatenation of <kind> and <data>, such
	// that mulitple permissions for the same resource are
	// represented as a single instance. The value is the
	// permission itself, e.g. [AliasPermission].
	lookup := make(map[string]any)
	for _, p := range role.Permissions {
		switch p.Action {
		case rest.CreateAliases:
			find(lookup, p.Action, p.Alias, &r.Aliases).Create = true
		case rest.ReadAliases:
			find(lookup, p.Action, p.Alias, &r.Aliases).Read = true
		case rest.UpdateAliases:
			find(lookup, p.Action, p.Alias, &r.Aliases).Update = true
		case rest.DeleteAliases:
			find(lookup, p.Action, p.Alias, &r.Aliases).Delete = true
		case rest.ManageBackups:
			find(lookup, p.Action, p.Backup, &r.Backups).Manage = true
		case rest.ReadCluster:
			find(lookup, p.Action, p.Cluster, &r.Cluster).Read = true
		case rest.CreateCollections:
			find(lookup, p.Action, p.Collection, &r.Collections).Create = true
		case rest.ReadCollections:
			find(lookup, p.Action, p.Collection, &r.Collections).Read = true
		case rest.UpdateCollections:
			find(lookup, p.Action, p.Collection, &r.Collections).Update = true
		case rest.DeleteCollections:
			find(lookup, p.Action, p.Collection, &r.Collections).Delete = true
		case rest.CreateData:
			find(lookup, p.Action, p.Data, &r.Data).Create = true
		case rest.ReadData:
			find(lookup, p.Action, p.Data, &r.Data).Read = true
		case rest.UpdateData:
			find(lookup, p.Action, p.Data, &r.Data).Update = true
		case rest.DeleteData:
			find(lookup, p.Action, p.Data, &r.Data).Delete = true
		case rest.ReadGroups:
			find(lookup, p.Action, p.Group, &r.Groups).Read = true
		case rest.AssignAndRevokeGroups:
			find(lookup, p.Action, p.Group, &r.Groups).AssignAndRevoke = true
		case rest.ManageNamespaces:
			find(lookup, p.Action, p.Namespace, &r.Namespaces).Manage = true
		case rest.ReadNodes:
			find(lookup, p.Action, p.Node, &r.Nodes).Read = true
		case rest.CreateMcp:
			find(lookup, p.Action, p.MCP, &r.MCP).Create = true
		case rest.ReadMcp:
			find(lookup, p.Action, p.MCP, &r.MCP).Read = true
		case rest.UpdateMcp:
			find(lookup, p.Action, p.MCP, &r.MCP).Update = true
		case rest.CreateReplicate:
			find(lookup, p.Action, p.Replication, &r.Replication).Create = true
		case rest.ReadReplicate:
			find(lookup, p.Action, p.Replication, &r.Replication).Read = true
		case rest.UpdateReplicate:
			find(lookup, p.Action, p.Replication, &r.Replication).Update = true
		case rest.DeleteReplicate:
			find(lookup, p.Action, p.Replication, &r.Replication).Delete = true
		case rest.CreateRoles:
			find(lookup, p.Action, p.Role, &r.Roles).Create = true
		case rest.ReadRoles:
			find(lookup, p.Action, p.Role, &r.Roles).Read = true
		case rest.UpdateRoles:
			find(lookup, p.Action, p.Role, &r.Roles).Update = true
		case rest.DeleteRoles:
			find(lookup, p.Action, p.Role, &r.Roles).Delete = true
		case rest.CreateTenants:
			find(lookup, p.Action, p.Tenant, &r.Tenants).Create = true
		case rest.ReadTenants:
			find(lookup, p.Action, p.Tenant, &r.Tenants).Read = true
		case rest.UpdateTenants:
			find(lookup, p.Action, p.Tenant, &r.Tenants).Update = true
		case rest.DeleteTenants:
			find(lookup, p.Action, p.Tenant, &r.Tenants).Delete = true
		case rest.CreateUsers:
			find(lookup, p.Action, p.User, &r.Users).Create = true
		case rest.ReadUsers:
			find(lookup, p.Action, p.User, &r.Users).Read = true
		case rest.UpdateUsers:
			find(lookup, p.Action, p.User, &r.Users).Update = true
		case rest.DeleteUsers:
			find(lookup, p.Action, p.User, &r.Users).Delete = true
		case rest.AssignAndRevokeUsers:
			find(lookup, p.Action, p.User, &r.Users).AssignAndRevoke = true
		}
	}

	return nil
}

func permissionKind(action rest.PermissionAction) string {
	if parts := strings.Split(string(action), "_"); len(parts) != 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// find retrieves the data for the permission for the given resource.
// The key is a concatenation of permission <kind> and <data>. The kind
// is assumed to be the last part of the action, as per the naming convention.
// If no permission for the given key exists, one is created by unmarshaling
// data into a new instance of T and is appended to the dest.
//
// If kind could not be found (action contains no "_") or data cannot be
// unmarshaled, find returns new(T), such that it is always safe to use
// the return value without a nil-check. Some permissions might not have
// any resource description, i.e. data is empty; for these, the unmarshaling
// is skipped, and a zero value of T is appended to dest.
func find[T any](lookup map[string]any, action rest.PermissionAction, data json.RawMessage, dest *[]T) *T {
	kind := permissionKind(action)
	if kind == "" {
		return new(T)
	}

	key := kind + string(data)
	addr, ok := lookup[key]
	if !ok {
		var t T
		if len(data) > 0 {
			if err := json.Unmarshal(data, &t); err != nil {
				return new(T)
			}
		}
		*dest = append(*dest, t)
		addr = &(*dest)[len(*dest)-1]
		lookup[key] = addr
	}
	return addr.(*T)
}
