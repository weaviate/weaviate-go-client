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
	role := Role{Name: *r.Name}

	for _, perm := range r.Permissions {
		switch {
		case perm.Backups != nil:
			role.Backups = append(role.Backups, BackupsPermission{
				Action:     *perm.Action,
				Collection: *perm.Backups.Collection,
			})
		case perm.Collections != nil:
			role.Collections = append(role.Collections, CollectionsPermission{
				Action:     *perm.Action,
				Collection: *perm.Collections.Collection,
			})
		case perm.Data != nil:
			role.Data = append(role.Data, DataPermission{
				Action:     *perm.Action,
				Collection: *perm.Data.Collection,
				Object:     *perm.Data.Object,
			})
		case perm.Nodes != nil:
			role.Nodes = append(role.Nodes, NodesPermission{
				Action:     *perm.Action,
				Collection: *perm.Nodes.Collection,
				Verbosity:  *perm.Nodes.Verbosity,
			})
		case perm.Roles != nil:
			role.Roles = append(role.Roles, RolesPermission{
				Action: *perm.Action,
				Role:   *perm.Roles.Role,
				Scope:  *perm.Roles.Scope,
			})

		// Weaviate v1.30 may defined additional actions for these permission groups
		// and we want to ensure they can be handled elegantly.
		// While somewhat crude, this method makes sure any cluster/tenants/users
		// action are read correctly without requiring the latest client version.
		case strings.HasSuffix(*perm.Action, "cluster"):
			role.Cluster = append(role.Cluster, ClusterPermission{Action: *perm.Action})
		case strings.HasSuffix(*perm.Action, "tenants"):
			role.Tenants = append(role.Tenants, TenantsPermission{Action: *perm.Action})
		case strings.HasSuffix(*perm.Action, "users"):
			role.Users = append(role.Users, UsersPermission{Action: *perm.Action})
		default:
			// New permission group may have been introduced on the server,
			// e.g. "manage_indices", which aren't reflected in this version of the client,
			// so it doesn't have a good way of presenting them to the user.
			log.Printf("WARN: %q action belongs to an unrecognized group, try updating the client to the latest version", *perm.Action)
		}
	}
	return role
}
