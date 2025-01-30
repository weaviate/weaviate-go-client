package rbac

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleGetter struct {
	connection *connection.Connection

	name string
}

func (rg *RoleGetter) WithName(name string) *RoleGetter {
	rg.name = name
	return rg
}

func (rg *RoleGetter) Do(ctx context.Context) (*models.Role, error) {
	res, err := rg.connection.RunREST(ctx, "/authz/roles/"+rg.name, http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var role models.Role
		decodeErr := res.DecodeBodyIntoTarget(&role)
		return &role, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
