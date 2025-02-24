package users

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

	userID string
}

func (urg *UserRolesGetter) WithUserID(id string) *UserRolesGetter {
	urg.userID = id
	return urg
}

func (urg *UserRolesGetter) Do(ctx context.Context) ([]*models.Role, error) {
	res, err := urg.connection.RunREST(ctx, urg.path(), http.MethodGet, nil)
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
	return fmt.Sprintf("/authz/users/%s/roles", urg.userID)
}
