package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/rbac"
)

type MyUserGetter struct {
	connection *connection.Connection
}

type UserInfo struct {
	UserID string
	Roles  []*rbac.Role
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

func (info *UserInfo) UnmarshalJSON(data []byte) error {
	var tmp struct {
		UserID   string       `json:"user_id"`
		Username string       `json:"username"`
		Roles    []*rbac.Role `json:"roles"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	id := tmp.UserID
	if id == "" {
		id = tmp.Username
	}
	*info = UserInfo{
		UserID: id,
		Roles:  tmp.Roles,
	}
	return nil
}
