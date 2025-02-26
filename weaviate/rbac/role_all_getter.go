package rbac

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleAllGetter struct {
	connection *connection.Connection
}

func (rag *RoleAllGetter) Do(ctx context.Context) ([]Role, error) {
	res, err := rag.connection.RunREST(ctx, "/authz/roles", http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var roles []*models.Role
		decodeErr := res.DecodeBodyIntoTarget(&roles)

		var out []Role
		for _, role := range roles {
			out = append(out, roleFromWeaviate(role))
		}
		return out, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
