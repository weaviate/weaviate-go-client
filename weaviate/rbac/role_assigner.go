package rbac

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

	user  string
	roles []string
}

func (ra *RoleAssigner) WithName(name string) *RoleAssigner {
	ra.user = name
	return ra
}

func (ra *RoleAssigner) WithRoles(roles ...string) *RoleAssigner {
	ra.roles = append([]string(nil), roles...)
	return ra
}

func (ra *RoleAssigner) Do(ctx context.Context) error {
	res, err := ra.connection.RunREST(ctx, ra.path(), http.MethodPost, authz.AssignRoleBody{
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
	return fmt.Sprintf("/authz/roles/%s/assign", ra.user)
}
