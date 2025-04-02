package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserDBDeactivator struct {
	connection *connection.Connection

	userID string
}

func (r *UserDBDeactivator) WithUserID(id string) *UserDBDeactivator {
	r.userID = id
	return r
}

func (r *UserDBDeactivator) Do(ctx context.Context) error {
	payload := users.NewDeactivateUserParams().WithUserID(r.userID)

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, payload)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDBDeactivator) path() string {
	return fmt.Sprintf("/users/db/%s/deactivate", r.userID)
}
