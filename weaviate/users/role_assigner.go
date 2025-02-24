package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/client/authz"
)

type RoleAssigner struct {
	connection *connection.Connection

	userID string
	roles  []string
}

func (ra *RoleAssigner) WithUserID(id string) *RoleAssigner {
	ra.userID = id
	return ra
}

func (ra *RoleAssigner) WithRoles(roles ...string) *RoleAssigner {
	ra.roles = roles
	return ra
}

func (ra *RoleAssigner) Do(ctx context.Context) error {
	res, err := ra.connection.RunREST(ctx, ra.path(), http.MethodPost, authz.AssignRoleToUserBody{
		Roles: ra.roles,
	})
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (ra *RoleAssigner) path() string {
	return fmt.Sprintf("/authz/users/%s/assign", ra.userID)
}
