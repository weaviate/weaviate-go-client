package rbac

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleCreator struct {
	connection *connection.Connection

	role Role
}

func (rc *RoleCreator) WithRole(role Role) *RoleCreator {
	rc.role = role
	return rc
}

func (rc *RoleCreator) Do(ctx context.Context) error {
	res, err := rc.connection.RunREST(ctx, "/authz/roles", http.MethodPost, &models.Role{
		Name:        &rc.role.Name,
		Permissions: rc.role.makeWeaviatePermissions(),
	})
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}

	if res.StatusCode == http.StatusCreated {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
