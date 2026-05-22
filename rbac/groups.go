package rbac

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type GroupsClient struct {
	*kindClient[api.GroupID]
}

func NewGroupsClient(t internal.Transport) *GroupsClient {
	dev.AssertNotNil(t, "transport")
	return &GroupsClient{
		kindClient: newKindClient[api.GroupID](t, api.RBACKindOIDC),
	}
}
