package groups

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate/entities/models"
)

type GroupsOIDC struct {
	connection *connection.Connection
}

func (r *GroupsOIDC) RolesGetter() *GroupRolesGetter {
	return &GroupRolesGetter{connection: r.connection, groupType: models.GroupTypeOidc}
}

func (r *GroupsOIDC) RolesAssigner() *RoleAssigner {
	return &RoleAssigner{connection: r.connection, groupType: models.GroupTypeOidc}
}

func (r *GroupsOIDC) RolesRevoker() *RoleRevoker {
	return &RoleRevoker{connection: r.connection, groupType: models.GroupTypeOidc}
}

func (r *GroupsOIDC) GetKnownGroups() *KnownGroupLister {
	return &KnownGroupLister{connection: r.connection, groupType: models.GroupTypeOidc}
}
