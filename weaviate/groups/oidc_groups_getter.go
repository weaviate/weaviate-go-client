package groups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type KnownGroupList []string

type KnownGroupLister struct {
	connection *connection.Connection
	groupType  models.GroupType
}

func (r *KnownGroupLister) Do(ctx context.Context) (KnownGroupList, error) {
	res, err := r.connection.RunREST(ctx, r.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var response []string
		err := res.DecodeBodyIntoTarget(&response)
		if err != nil {
			return nil, except.NewDerivedWeaviateClientError(err)
		}

		return response, err
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *KnownGroupLister) path() string {
	return fmt.Sprintf("/authz/groups/%s", string(r.groupType))
}
