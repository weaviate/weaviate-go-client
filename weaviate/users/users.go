package users

import "github.com/weaviate/weaviate-go-client/v5/weaviate/connection"

type API struct {
	connection *connection.Connection
}

func New(connection *connection.Connection) *API {
	return &API{connection}
}

// Get user info for the current user.
func (api *API) MyUserGetter() *MyUserGetter {
	return &MyUserGetter{connection: api.connection}
}

// Get roles assigned to a user.
func (api *API) UserRolesGetter() *UserRolesGetter {
	return &UserRolesGetter{connection: api.connection}
}

// Assign a role to a user. Note that 'root' cannot be assigned.
func (api *API) Assigner() *RoleAssigner {
	return &RoleAssigner{connection: api.connection}
}

// Revoke a role from a user. Note that 'root' cannot be revoked.
func (api *API) Revoker() *RoleRevoker {
	return &RoleRevoker{connection: api.connection}
}
