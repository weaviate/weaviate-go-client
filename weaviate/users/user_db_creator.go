package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserDbCreator struct {
	connection *connection.Connection

	userID string
}

func (r *UserDbCreator) WithUserID(id string) *UserDbCreator {
	r.userID = id
	return r
}

func (r *UserDbCreator) Do(ctx context.Context) error {
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

func (r *UserDbCreator) path() string {
	return fmt.Sprintf("/users/db/%s", r.userID)
}
