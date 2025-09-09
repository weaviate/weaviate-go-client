package groups

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate/entities/models"
)

type GroupRolesGetter struct {
	connection *connection.Connection

	groupID          string
	groupType        models.GroupType
	includeFullRoles bool
}

func (grg *GroupRolesGetter) WithGroupID(id string) *GroupRolesGetter {
	grg.groupID = id
	return grg
}

func (grg *GroupRolesGetter) WithIncludeFullRoles(include bool) *GroupRolesGetter {
	grg.includeFullRoles = include
	return grg
}

func (grg *GroupRolesGetter) Do(ctx context.Context) ([]*rbac.Role, error) {
	res, err := grg.connection.RunREST(ctx, grg.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var roles []*rbac.Role
		decodeErr := res.DecodeBodyIntoTarget(&roles)
		return roles, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (grg *GroupRolesGetter) path() string {
	path := fmt.Sprintf("/authz/groups/%s/roles/%s", url.PathEscape(grg.groupID), grg.groupType)

	if grg.includeFullRoles {
		path += "?includeFullRoles=" + fmt.Sprint(grg.includeFullRoles)
	}

	return path
}
