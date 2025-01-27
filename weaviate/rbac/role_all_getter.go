package rbac

import (
	"context"
	"log"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type RoleAllGetter struct {
	connection *connection.Connection
}

func (rag *RoleAllGetter) Do(ctx context.Context) ([]*models.Role, error) {
	res, err := rag.connection.RunREST(ctx, "/authz/roles", http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	log.Print("status code: ", res.StatusCode)
	if res.StatusCode == http.StatusOK {
		log.Print(string(res.Body))
		var roles []*models.Role
		decodeErr := res.DecodeBodyIntoTarget(&roles)
		log.Print("decoded successfully: ", decodeErr == nil)
		return roles, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
