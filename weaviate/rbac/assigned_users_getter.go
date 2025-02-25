package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type AssignedUsersGetter struct {
	connection *connection.Connection

	role string
}

func (aug *AssignedUsersGetter) WithRole(role string) *AssignedUsersGetter {
	aug.role = role
	return aug
}

func (aug *AssignedUsersGetter) Do(ctx context.Context) ([]string, error) {
	res, err := aug.connection.RunREST(ctx, aug.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var users []string
		decodeErr := res.DecodeBodyIntoTarget(&users)
		return users, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (aug *AssignedUsersGetter) path() string {
	return fmt.Sprintf("/authz/roles/%s/users", aug.role)
}
