package groups

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/authz"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleAssigner struct {
	connection *connection.Connection

	groupID   string
	roles     []string
	groupType models.GroupType
}

func (ra *RoleAssigner) WithGroupId(id string) *RoleAssigner {
	ra.groupID = id
	return ra
}

func (ra *RoleAssigner) WithRoles(roles ...string) *RoleAssigner {
	ra.roles = roles
	return ra
}

func (ra *RoleAssigner) Do(ctx context.Context) error {
	payload := authz.AssignRoleToGroupBody{
		Roles:     ra.roles,
		GroupType: ra.groupType,
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
	return fmt.Sprintf("/authz/groups/%s/assign", url.PathEscape(ra.groupID))
}
