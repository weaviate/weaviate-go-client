package users

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserInfoList []UserInfo

type UserDbLister struct {
	connection *connection.Connection

	userID string
}

func (r *UserDbLister) WithUserID(id string) *UserDbLister {
	r.userID = id
	return r
}

func (r *UserDbLister) Do(ctx context.Context) (UserInfoList, error) {
	payload := users.NewListAllUsersParams()

	res, err := r.connection.RunREST(ctx, r.path(), http.MethodGet, payload)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var data UserInfoList
		err := res.DecodeBodyIntoTarget(&data)
		return data, err
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDbLister) path() string {
	return "/users/db"
}
