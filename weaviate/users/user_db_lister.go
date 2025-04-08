package users

import (
	"context"
	"net/http"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate/entities/models"
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
	res, err := r.connection.RunREST(ctx, r.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var response []*models.DBUserInfo
		err := res.DecodeBodyIntoTarget(&response)
		data := make(UserInfoList, len(response))
		for i, user := range response {
			data[i] = UserInfo{
				Active:    *user.Active,
				CreatedAt: time.Time(user.CreatedAt),
				UserType:  rbac.UserType(*user.DbUserType),
				UserID:    *user.UserID,
				Roles:     []*rbac.Role{},
			}
			for _, role := range user.Roles {
				data[i].Roles = append(data[i].Roles, &rbac.Role{
					Name: role,
				})
			}
		}
		return data, err
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (r *UserDBLister) path() string {
	return "/users/db"
}
