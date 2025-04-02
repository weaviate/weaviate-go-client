package users

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/client/users"
)

type UserInfoList []UserInfo

type UserDBLister struct {
	connection *connection.Connection

	userID string
}

func (r *UserDBLister) WithUserID(id string) *UserDBLister {
	r.userID = id
	return r
}

func (r *UserDBLister) Do(ctx context.Context) (UserInfoList, error) {
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

func (r *UserDBLister) path() string {
	return "/users/db"
}
