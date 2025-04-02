package users

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate/entities/models"
)

type UserType string

const (
	UserTypeDB   UserType = UserType(models.UserTypeDb)
	UserTypeOIDC UserType = UserType(models.UserTypeOidc)
)

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
//
// Deprecated: This method is deprecated and will be removed in Q4 25.
// Please use DB().UserRolesGetter() and/or OIDC().UserRolesGetter() instead.
func (api *API) UserRolesGetter() *UserRolesGetter {
	return &UserRolesGetter{connection: api.connection}
}

// Assign a role to a user. Note that 'root' cannot be assigned.
//
// Deprecated: This method is deprecated and will be removed in Q4 25.
// Please use DB().Assigner() and/or OIDC().Assigner() instead.
func (api *API) Assigner() *RoleAssigner {
	return &RoleAssigner{connection: api.connection}
}

// Revoke a role from a user. Note that 'root' cannot be revoked.
//
// Deprecated: This method is deprecated and will be removed in Q4 25.
// Please use DB().Revoker() and/or OIDC().Revoker() instead.
func (api *API) Revoker() *RoleRevoker {
	return &RoleRevoker{connection: api.connection}
}

func (api *API) DB() *UserOperationsDB {
	return &UserOperationsDB{api.connection}
}

func (api *API) OIDC() *UserOperationsOIDC {
	return &UserOperationsOIDC{api.connection}
}
