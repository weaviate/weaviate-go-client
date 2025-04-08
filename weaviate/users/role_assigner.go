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

type RoleAssigner struct {
	connection *connection.Connection

	userID   string
	roles    []string
	userType UserTypeInput
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
	payload := authz.AssignRoleToUserBody{
		Roles:    ra.roles,
		UserType: models.UserTypeInput(ra.userType),
	}
	res, err := ra.connection.RunREST(ctx, ra.path(), http.MethodPost, payload)
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
