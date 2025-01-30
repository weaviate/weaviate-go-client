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

func (api *API) Creator() *RoleCreator {
	return &RoleCreator{connection: api.connection}
}

func (api *API) Deleter() *RoleDeleter {
	return &RoleDeleter{connection: api.connection}
}

func (api *API) PermissionAdder() *PermissionAdder {
	return &PermissionAdder{connection: api.connection}
}

func (api *API) PermissionRemover() *PermissionRemover {
	return &PermissionRemover{connection: api.connection}
}

func (api *API) PermissionChecker() *PermissionChecker {
	return &PermissionChecker{connection: api.connection}
}

func (api *API) AllGetter() *RoleAllGetter {
	return &RoleAllGetter{connection: api.connection}
}

func (api *API) Getter() *RoleGetter {
	return &RoleGetter{connection: api.connection}
}

func (api *API) UserRolesGetter() *UserRolesGetter {
	return &UserRolesGetter{connection: api.connection}
}

func (api *API) AssignedUsersGetter() *AssignedUsersGetter {
	return &AssignedUsersGetter{connection: api.connection}
}

func (api *API) Exists() *RoleExists {
	return &RoleExists{connection: api.connection, getter: api.Getter()}
}

func (api *API) Assigner() *RoleAssigner {
	return &RoleAssigner{connection: api.connection}
}

func (api *API) Revoker() *RoleRevoker {
	return &RoleRevoker{connection: api.connection}
}
