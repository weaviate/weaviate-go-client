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

	add := func(action rest.PermissionAction, data any) {
		permission := map[string]any{"action": action}
		if data != nil {
			var kind string
			if parts := strings.Split(string(action), "_"); len(parts) != 0 {
				kind = parts[len(parts)-1]
				permission[kind] = data
			}
		}
		permissions = append(permissions, permission)
	}

	for _, p := range r.Permissions.Aliases {
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

	for _, p := range r.Permissions.Backups {
		if p.Manage {
			add(rest.ManageBackups, p)
		}
	}

	for _, p := range r.Permissions.Cluster {
		if p.Read {
			add(rest.ReadCluster, nil)
		}
	}

	for _, p := range r.Permissions.Collections {
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

	for _, p := range r.Permissions.Data {
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

	for _, p := range r.Permissions.Groups {
		if p.Read {
			add(rest.ReadGroups, p)
		}
		if p.AssignAndRevoke {
			add(rest.AssignAndRevokeGroups, p)
		}
	}

	for _, p := range r.Permissions.Namespaces {
		if p.Manage {
			add(rest.ManageNamespaces, p)
		}
	}

	for _, p := range r.Permissions.Nodes {
		if p.Read {
			add(rest.ReadNodes, p)
		}
	}

	for _, p := range r.Permissions.Replication {
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

	for _, p := range r.Permissions.Roles {
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

	for _, p := range r.Permissions.Tenants {
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

	for _, p := range r.Permissions.Users {
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
