package rbac

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

	user  string
	roles []string
}

func (rr *RoleRevoker) WithUser(user string) *RoleRevoker {
	rr.user = user
	return rr
}

func (rr *RoleRevoker) WithRoles(roles ...string) *RoleRevoker {
	rr.roles = append([]string(nil), roles...)
	return rr
}

func (rr *RoleRevoker) Do(ctx context.Context) error {
	res, err := rr.connection.RunREST(ctx, rr.path(), http.MethodPost, authz.RevokeRoleBody{
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
	return fmt.Sprintf("/authz/users/%s/revoke", rr.user)
}
