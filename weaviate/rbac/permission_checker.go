package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
)

type PermissionChecker struct {
	connection *connection.Connection

	role Role
}

func (pc *PermissionChecker) WithRole(role string) *PermissionChecker {
	pc.role.Name = role
	return pc
}

// WithPermission specifies the permission (singular) to be checked.
// Note that, while this method accepts a permission group for convenience,
// it will only check against the first permission in the group.
func (pc *PermissionChecker) WithPermission(permission PermissionGroup) *PermissionChecker {
	permission.ExtendRole(&pc.role)
	return pc
}

func (pc *PermissionChecker) Do(ctx context.Context) (bool, error) {
	checkPermission := pc.role.Permissions.toWeaviate()[0]
	res, err := pc.connection.RunREST(ctx, pc.path(), http.MethodPost, checkPermission)
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
	return fmt.Sprintf("/authz/roles/%s/has-permission", pc.role.Name)
}
