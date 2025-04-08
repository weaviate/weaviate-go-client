package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type UserDBActivator struct {
	connection *connection.Connection

	userID string
}

func (r *UserDBActivator) WithUserID(id string) *UserDBActivator {
	r.userID = id
	return r
}

func (r *UserDBActivator) Do(ctx context.Context) (bool, error) {
	res, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, nil)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusConflict:
		fallthrough
	case http.StatusNotFound:
		return false, nil
	}
	return false, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDBActivator) path() string {
	return fmt.Sprintf("/users/db/%s/activate", r.userID)
}
