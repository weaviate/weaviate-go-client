package rbac

import "github.com/weaviate/weaviate-go-client/v4/weaviate/connection"

const (
	BACKEND_FILESYSTEM = "filesystem"
	BACKEND_S3         = "s3"
	BACKEND_GCS        = "gcs"
	BACKEND_AZURE      = "azure"
)

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

// Get roles assigned to a user.
func (api *API) UserRolesGetter() *UserRolesGetter {
	return &UserRolesGetter{connection: api.connection}
}

// Get users assigned to a role.
func (api *API) AssignedUsersGetter() *AssignedUsersGetter {
	return &AssignedUsersGetter{connection: api.connection}
}

// Check if a role exists.
func (api *API) Exists() *RoleExists {
	return &RoleExists{connection: api.connection, getter: api.Getter()}
}

// Assign a role to a user. Note that 'root' cannot be assigned.
func (api *API) Assigner() *RoleAssigner {
	return &RoleAssigner{connection: api.connection}
}

// Revoke a role from a user. Note that 'root' cannot be revoked.
func (api *API) Revoker() *RoleRevoker {
	return &RoleRevoker{connection: api.connection}
}
