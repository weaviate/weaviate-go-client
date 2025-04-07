package users

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

type MyUserGetter struct {
	connection *connection.Connection
}

func (mug *MyUserGetter) Do(ctx context.Context) (UserInfo, error) {
	path := "/users/own-info"
	res, err := mug.connection.RunREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		return UserInfo{}, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var user UserInfo
		decodeErr := res.DecodeBodyIntoTarget(&user)
		return user, decodeErr
	}
	return UserInfo{}, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
