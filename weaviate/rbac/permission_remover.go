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

type PermissionRemover struct {
	connection *connection.Connection

	role        string
	permissions []*models.Permission
}

func (pr *PermissionRemover) WithRole(role string) *PermissionRemover {
	pr.role = role
	return pr
}

func (pr *PermissionRemover) WithPermissions(permissions ...*models.Permission) *PermissionRemover {
	pr.permissions = append([]*models.Permission(nil), permissions...)
	return pr
}

func (pr *PermissionRemover) Do(ctx context.Context) error {
	res, err := pr.connection.RunREST(ctx, pr.path(), http.MethodPost, authz.RemovePermissionsBody{
		Permissions: pr.permissions,
	})
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (pr *PermissionRemover) path() string {
	return fmt.Sprintf("/authz/roles/%s/remove-permissions", pr.role)
}
