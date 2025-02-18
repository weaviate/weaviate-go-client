package users

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type MyUserGetter struct {
	connection *connection.Connection
}

func (mug *MyUserGetter) Do(ctx context.Context) (*models.UserInfo, error) {
	path := "/users/own-info"
	res, err := mug.connection.RunREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var user *models.UserInfo
		decodeErr := res.DecodeBodyIntoTarget(&user)
		return user, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}
