package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
	"github.com/weaviate/weaviate/entities/models"
)

type UserDbKeyRotator struct {
	connection *connection.Connection

	userID string
}

func (r *UserDbKeyRotator) WithUserID(id string) *UserDbKeyRotator {
	r.userID = id
	return r
}

func (r *UserDbKeyRotator) Do(ctx context.Context) (string, error) {
	payload := users.NewCreateUserParams().WithUserID(r.userID)

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, payload)
	if err != nil {
		return "", except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var data models.UserAPIKey
		err := res.DecodeBodyIntoTarget(&data)
		return *data.Apikey, err
	}
	return "", except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDbKeyRotator) path() string {
	return fmt.Sprintf("/users/db/%s/rotate-key", r.userID)
}
