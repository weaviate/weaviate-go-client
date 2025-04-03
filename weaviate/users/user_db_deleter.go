package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserDBDeleter struct {
	connection *connection.Connection

	userID string
}

func (r *UserDBDeleter) WithUserID(id string) *UserDBDeleter {
	r.userID = id
	return r
}

func (r *UserDBDeleter) Do(ctx context.Context) (bool, error) {
	payload := users.NewDeleteUserParams().WithUserID(r.userID)

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodDelete, payload)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	switch res.StatusCode {
	case http.StatusNoContent:
		return true, nil
	case http.StatusNotFound:
		return false, except.NewExpectedStatusCodeErrorFromRESTResponse(res)
	}

	return false, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDBDeleter) path() string {
	return fmt.Sprintf("/users/db/%s", r.userID)
}
