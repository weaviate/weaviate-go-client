package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/client/authz"
)

type RoleRevoker struct {
	connection *connection.Connection

	userID string
	roles  []string
}

func (rr *RoleRevoker) WithUserID(id string) *RoleRevoker {
	rr.userID = id
	return rr
}

func (rr *RoleRevoker) WithRoles(roles ...string) *RoleRevoker {
	rr.roles = roles
	return rr
}

func (rr *RoleRevoker) Do(ctx context.Context) error {
	res, err := rr.connection.RunREST(ctx, rr.path(), http.MethodPost, authz.RevokeRoleFromUserBody{
		Roles: rr.roles,
	})
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
