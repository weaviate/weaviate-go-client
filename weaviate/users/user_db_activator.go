package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserDbActivator struct {
	connection *connection.Connection

	userID string
}

func (r *UserDbActivator) WithUserID(id string) *UserDbActivator {
	r.userID = id
	return r
}

func (r *UserDbActivator) Do(ctx context.Context) error {
	payload := users.NewActivateUserParams().WithUserID(r.userID)

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, payload)
	if err != nil {
		return except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	return except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDbActivator) path() string {
	return fmt.Sprintf("/users/db/%s/activate", r.userID)
}
