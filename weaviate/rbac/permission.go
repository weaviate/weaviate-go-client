package rbac

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/weaviate/weaviate/entities/models"
)

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

var _ json.Unmarshaler = (*Role)(nil)

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

func (r *Role) UnmarshalJSON(data []byte) error {
	var tmp *models.Role
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	role := roleFromWeaviate(tmp)
	*r = role
	return nil
}

// roleFromWeaviate groups permissions by resource and creates an rbac.Role.
func roleFromWeaviate(r *models.Role) Role {
	backups := make(mergedPermissions)
	collections := make(mergedPermissions)
	data := make(mergedPermissions)
	nodes := make(mergedPermissions)
	roles := make(mergedPermissions)
	clusters := make(mergedPermissions)
	tenants := make(mergedPermissions)
	users := make(mergedPermissions)

	for _, perm := range r.Permissions {
		switch {
		case perm.Backups != nil:
			backups.Add(func(actions []string, resources ...string) Permission {
				return BackupsPermission{
					Actions:    actions,
					Collection: resources[0],
				}
			}, *perm.Action, *perm.Backups.Collection)
		case perm.Collections != nil:
			collections.Add(func(actions []string, resources ...string) Permission {
				return CollectionsPermission{
					Actions:    actions,
					Collection: resources[0],
				}
			}, *perm.Action, *perm.Collections.Collection)
		case perm.Data != nil:
			data.Add(func(actions []string, resources ...string) Permission {
				return DataPermission{
					Actions:    actions,
					Collection: resources[0],
				}
			}, *perm.Action, *perm.Data.Collection)
		case perm.Nodes != nil:
			nodes.Add(func(actions []string, resources ...string) Permission {
				return NodesPermission{
					Actions:    actions,
					Collection: resources[0],
					Verbosity:  resources[1],
				}
			}, *perm.Action, *perm.Nodes.Collection, *perm.Nodes.Verbosity)
		case perm.Roles != nil:
			roles.Add(func(actions []string, resources ...string) Permission {
				return RolesPermission{
					Actions: actions,
					Role:    resources[0],
					Scope:   resources[1],
				}
			}, *perm.Action, *perm.Roles.Role, *perm.Roles.Scope)

		// Weaviate v1.30 may define additional actions for these permission groups
		// and we want to ensure they can be handled elegantly.
		// While somewhat crude, this method makes sure any cluster/tenants/users
		// action are read correctly without requiring the latest client version.
		case strings.HasSuffix(*perm.Action, "cluster"):
			clusters.Add(func(actions []string, _ ...string) Permission {
				return ClusterPermission{Actions: actions}
			}, *perm.Action)
		case strings.HasSuffix(*perm.Action, "tenants"):
			tenants.Add(func(actions []string, _ ...string) Permission {
				return TenantsPermission{Actions: actions}
			}, *perm.Action)
		case strings.HasSuffix(*perm.Action, "users"):
			users.Add(func(actions []string, _ ...string) Permission {
				return UsersPermission{Actions: actions}
			}, *perm.Action)
		default:
			// New permission group may have been introduced on the server,
			// e.g. "manage_indices", which aren't reflected in this version of the client,
			// so it doesn't have a good way of presenting them to the user.
			log.Printf("WARN: %q action belongs to an unrecognized group, try updating the client to the latest version", *perm.Action)
		}
	}
	return NewRole(*r.Name, backups, collections, data, nodes, roles, clusters, tenants, users)
}

// mergedPermissions groups permissions by resource.
type mergedPermissions map[string]*genericPermission

func (mp mergedPermissions) Add(
	permFunc func(actions []string, resources ...string) Permission,
	action string, resources ...string,
) {
	key := strings.Join(resources, "#")
	if v, ok := mp[key]; !ok {
		mp[key] = &genericPermission{actions: []string{action}, resources: resources, permFunc: permFunc}
	} else {
		v.actions = append(v.actions, action)
	}
}

// ExtendRole with all permissions in this group.
func (mp mergedPermissions) ExtendRole(r *Role) {
	for _, generic := range mp {
		generic.ExtendRole(r)
	}
}

// ExtendRole with a concrete action derived from permFunc.
func (gp *genericPermission) ExtendRole(r *Role) {
	concrete := gp.permFunc(gp.actions, gp.resources...)
	concrete.ExtendRole(r)
}

// genericPermission is a helper container for information
// necessary to construct a concrete Permission for the specified resources.
type genericPermission struct {
	actions   []string
	resources []string

	// permFunc creates a concrete permission with given actions and filters.
	permFunc func(actions []string, resources ...string) Permission
}
