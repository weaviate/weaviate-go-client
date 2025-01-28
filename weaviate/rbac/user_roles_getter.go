package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type UserRolesGetter struct {
	connection *connection.Connection

	user string
}

func (urg *UserRolesGetter) WithUser(user string) *UserRolesGetter {
	urg.user = user
	return urg
}

func (urg *UserRolesGetter) Do(ctx context.Context) ([]*models.Role, error) {
	path := "/authz/users/own-roles"
	if urg.user != "" {
		path = urg.path()
	}
	res, err := urg.connection.RunREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var roles []*models.Role
		decodeErr := res.DecodeBodyIntoTarget(&roles)
		return roles, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (urg *UserRolesGetter) path() string {
	return fmt.Sprintf("/authz/users/%s/roles", urg.user)
}
