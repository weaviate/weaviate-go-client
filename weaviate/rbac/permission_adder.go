package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/client/authz"
	"github.com/weaviate/weaviate/entities/models"
)

type PermissionAdder struct {
	connection *connection.Connection

	role        string
	permissions []*models.Permission
}

func (pa *PermissionAdder) WithRole(role string) *PermissionAdder {
	pa.role = role
	return pa
}

func (pa *PermissionAdder) WithPermissions(permissions ...*models.Permission) *PermissionAdder {
	pa.permissions = append([]*models.Permission(nil), permissions...)
	return pa
}

func (pa *PermissionAdder) Do(ctx context.Context) error {
	res, err := pa.connection.RunREST(ctx, pa.path(), http.MethodPost, authz.AddPermissionsBody{
		Permissions: pa.permissions,
	})
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (pa *PermissionAdder) path() string {
	return fmt.Sprintf("/authz/roles/%s/add-permissions", pa.role)
}
