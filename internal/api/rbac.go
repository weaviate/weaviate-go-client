package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type Role struct {
	ID string
	Permissions
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
	Replication []ReplicationPermission
	Roles       []RolePermission
	Tenants     []TenantPermission
	Users       []UserPermission
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
	Collections CollectionPermission
	Data        DataPermission
	Groups      GroupPermission
	Namespaces  NamespacePermission
	Nodes       NodesPermission
	Replication ReplicationPermission
	Roles       RolePermission
	Tenants     TenantPermission
	Users       UserPermission
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
		Object     string `json:"object,omitempty"`

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
	NodeVerbosity   rest.PermissionNodesVerbosity // api.NodeVerbosity
	NodesPermission struct {
		Collection string        `json:"collection,omitempty"`
		Verbosity  NodeVerbosity `json:"verbosity,omitempty"`

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

		Create bool `json:"-"`
		Read   bool `json:"-"`
		Update bool `json:"-"`
		Delete bool `json:"-"`
	}
)

const (
	NodeVerbosityMinimal = NodeVerbosity(rest.Minimal)
	NodeVerbosityVerbose = NodeVerbosity(rest.Verbose)

	RoleScopeAll   = RoleScope(rest.PermissionRolesScopeAll)
	RoleScopeMatch = RoleScope(rest.PermissionRolesScopeMatch)
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

	for _, p := range r.Aliases {
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

	for _, p := range r.Backups {
		if p.Manage {
			add(rest.ManageBackups, p)
		}
	}

	for _, p := range r.Cluster {
		if p.Read {
			add(rest.ReadCluster, nil)
		}
	}

	for _, p := range r.Collections {
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

	for _, p := range r.Data {
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

	for _, p := range r.Groups {
		if p.Read {
			add(rest.ReadGroups, p)
		}
		if p.AssignAndRevoke {
			add(rest.AssignAndRevokeGroups, p)
		}
	}

	for _, p := range r.Namespaces {
		if p.Manage {
			add(rest.ManageNamespaces, p)
		}
	}

	for _, p := range r.Nodes {
		if p.Read {
			add(rest.ReadNodes, p)
		}
	}

	for _, p := range r.Replication {
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

	for _, p := range r.Roles {
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

	for _, p := range r.Tenants {
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

	for _, p := range r.Users {
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
		} `json:"permissions"`
	}

	if err := json.Unmarshal(data, &role); err != nil {
		return err
	}

	*r = Role{ID: role.ID}

	lookup := make(map[string]any)
	for _, p := range role.Permissions {
		var kind string
		if parts := strings.Split(string(p.Action), "_"); len(parts) != 0 {
			kind = parts[len(parts)-1]
		} else {
			continue
		}

		switch p.Action {
		case rest.CreateAliases:
			find(lookup, kind, p.Alias, &r.Aliases).Create = true
		case rest.ReadAliases:
			find(lookup, kind, p.Alias, &r.Aliases).Read = true
		case rest.UpdateAliases:
			find(lookup, kind, p.Alias, &r.Aliases).Update = true
		case rest.DeleteAliases:
			find(lookup, kind, p.Alias, &r.Aliases).Delete = true
		case rest.ManageBackups:
			find(lookup, kind, p.Backup, &r.Backups).Manage = true
		case rest.ReadCluster:
			find(lookup, kind, p.Cluster, &r.Cluster).Read = true
		case rest.CreateCollections:
			find(lookup, kind, p.Collection, &r.Collections).Create = true
		case rest.ReadCollections:
			find(lookup, kind, p.Collection, &r.Collections).Read = true
		case rest.UpdateCollections:
			find(lookup, kind, p.Collection, &r.Collections).Update = true
		case rest.DeleteCollections:
			find(lookup, kind, p.Collection, &r.Collections).Delete = true
		case rest.CreateData:
			find(lookup, kind, p.Data, &r.Data).Create = true
		case rest.ReadData:
			find(lookup, kind, p.Data, &r.Data).Read = true
		case rest.UpdateData:
			find(lookup, kind, p.Data, &r.Data).Update = true
		case rest.DeleteData:
			find(lookup, kind, p.Data, &r.Data).Delete = true
		case rest.ReadGroups:
			find(lookup, kind, p.Group, &r.Groups).Read = true
		case rest.AssignAndRevokeGroups:
			find(lookup, kind, p.Group, &r.Groups).AssignAndRevoke = true
		case rest.ManageNamespaces:
			find(lookup, kind, p.Namespace, &r.Namespaces).Manage = true
		case rest.ReadNodes:
			find(lookup, kind, p.Node, &r.Nodes).Read = true
		case rest.CreateReplicate:
			find(lookup, kind, p.Replication, &r.Replication).Create = true
		case rest.ReadReplicate:
			find(lookup, kind, p.Replication, &r.Replication).Read = true
		case rest.UpdateReplicate:
			find(lookup, kind, p.Replication, &r.Replication).Update = true
		case rest.DeleteReplicate:
			find(lookup, kind, p.Replication, &r.Replication).Delete = true
		case rest.CreateRoles:
			find(lookup, kind, p.Role, &r.Roles).Create = true
		case rest.ReadRoles:
			find(lookup, kind, p.Role, &r.Roles).Read = true
		case rest.UpdateRoles:
			find(lookup, kind, p.Role, &r.Roles).Update = true
		case rest.DeleteRoles:
			find(lookup, kind, p.Role, &r.Roles).Delete = true
		case rest.CreateTenants:
			find(lookup, kind, p.Tenant, &r.Tenants).Create = true
		case rest.ReadTenants:
			find(lookup, kind, p.Tenant, &r.Tenants).Read = true
		case rest.UpdateTenants:
			find(lookup, kind, p.Tenant, &r.Tenants).Update = true
		case rest.DeleteTenants:
			find(lookup, kind, p.Tenant, &r.Tenants).Delete = true
		case rest.CreateUsers:
			find(lookup, kind, p.User, &r.Users).Create = true
		case rest.ReadUsers:
			find(lookup, kind, p.User, &r.Users).Read = true
		case rest.UpdateUsers:
			find(lookup, kind, p.User, &r.Users).Update = true
		case rest.DeleteUsers:
			find(lookup, kind, p.User, &r.Users).Delete = true
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
