package rbac

import "github.com/weaviate/weaviate/entities/models"

type Role struct {
	Permissions
	Name string
}

// Permissions holds for all permissions for all resources associated with the role.
type Permissions struct {
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
func NewRole(name string, permissionGroups ...PermissionGroup) Role {
	role := Role{Name: name}

	for _, group := range permissionGroups {
		group.ExtendRole(&role)
	}
	return role
}

func BackupPermissions(collection string, actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return BackupsPermission{
			Action:     action,
			Collection: collection,
		}
	}, actions...)
}

func ClusterPermissions(actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return ClusterPermission{Action: action}
	}, actions...)
}

func CollectionsPermissions(collection string, actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return CollectionsPermission{Action: action, Collection: collection}
	}, actions...)
}

func DataPermissions(collection, object string, actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return DataPermission{
			Action:     action,
			Collection: collection,
			Object:     object,
		}
	}, actions...)
}

// NodesPermissionsVerbose creates nodes permission for a specified collection.
// Verbosity is set to "verbose".
func NodesPermissionsVerbose(collection string, actions ...string) PermissionGroup {
	return nodesPermissions(collection, models.PermissionNodesVerbosityVerbose, actions...)
}

// NodesPermissionsMinimal creates nodes permission for all collections.
// Verbosity is set to "minimal".
func NodesPermissionsMinimal(collection string, actions ...string) PermissionGroup {
	return nodesPermissions("*", models.PermissionNodesVerbosityMinimal, actions...)
}

func nodesPermissions(collection, verbosity string, actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return NodesPermission{
			Action:     action,
			Collection: collection,
			Verbosity:  verbosity,
		}
	}, actions...)
}

func RolesPermissions(role, scope string, actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return RolesPermission{
			Action: action,
			Role:   role,
			Scope:  scope,
		}
	}, actions...)
}

func TenantsPermissions(actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return TenantsPermission{Action: action}
	}, actions...)
}

func UsersPermissions(actions ...string) PermissionGroup {
	return newPermissionGroup(func(action string) RoleOption {
		return UsersPermission{Action: action}
	}, actions...)
}

// RoleOption configures a role.
type RoleOption interface {
	ExtendRole(*Role) // Extend the role with permissions.
}

// PermissionGroup describes a collection of permissions for the same resource.
type PermissionGroup []RoleOption

var _ RoleOption = (PermissionGroup)(nil)

// newPermissionGroup creates permissions for multiple actions for the same resource.
// permFunc can return concrete permission types, which implement RoleOption.
func newPermissionGroup(permFunc func(action string) RoleOption, actions ...string) PermissionGroup {
	permissions := make(PermissionGroup, len(actions))
	for i, action := range actions {
		permissions[i] = permFunc(action)
	}
	return permissions
}

func (pg PermissionGroup) ExtendRole(r *Role) {
	for _, perm := range pg {
		perm.ExtendRole(r)
	}
}

type BackupsPermission struct {
	Action     string
	Collection string
}

func (p BackupsPermission) ExtendRole(r *Role) {
	r.Backups = append(r.Backups, p)
}

func (p BackupsPermission) toWeaviate() *models.Permission {
	return &models.Permission{
		Action: &p.Action,
		Backups: &models.PermissionBackups{
			Collection: &p.Collection,
		},
	}
}

type ClusterPermission struct {
	Action string
}

func (p ClusterPermission) ExtendRole(r *Role) {
	r.Cluster = append(r.Cluster, p)
}

func (p ClusterPermission) toWeaviate() *models.Permission {
	return &models.Permission{Action: &p.Action}
}

type CollectionsPermission struct {
	Action     string
	Collection string
}

func (p CollectionsPermission) ExtendRole(r *Role) {
	r.Collections = append(r.Collections, p)
}

func (p CollectionsPermission) toWeaviate() *models.Permission {
	return &models.Permission{
		Action: &p.Action,
		Collections: &models.PermissionCollections{
			Collection: &p.Collection,
		},
	}
}

type DataPermission struct {
	Action     string
	Collection string
	Object     string
}

func (p DataPermission) ExtendRole(r *Role) {
	r.Data = append(r.Data, p)
}

func (p DataPermission) toWeaviate() *models.Permission {
	return &models.Permission{
		Action: &p.Action,
		Data: &models.PermissionData{
			Collection: &p.Collection,
			Object:     &p.Object,
		},
	}
}

type NodesPermission struct {
	Action     string
	Collection string
	Verbosity  string
}

func (p NodesPermission) ExtendRole(r *Role) {
	r.Nodes = append(r.Nodes, p)
}

func (p NodesPermission) toWeaviate() *models.Permission {
	return &models.Permission{
		Action: &p.Action,
		Nodes: &models.PermissionNodes{
			Collection: &p.Collection,
			Verbosity:  &p.Verbosity,
		},
	}
}

type RolesPermission struct {
	Action string
	Role   string
	Scope  string
}

func (p RolesPermission) ExtendRole(r *Role) {
	r.Roles = append(r.Roles, p)
}

func (p RolesPermission) toWeaviate() *models.Permission {
	return &models.Permission{
		Action: &p.Action,
		Roles: &models.PermissionRoles{
			Role:  &p.Role,
			Scope: &p.Scope,
		},
	}
}

type TenantsPermission struct {
	Action string
}

func (p TenantsPermission) ExtendRole(r *Role) {
	r.Tenants = append(r.Tenants, p)
}

func (p TenantsPermission) toWeaviate() *models.Permission {
	return &models.Permission{Action: &p.Action}
}

type UsersPermission struct {
	Action string
}

func (p UsersPermission) ExtendRole(r *Role) {
	r.Users = append(r.Users, p)
}

func (p UsersPermission) toWeaviate() *models.Permission {
	return &models.Permission{Action: &p.Action}
}

// toWeaviate converts Permissions to the REST API format.
func (p Permissions) toWeaviate() []*models.Permission {
	var out []*models.Permission

	appendPermissions := func(n int, toWeaviate func(int) *models.Permission) {
		for i := 0; i < n; i++ {
			out = append(out, toWeaviate(i))
		}
	}

	appendPermissions(len(p.Backups), func(i int) *models.Permission { return p.Backups[i].toWeaviate() })

	appendPermissions(len(p.Cluster), func(i int) *models.Permission { return p.Cluster[i].toWeaviate() })
	appendPermissions(len(p.Collections), func(i int) *models.Permission { return p.Collections[i].toWeaviate() })
	appendPermissions(len(p.Data), func(i int) *models.Permission { return p.Data[i].toWeaviate() })
	appendPermissions(len(p.Nodes), func(i int) *models.Permission { return p.Nodes[i].toWeaviate() })
	appendPermissions(len(p.Roles), func(i int) *models.Permission { return p.Roles[i].toWeaviate() })
	appendPermissions(len(p.Tenants), func(i int) *models.Permission { return p.Tenants[i].toWeaviate() })
	appendPermissions(len(p.Users), func(i int) *models.Permission { return p.Users[i].toWeaviate() })
	return out
}
