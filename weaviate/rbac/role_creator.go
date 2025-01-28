package rbac

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleCreator struct {
	connection *connection.Connection

	name        string
	permissions []*models.Permission
}

func (rc *RoleCreator) WithName(name string) *RoleCreator {
	rc.name = name
	return rc
}

func (rc *RoleCreator) WithPermissions(permissions ...*models.Permission) *RoleCreator {
	rc.permissions = append([]*models.Permission(nil), permissions...)
	return rc
}

func (rc *RoleCreator) Do(ctx context.Context) error {
	res, err := rc.connection.RunREST(ctx, "/authz/roles", http.MethodPost, models.Role{
		Name:        &rc.name,
		Permissions: rc.permissions,
	})
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}

	if res.StatusCode == http.StatusCreated {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
