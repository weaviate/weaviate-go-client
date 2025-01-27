package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type PermissionChecker struct {
	connection *connection.Connection

	role       string
	permission *models.Permission
}

func (pc *PermissionChecker) WithRole(role string) *PermissionChecker {
	pc.role = role
	return pc
}

func (pc *PermissionChecker) WithPermissions(permission *models.Permission) *PermissionChecker {
	pc.permission = permission
	return pc
}

func (pc *PermissionChecker) Do(ctx context.Context) (bool, error) {
	res, err := pc.connection.RunREST(ctx, pc.path(), http.MethodPost, pc.permission)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var hasPerm bool
		decodeErr := res.DecodeBodyIntoTarget(&hasPerm)
		return hasPerm, decodeErr
	}
	return false, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (pc *PermissionChecker) path() string {
	return fmt.Sprintf("/authz/roles/%s/has-permission", pc.role)
}
