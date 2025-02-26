package rbac

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type RoleDeleter struct {
	connection *connection.Connection

	name string
}

func (rc *RoleDeleter) WithName(name string) *RoleDeleter {
	rc.name = name
	return rc
}

func (rc *RoleDeleter) Do(ctx context.Context) error {
	res, err := rc.connection.RunREST(ctx, "/authz/roles/"+rc.name, http.MethodDelete, nil)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusNoContent { // 204 - Successfully deleted
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
