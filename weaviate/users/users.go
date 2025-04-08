package users

import (
	"encoding/json"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate/entities/models"
)

type UserTypeInput string

const (
	UserTypeInputDB   UserTypeInput = UserTypeInput(models.UserTypeInputDb)
	UserTypeInputOIDC UserTypeInput = UserTypeInput(models.UserTypeInputOidc)
)

type API struct {
	connection *connection.Connection
}

type UserInfo struct {
	Active    bool
	CreatedAt time.Time
	UserType  rbac.UserType
	UserID    string
	Roles     []*rbac.Role
}

func (info *UserInfo) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Active    bool            `json:"active"`
		CreatedAt strfmt.DateTime `json:"createdAt,omitempty"`
		UserType  string          `json:"dbUserType"`
		UserID    string          `json:"user_id"`
		Username  string          `json:"username"`
		Roles     []*rbac.Role    `json:"roles"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	id := tmp.UserID
	if id == "" {
		id = tmp.Username
	}

	*info = UserInfo{
		Active:    tmp.Active,
		CreatedAt: time.Time(tmp.CreatedAt),
		UserType:  rbac.UserType(tmp.UserType),
		UserID:    id,
		Roles:     tmp.Roles,
	}
	return nil
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
