package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/authz"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleRevoker struct {
	connection *connection.Connection

	userID   string
	roles    []string
	userType models.UserTypeInput
}

func (rr *RoleRevoker) WithUserID(id string) *RoleRevoker {
	rr.userID = id
	return rr
}

func (rr *RoleRevoker) WithRoles(roles ...string) *RoleRevoker {
	rr.roles = roles
	return rr
}

func (rr *RoleRevoker) WithUserType(userType UserType) *RoleRevoker {
	rr.userType = models.UserTypeInput(userType)
	return rr
}

func (rr *RoleRevoker) Do(ctx context.Context) error {
	payload := authz.RevokeRoleFromUserBody{
		Roles:    rr.roles,
		UserType: rr.userType,
	}
	res, err := rr.connection.RunREST(ctx, rr.path(), http.MethodPost, payload)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (rr *RoleRevoker) path() string {
	return fmt.Sprintf("/authz/users/%s/revoke", rr.userID)
}
