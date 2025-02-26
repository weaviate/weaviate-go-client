package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/authz"
)

type PermissionAdder struct {
	connection *connection.Connection

	role Role
}

func (pa *PermissionAdder) WithRole(role string) *PermissionAdder {
	pa.role.Name = role
	return pa
}

func (pa *PermissionAdder) WithPermissions(permissions ...Permission) *PermissionAdder {
	for _, perm := range permissions {
		perm.ExtendRole(&pa.role)
	}
	return pa
}

func (pa *PermissionAdder) Do(ctx context.Context) error {
	res, err := pa.connection.RunREST(ctx, pa.path(), http.MethodPost, authz.AddPermissionsBody{
		Permissions: pa.role.makeWeaviatePermissions(),
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
	return fmt.Sprintf("/authz/roles/%s/add-permissions", pa.role.Name)
}
