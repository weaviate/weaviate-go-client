package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type UserDBDeactivator struct {
	connection *connection.Connection

	userID    string
	revokeKey bool
}

func (r *UserDBDeactivator) WithUserID(id string) *UserDBDeactivator {
	r.userID = id
	return r
}

func (r *UserDBDeactivator) WithRevokeKey(revokeKey bool) *UserDBDeactivator {
	r.revokeKey = revokeKey
	return r
}

func (r *UserDBDeactivator) Do(ctx context.Context) (bool, error) {
	payload := struct {
		RevokeKey bool `json:"revoke_key"`
	}{
		RevokeKey: r.revokeKey,
	}

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodPost, payload)
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

func (r *UserDBDeactivator) path() string {
	return fmt.Sprintf("/users/db/%s/deactivate", r.userID)
}
