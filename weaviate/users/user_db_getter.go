package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type UserDBGetter struct {
	connection *connection.Connection

	userID string
}

func (r *UserDBGetter) WithUserID(id string) *UserDBGetter {
	r.userID = id
	return r
}

func (r *UserDBGetter) Do(ctx context.Context) (UserInfo, error) {
	res, err := r.connection.RunREST(ctx, r.path(), http.MethodGet, nil)
	if err != nil {
		return UserInfo{}, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var data UserInfo
		err := res.DecodeBodyIntoTarget(&data)
		return data, err
	}
	return UserInfo{}, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDBGetter) path() string {
	return fmt.Sprintf("/users/db/%s", r.userID)
}
