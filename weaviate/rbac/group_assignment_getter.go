package rbac

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

type GroupAssignmentGetter struct {
	connection *connection.Connection

	role string
}

func (aug *GroupAssignmentGetter) WithRole(role string) *GroupAssignmentGetter {
	aug.role = role
	return aug
}

func (aug *GroupAssignmentGetter) Do(ctx context.Context) ([]GroupAssignment, error) {
	res, err := aug.connection.RunREST(ctx, aug.path(), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if res.StatusCode == http.StatusOK {
		var groups []struct {
			GroupID   string           `json:"groupId"`
			GroupType models.GroupType `json:"groupType"`
		}
		decodeErr := res.DecodeBodyIntoTarget(&groups)
		res := make([]GroupAssignment, len(groups))
		for i, group := range groups {
			res[i] = GroupAssignment{
				Group:     group.GroupID,
				GroupType: GroupType(group.GroupType),
			}
		}
		return res, decodeErr
	}
	return nil, except.NewUnexpectedStatusCodeErrorFromRESTResponse(res)
}

func (aug *GroupAssignmentGetter) path() string {
	return fmt.Sprintf("/authz/roles/%s/group-assignments", aug.role)
}
