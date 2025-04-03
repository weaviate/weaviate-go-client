package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type UserDBKeyRotator struct {
	connection *connection.Connection

	userID string
}

func (r *UserDBKeyRotator) WithUserID(id string) *UserDBKeyRotator {
	r.userID = id
	return r
}

func (r *UserDBKeyRotator) Do(ctx context.Context) (string, error) {
	res, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, nil)
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

func (r *UserDBKeyRotator) path() string {
	return fmt.Sprintf("/users/db/%s/rotate-key", r.userID)
}
