package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type Role struct {
	ID          string
	Permissions Permissions
}

var (
	_ json.Marshaler   = (*Role)(nil)
	_ json.Unmarshaler = (*Role)(nil)
)

type Permissions struct {
	Aliases     []AliasPermission
	Backups     []BackupsPermission
	Cluster     []ClusterPermission
	Collections []CollectionsPermission
	Data        []DataPermission
	Groups      []GroupsPermission
	Namespaces  []NamespacePermission
	Nodes       []NodesPermission
	Replication []ReplicationPermission
	Roles       []RolesPermission
	Tenants     []TenantsPermission
	Users       []UsersPermission
}

type CreateRoleRequest struct {
	transports.BaseEndpoint
	Role Role
}

var _ transports.Endpoint = (*CreateRoleRequest)(nil)

func (r *CreateRoleRequest) Method() string { return http.MethodPost }
func (r *CreateRoleRequest) Path() string   { return "/authz/roles" }
func (r *CreateRoleRequest) Body() any      { return &r.Role }

var GetRoleRequest = transports.IdentityEndpoint[string](http.MethodGet, "/authz/roles/%s")

var ListRolesRequest = transports.StaticEndpoint(http.MethodGet, "/authz/roles")

var DeleteRoleRequest = transports.IdentityEndpoint[string](http.MethodDelete, "/authz/roles/%s")

type Permission struct {
	Alias       AliasPermission
	Backups     BackupsPermission
	Cluster     ClusterPermission
	Collections CollectionsPermission
	Data        DataPermission
	Groups      GroupsPermission
	Namespaces  NamespacePermission
	Nodes       NodesPermission
	Replication ReplicationPermission
	Roles       RolesPermission
	Tenants     TenantsPermission
	Users       UsersPermission
}

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
	CollectionsPermission struct {
		Collection string `json:"collection,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	DataPermission struct {
		Collection string `json:"collection,omitempty"`
		Tenant     string `json:"tenant,omitempty"`
		Object     string `json:"object,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	GroupsPermission struct {
		GroupID string `json:"group,omitempty"`
		Type    string `json:"groupType,omitempty"`

		Read            bool `json:"-"`
		AssignAndRevoke bool `json:"-"`
	}
	NamespacePermission struct {
		Namespace string `json:"namespace,omitempty"`

		Manage bool `json:"-"`
	}
	NodesVerbosity  rest.PermissionNodesVerbosity // api.NodeVerbosity
	NodesPermission struct {
		Collection string         `json:"collection,omitempty"`
		Verbosity  NodesVerbosity `json:"verbosity,omitempty"`

		Read bool `json:"-"`
	}
	ReplicationPermission struct {
		Collection string `json:"collection,omitempty"`
		Shard      string `json:"shard,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	RolesScope      rest.PermissionRolesScope
	RolesPermission struct {
		RoleID string     `json:"role,omitempty"`
		Scope  RolesScope `json:"scope,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	TenantsPermission struct {
		Collection string `json:"collection,omitempty"`
		Tenant     string `json:"tenant,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
	UsersPermission struct {
		UserID string `json:"user,omitempty"`

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
)

const (
	NodeVerbosityMinimal = NodesVerbosity(rest.Minimal)
	NodeVerbosityVerbose = NodesVerbosity(rest.Verbose)

	RolesScopeAll   = RolesScope(rest.PermissionRolesScopeAll)
	RolesScopeMatch = RolesScope(rest.PermissionRolesScopeMatch)
)

func (r *Role) MarshalJSON() ([]byte, error) {
	permissions := make([]map[string]any, 0)
	for _, p := range r.Permissions.Aliases {
		if p.Create {
			permissions = append(permissions, map[string]any{
				"action":  rest.CreateAliases,
				"aliases": p,
			})
		}
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action":  rest.ReadAliases,
				"aliases": p,
			})
		}
		if p.Update {
			permissions = append(permissions, map[string]any{
				"action":  rest.UpdateAliases,
				"aliases": p,
			})
		}
		if p.Delete {
			permissions = append(permissions, map[string]any{
				"action":  rest.DeleteAliases,
				"aliases": p,
			})
		}
	}

	for _, p := range r.Permissions.Backups {
		if p.Manage {
			permissions = append(permissions, map[string]any{
				"action":  rest.ManageBackups,
				"backups": p,
			})
		}
	}

	for _, p := range r.Permissions.Cluster {
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action": rest.ReadCluster,
			})
		}
	}

	for _, p := range r.Permissions.Collections {
		if p.Create {
			permissions = append(permissions, map[string]any{
				"action":      rest.CreateCollections,
				"collections": p,
			})
		}
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action":      rest.ReadCollections,
				"collections": p,
			})
		}
		if p.Update {
			permissions = append(permissions, map[string]any{
				"action":      rest.UpdateCollections,
				"collections": p,
			})
		}
		if p.Delete {
			permissions = append(permissions, map[string]any{
				"action":      rest.DeleteCollections,
				"collections": p,
			})
		}
	}

	for _, p := range r.Permissions.Data {
		if p.Create {
			permissions = append(permissions, map[string]any{
				"action": rest.CreateData,
				"data":   p,
			})
		}
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action": rest.ReadData,
				"data":   p,
			})
		}
		if p.Update {
			permissions = append(permissions, map[string]any{
				"action": rest.UpdateData,
				"data":   p,
			})
		}
		if p.Delete {
			permissions = append(permissions, map[string]any{
				"action": rest.DeleteData,
				"data":   p,
			})
		}
	}
	for _, p := range r.Permissions.Groups {
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action": rest.ReadGroups,
				"groups": p,
			})
		}
		if p.AssignAndRevoke {
			permissions = append(permissions, map[string]any{
				"action": rest.AssignAndRevokeGroups,
				"groups": p,
			})
		}
	}

	for _, p := range r.Permissions.Namespaces {
		if p.Manage {
			permissions = append(permissions, map[string]any{
				"action":     rest.ManageNamespaces,
				"namespaces": p,
			})
		}
	}

	for _, p := range r.Permissions.Nodes {
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action": rest.ReadNodes,
				"nodes":  p,
			})
		}
	}

	for _, p := range r.Permissions.Replication {
		if p.Create {
			permissions = append(permissions, map[string]any{
				"action":    rest.CreateReplicate,
				"replicate": p,
			})
		}
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action":    rest.ReadReplicate,
				"replicate": p,
			})
		}
		if p.Update {
			permissions = append(permissions, map[string]any{
				"action":    rest.UpdateReplicate,
				"replicate": p,
			})
		}
		if p.Delete {
			permissions = append(permissions, map[string]any{
				"action":    rest.DeleteReplicate,
				"replicate": p,
			})
		}
	}

	for _, p := range r.Permissions.Roles {
		if p.Create {
			permissions = append(permissions, map[string]any{
				"action": rest.CreateRoles,
				"roles":  p,
			})
		}
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action": rest.ReadRoles,
				"roles":  p,
			})
		}
		if p.Update {
			permissions = append(permissions, map[string]any{
				"action": rest.UpdateRoles,
				"roles":  p,
			})
		}
		if p.Delete {
			permissions = append(permissions, map[string]any{
				"action": rest.DeleteRoles,
				"roles":  p,
			})
		}
	}

	for _, p := range r.Permissions.Tenants {
		if p.Create {
			permissions = append(permissions, map[string]any{
				"action":  rest.CreateTenants,
				"tenants": p,
			})
		}
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action":  rest.ReadTenants,
				"tenants": p,
			})
		}
		if p.Update {
			permissions = append(permissions, map[string]any{
				"action":  rest.UpdateTenants,
				"tenants": p,
			})
		}
		if p.Delete {
			permissions = append(permissions, map[string]any{
				"action":  rest.DeleteTenants,
				"tenants": p,
			})
		}
	}

	for _, p := range r.Permissions.Users {
		if p.Create {
			permissions = append(permissions, map[string]any{
				"action": rest.CreateUsers,
				"users":  p,
			})
		}
		if p.Read {
			permissions = append(permissions, map[string]any{
				"action": rest.ReadUsers,
				"users":  p,
			})
		}
		if p.Update {
			permissions = append(permissions, map[string]any{
				"action": rest.UpdateUsers,
				"users":  p,
			})
		}
		if p.Delete {
			permissions = append(permissions, map[string]any{
				"action": rest.DeleteUsers,
				"users":  p,
			})
		}
	}

	return json.Marshal(&struct {
		ID          string           `json:"name"`
		Permissions []map[string]any `json:"permissions"`
	}{
		ID:          r.ID,
		Permissions: permissions,
	})
}

func (r *Role) UnmarshalJSON(data []byte) error {
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
			Replication json.RawMessage       `json:"replicate"`
			Role        json.RawMessage       `json:"roles"`
			Tenant      json.RawMessage       `json:"tenants"`
			User        json.RawMessage       `json:"users"`
		}
	}

	if err := json.Unmarshal(data, &role); err != nil {
		return err
	}

	*r = Role{ID: role.ID}

	lookup := make(map[string]any)
	ps := &r.Permissions
	for _, p := range role.Permissions {
		var kind string
		if parts := strings.Split(string(p.Action), "_"); len(parts) != 0 {
			kind = parts[len(parts)-1]
		} else {
			continue
		}

		switch p.Action {
		case rest.CreateAliases:
			find(lookup, kind, p.Alias, &ps.Aliases).Create = true
		case rest.ReadAliases:
			find(lookup, kind, p.Alias, &ps.Aliases).Read = true
		case rest.UpdateAliases:
			find(lookup, kind, p.Alias, &ps.Aliases).Update = true
		case rest.DeleteAliases:
			find(lookup, kind, p.Alias, &ps.Aliases).Delete = true
		case rest.ManageBackups:
			find(lookup, kind, p.Backup, &ps.Backups).Manage = true
		case rest.ReadCluster:
			find(lookup, kind, p.Cluster, &ps.Cluster).Read = true
		case rest.CreateCollections:
			find(lookup, kind, p.Collection, &ps.Collections).Create = true
		case rest.ReadCollections:
			find(lookup, kind, p.Collection, &ps.Collections).Read = true
		case rest.UpdateCollections:
			find(lookup, kind, p.Collection, &ps.Collections).Update = true
		case rest.DeleteCollections:
			find(lookup, kind, p.Collection, &ps.Collections).Delete = true
		case rest.CreateData:
			find(lookup, kind, p.Data, &ps.Data).Create = true
		case rest.ReadData:
			find(lookup, kind, p.Data, &ps.Data).Read = true
		case rest.UpdateData:
			find(lookup, kind, p.Data, &ps.Data).Update = true
		case rest.DeleteData:
			find(lookup, kind, p.Data, &ps.Data).Delete = true
		case rest.ReadGroups:
			find(lookup, kind, p.Group, &ps.Groups).Read = true
		case rest.AssignAndRevokeGroups:
			find(lookup, kind, p.Group, &ps.Groups).AssignAndRevoke = true
		case rest.ManageNamespaces:
			find(lookup, kind, p.Namespace, &ps.Namespaces).Manage = true
		case rest.ReadNodes:
			find(lookup, kind, p.Node, &ps.Nodes).Read = true
		case rest.CreateReplicate:
			find(lookup, kind, p.Replication, &ps.Replication).Create = true
		case rest.ReadReplicate:
			find(lookup, kind, p.Replication, &ps.Replication).Read = true
		case rest.UpdateReplicate:
			find(lookup, kind, p.Replication, &ps.Replication).Update = true
		case rest.DeleteReplicate:
			find(lookup, kind, p.Replication, &ps.Replication).Delete = true
		case rest.CreateRoles:
			find(lookup, kind, p.Role, &ps.Roles).Create = true
		case rest.ReadRoles:
			find(lookup, kind, p.Role, &ps.Roles).Read = true
		case rest.UpdateRoles:
			find(lookup, kind, p.Role, &ps.Roles).Update = true
		case rest.DeleteRoles:
			find(lookup, kind, p.Role, &ps.Roles).Delete = true
		case rest.CreateTenants:
			find(lookup, kind, p.Tenant, &ps.Tenants).Create = true
		case rest.ReadTenants:
			find(lookup, kind, p.Tenant, &ps.Tenants).Read = true
		case rest.UpdateTenants:
			find(lookup, kind, p.Tenant, &ps.Tenants).Update = true
		case rest.DeleteTenants:
			find(lookup, kind, p.Tenant, &ps.Tenants).Delete = true
		case rest.CreateUsers:
			find(lookup, kind, p.User, &ps.Users).Create = true
		case rest.ReadUsers:
			find(lookup, kind, p.User, &ps.Users).Read = true
		case rest.UpdateUsers:
			find(lookup, kind, p.User, &ps.Users).Update = true
		case rest.DeleteUsers:
			find(lookup, kind, p.User, &ps.Users).Delete = true
		}
	}

	return nil
}

func find[T any](lookup map[string]any, kind string, data json.RawMessage, dest *[]T) *T {
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
