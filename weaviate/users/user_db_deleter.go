package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserDbDeleter struct {
	connection *connection.Connection

	userID string
}

func (r *UserDbDeleter) WithUserID(id string) *UserDbDeleter {
	r.userID = id
	return r
}

func (r *UserDbDeleter) Do(ctx context.Context) error {
	payload := users.NewDeleteUserParams().WithUserID(r.userID)

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodDelete, payload)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusNoContent {
		return err
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDbDeleter) path() string {
	return fmt.Sprintf("/users/db/%s", r.userID)
}
