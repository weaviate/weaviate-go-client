package rbac

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleGetter struct {
	connection *connection.Connection

	name string
}

func (rg *RoleGetter) WithName(name string) *RoleGetter {
	rg.name = name
	return rg
}

func (rg *RoleGetter) Do(ctx context.Context) (Role, error) {
	res, err := rg.connection.RunREST(ctx, "/authz/roles/"+rg.name, http.MethodGet, nil)
	if err != nil {
		return Role{}, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var role models.Role
		decodeErr := res.DecodeBodyIntoTarget(&role)
		return roleFromWeaviate(role), decodeErr
	}
	return Role{}, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func roleFromWeaviate(r models.Role) Role {
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
			backups.Add(func(actions []string, filters ...string) Permission {
				return BackupsPermission{
					Actions:    actions,
					Collection: filters[0],
				}
			}, *perm.Action, *perm.Backups.Collection)
		case perm.Collections != nil:
			collections.Add(func(actions []string, filters ...string) Permission {
				return CollectionsPermission{
					Actions:    actions,
					Collection: filters[0],
				}
			}, *perm.Action, *perm.Collections.Collection)
		case perm.Data != nil:
			data.Add(func(actions []string, filters ...string) Permission {
				return DataPermission{
					Actions:    actions,
					Collection: filters[0],
				}
			}, *perm.Action, *perm.Data.Collection)
		case perm.Nodes != nil:
			nodes.Add(func(actions []string, filters ...string) Permission {
				return NodesPermission{
					Actions:    actions,
					Collection: filters[0],
					Verbosity:  filters[1],
				}
			}, *perm.Action, *perm.Nodes.Collection, *perm.Nodes.Verbosity)
			// Scope comes back as `nil` if not set.
		case perm.Roles != nil:
			// scope := ""
			// if (perm.Roles.Scope != nil) {
			// 	scope = *perm.Roles.Scope
			// }
			roles.Add(func(actions []string, filters ...string) Permission {
				return RolesPermission{
					Actions: actions,
					Role:    filters[0],
					Scope:   filters[1],
				}
			}, *perm.Action, *perm.Roles.Role, *perm.Roles.Scope)

		// Weaviate v1.30 may defined additional actions for these permission groups
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
	permFunc func(actions []string, filters ...string) Permission,
	action string, parts ...string,
) {
	key := strings.Join(parts, "#")
	if v, ok := mp[key]; !ok {
		mp[key] = &genericPermission{fields: parts, permFunc: permFunc}
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
	concrete := gp.permFunc(gp.actions, gp.fields...)
	concrete.ExtendRole(r)
}

// genericPermission is a helper type that has information necessary
type genericPermission struct {
	actions []string
	fields  []string

	// permFunc creates a concrete permission with given actions and filters.
	permFunc func(actions []string, filters ...string) Permission
}
