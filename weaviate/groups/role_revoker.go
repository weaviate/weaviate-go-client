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

type RoleRevoker struct {
	connection *connection.Connection

	groupID   string
	roles     []string
	groupType models.GroupType
}

func (ra *RoleRevoker) WithGroupId(id string) *RoleRevoker {
	ra.groupID = id
	return ra
}

func (ra *RoleRevoker) WithRoles(roles ...string) *RoleRevoker {
	ra.roles = roles
	return ra
}

func (ra *RoleRevoker) Do(ctx context.Context) error {
	payload := authz.RevokeRoleFromGroupBody{
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

func (ra *RoleRevoker) path() string {
	return fmt.Sprintf("/authz/groups/%s/revoke", url.PathEscape(ra.groupID))
}
