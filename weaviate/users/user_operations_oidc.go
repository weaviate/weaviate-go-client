package users

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate/entities/models"
)

type UserOperationsOIDC struct {
	connection *connection.Connection
}

func (r *UserOperationsOIDC) RolesGetter() *UserRolesGetter {
	return &UserRolesGetter{connection: r.connection, userType: models.UserTypeInputOidc}
}

func (r *UserOperationsOIDC) RolesAssigner() *RoleAssigner {
	return &RoleAssigner{connection: r.connection, userType: models.UserTypeInputOidc}
}

func (r *UserOperationsOIDC) RolesRevoker() *RoleRevoker {
	return &RoleRevoker{connection: r.connection, userType: models.UserTypeInputOidc}
}
