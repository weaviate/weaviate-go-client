package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/authz"
)

type PermissionRemover struct {
	connection *connection.Connection

	role Role
}

func (pc *PermissionRemover) WithRole(role string) *PermissionRemover {
	pc.role.Name = role
	return pc
}

func (pc *PermissionRemover) WithPermissions(permissions ...Permission) *PermissionRemover {
	for _, perm := range permissions {
		perm.ExtendRole(&pc.role)
	}
	return pc
}

func (pr *PermissionRemover) Do(ctx context.Context) error {
	res, err := pr.connection.RunREST(ctx, pr.path(), http.MethodPost, authz.RemovePermissionsBody{
		Permissions: pr.role.makeWeaviatePermissions(),
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
	return fmt.Sprintf("/authz/roles/%s/remove-permissions", pr.role.Name)
}
