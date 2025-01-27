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

func (rc *RoleGetter) WithName(name string) *RoleGetter {
	rc.name = name
	return rc
}

func (rc *RoleGetter) Do(ctx context.Context) (*models.Role, error) {
	res, err := rc.connection.RunREST(ctx, "/authz/roles/"+rc.name, http.MethodPost, nil)
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
