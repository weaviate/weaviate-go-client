package rbac

import "github.com/weaviate/weaviate/entities/models"

type Role struct {
	Name string

	Backups     []BackupsPermission
	Cluster     []ClusterPermission
	Collections []CollectionsPermission
	Data        []DataPermission
	Nodes       []NodesPermission
	Roles       []RolesPermission
	Tenants     []TenantsPermission
	Users       []UsersPermission
}

// NewRole creates a role with its associated permissions.
func NewRole(name string, permissions ...Permission) Role {
	role := Role{Name: name}

	for _, perm := range permissions {
		perm.ExtendRole(&role)
	}
	return role
}

// Permission configures a role.
type Permission interface {
	ExtendRole(*Role) // Extend the role with permissions.
}

type BackupsPermission struct {
	Actions    []string
	Collection string
}

func (p BackupsPermission) ExtendRole(r *Role) {
	r.Backups = append(r.Backups, p)
}

func (p BackupsPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{
			Action: &action,
			Backups: &models.PermissionBackups{
				Collection: &p.Collection,
			},
		}
	}
	return out
}

type ClusterPermission struct {
	Actions []string
}

func (p ClusterPermission) ExtendRole(r *Role) {
	r.Cluster = append(r.Cluster, p)
}

func (p ClusterPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{Action: &action}
	}
	return out
}

type CollectionsPermission struct {
	Actions    []string
	Collection string
}

func (p CollectionsPermission) ExtendRole(r *Role) {
	r.Collections = append(r.Collections, p)
}

func (p CollectionsPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{
			Action: &action,
			Collections: &models.PermissionCollections{
				Collection: &p.Collection,
			},
		}
	}
	return out
}

type DataPermission struct {
	Actions    []string
	Collection string
}

func (p DataPermission) ExtendRole(r *Role) {
	r.Data = append(r.Data, p)
}

func (p DataPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{
			Action: &action,
			Data: &models.PermissionData{
				Collection: &p.Collection,
			},
		}
	}
	return out
}

type NodesPermission struct {
	Actions    []string
	Collection string
	Verbosity  string
}

func (p NodesPermission) ExtendRole(r *Role) {
	r.Nodes = append(r.Nodes, p)
}

func (p NodesPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{
			Action: &action,
			Nodes: &models.PermissionNodes{
				Collection: &p.Collection,
				Verbosity:  &p.Verbosity,
			},
		}
	}
	return out
}

type RolesPermission struct {
	Actions []string
	Role    string
	Scope   string
}

func (p RolesPermission) ExtendRole(r *Role) {
	r.Roles = append(r.Roles, p)
}

func (p RolesPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{
			Action: &action,
			Roles: &models.PermissionRoles{
				Role:  &p.Role,
				Scope: &p.Scope,
			},
		}
	}
	return out
}

type TenantsPermission struct {
	Actions []string
}

func (p TenantsPermission) ExtendRole(r *Role) {
	r.Tenants = append(r.Tenants, p)
}

func (p TenantsPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{Action: &action}
	}
	return out
}

type UsersPermission struct {
	Actions []string
}

func (p UsersPermission) ExtendRole(r *Role) {
	r.Users = append(r.Users, p)
}

func (p UsersPermission) toWeaviate() []*models.Permission {
	out := make([]*models.Permission, len(p.Actions))
	for i, action := range p.Actions {
		out[i] = &models.Permission{Action: &action}
	}
	return out
}

// makeWeaviatePermissions converts Role's permissions to the REST API format.
func (r *Role) makeWeaviatePermissions() []*models.Permission {
	var out []*models.Permission

	appendPermissions := func(n int, toWeaviate func(int) []*models.Permission) {
		for i := 0; i < n; i++ {
			out = append(out, toWeaviate(i)...)
		}
	}

	appendPermissions(len(r.Backups), func(i int) []*models.Permission { return r.Backups[i].toWeaviate() })
	appendPermissions(len(r.Cluster), func(i int) []*models.Permission { return r.Cluster[i].toWeaviate() })
	appendPermissions(len(r.Collections), func(i int) []*models.Permission { return r.Collections[i].toWeaviate() })
	appendPermissions(len(r.Data), func(i int) []*models.Permission { return r.Data[i].toWeaviate() })
	appendPermissions(len(r.Nodes), func(i int) []*models.Permission { return r.Nodes[i].toWeaviate() })
	appendPermissions(len(r.Roles), func(i int) []*models.Permission { return r.Roles[i].toWeaviate() })
	appendPermissions(len(r.Tenants), func(i int) []*models.Permission { return r.Tenants[i].toWeaviate() })
	appendPermissions(len(r.Users), func(i int) []*models.Permission { return r.Users[i].toWeaviate() })
	return out
}
