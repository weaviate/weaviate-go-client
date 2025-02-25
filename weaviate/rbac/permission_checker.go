package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
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
// Only first action in the permission's list of actions will be used.
func (pc *PermissionChecker) WithPermission(permission Permission) *PermissionChecker {
	permission.ExtendRole(&pc.role)
	return pc
}

func (pc *PermissionChecker) Do(ctx context.Context) (bool, error) {
	checkPermission := pc.role.makeWeaviatePermissions()[0]
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
