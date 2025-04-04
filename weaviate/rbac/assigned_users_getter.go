package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type AssignedUsersGetter struct {
	connection *connection.Connection

	role string
}

func (aug *AssignedUsersGetter) WithRole(role string) *AssignedUsersGetter {
	aug.role = role
	return aug
}

func (aug *AssignedUsersGetter) Do(ctx context.Context) ([]UserAssignment, error) {
	res, err := aug.connection.RunREST(ctx, aug.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var users []struct {
			UserID   string                `json:"userId"`
			UserType models.UserTypeOutput `json:"userType"`
		}
		decodeErr := res.DecodeBodyIntoTarget(&users)
		res := make([]UserAssignment, len(users))
		for i, user := range users {
			res[i] = UserAssignment{
				UserID:   user.UserID,
				UserType: mapUserType(user.UserType),
			}
		}
		return res, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (aug *AssignedUsersGetter) path() string {
	return fmt.Sprintf("/authz/roles/%s/user-assignments", aug.role)
}
