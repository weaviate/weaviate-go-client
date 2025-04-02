package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
)

type UserRolesGetter struct {
	connection *connection.Connection

	userID           string
	userType         string
	IncludeFullRoles *bool
}

func (urg *UserRolesGetter) WithUserID(id string) *UserRolesGetter {
	urg.userID = id
	return urg
}

func (urg *UserRolesGetter) WithUserType(userType UserType) *UserRolesGetter {
	urg.userType = string(userType)
	return urg
}

func (urg *UserRolesGetter) WithIncludeFullRoles(include bool) *UserRolesGetter {
	urg.IncludeFullRoles = &include
	return urg
}

func (urg *UserRolesGetter) Do(ctx context.Context) ([]*rbac.Role, error) {
	// Assume DB user if no user type is specified
	if urg.userType == "" {
		urg = urg.WithUserType(UserTypeDB)
	}
	res, err := urg.connection.RunREST(ctx, urg.path(), http.MethodGet, nil)
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

func (urg *UserRolesGetter) path() string {
	path := fmt.Sprintf("/authz/users/%s/roles/%s", urg.userID, urg.userType)
	if urg.IncludeFullRoles != nil {
		path += "?includeFullRoles=" + fmt.Sprint(*urg.IncludeFullRoles)
	}
	return path
}
