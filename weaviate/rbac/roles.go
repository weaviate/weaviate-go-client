package rbac

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate/entities/models"
)

type (
	UserType string
)

const (
	UserTypeDB    UserType = UserType(models.UserTypeOutputDbUser)
	UserTypeDBEnv UserType = UserType(models.UserTypeOutputDbEnvUser)
	UserTypeOIDC  UserType = UserType(models.UserTypeOutputOidc)
)

type UserAssignment struct {
	UserID   string
	UserType UserType
}

type API struct {
	connection *connection.Connection
}

func New(connection *connection.Connection) *API {
	return &API{connection}
}

// Create a new role.
func (api *API) Creator() *RoleCreator {
	return &RoleCreator{connection: api.connection}
}

// Delete a role.
func (api *API) Deleter() *RoleDeleter {
	return &RoleDeleter{connection: api.connection}
}

// Add permissions to an existing role.
// Note: This method is an upsert operation. If the permission already exists,
// it will be updated. If it does not exist, it will be created.
func (api *API) PermissionAdder() *PermissionAdder {
	return &PermissionAdder{connection: api.connection}
}

// Remove permissions from a role.
// Note: This method is a downsert operation. If the permission does not
// exist, it will be ignored. If these permissions are the only permissions of
// the role, the role will be deleted.
func (api *API) PermissionRemover() *PermissionRemover {
	return &PermissionRemover{connection: api.connection}
}

// Check if a role has a permission.
func (api *API) PermissionChecker() *PermissionChecker {
	return &PermissionChecker{connection: api.connection}
}

// Get all existing roles.
func (api *API) AllGetter() *RoleAllGetter {
	return &RoleAllGetter{connection: api.connection}
}

// Get role and its associated permissions.
func (api *API) Getter() *RoleGetter {
	return &RoleGetter{connection: api.connection}
}

// Get users assigned to a role.
//
// Deprecated: Use UserAssignmentGetter() instead.
func (api *API) AssignedUsersGetter() *AssignedUsersGetter {
	return &AssignedUsersGetter{connection: api.connection}
}

// Get users assigned to a role.
func (api *API) UserAssignmentGetter() *UserAssignmentGetter {
	return &UserAssignmentGetter{connection: api.connection}
}

// Check if a role exists.
func (api *API) Exists() *RoleExists {
	return &RoleExists{connection: api.connection, getter: api.Getter()}
}
