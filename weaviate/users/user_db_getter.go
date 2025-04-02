package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserDbGetter struct {
	connection *connection.Connection

	userID string
}

func (r *UserDbGetter) WithUserID(id string) *UserDbGetter {
	r.userID = id
	return r
}

func (r *UserDbGetter) Do(ctx context.Context) (UserInfo, error) {
	payload := users.NewGetUserInfoParams().WithUserID(r.userID)

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodGet, payload)
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

func (r *UserDbGetter) path() string {
	return fmt.Sprintf("/users/db/%s", r.userID)
}
