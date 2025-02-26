package rbac

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type RoleGetter struct {
	connection *connection.Connection

	name string
}

func (rg *RoleGetter) WithName(name string) *RoleGetter {
	rg.name = name
	return rg
}

func (rg *RoleGetter) Do(ctx context.Context) (Role, error) {
	res, err := rg.connection.RunREST(ctx, "/authz/roles/"+rg.name, http.MethodGet, nil)
	if err != nil {
		return Role{}, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var role Role
		decodeErr := res.DecodeBodyIntoTarget(&role)
		return role, decodeErr
	}
	return Role{}, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
