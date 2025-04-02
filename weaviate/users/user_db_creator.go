package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserDBCreator struct {
	connection *connection.Connection

	userID string
}

func (r *UserDBCreator) WithUserID(id string) *UserDBCreator {
	r.userID = id
	return r
}

func (r *UserDBCreator) Do(ctx context.Context) error {
	payload := users.NewCreateUserParams().WithUserID(r.userID)

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, payload)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDBCreator) path() string {
	return fmt.Sprintf("/users/db/%s", r.userID)
}
