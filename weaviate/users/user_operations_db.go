package users

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

type UserOperationsDB struct {
	connection *connection.Connection
}

func (r *UserOperationsDB) RolesGetter() *UserRolesGetter {
	return &UserRolesGetter{connection: r.connection, userType: UserTypeInputDB}
}

func (r *UserOperationsDB) RolesAssigner() *RoleAssigner {
	return &RoleAssigner{connection: r.connection, userType: UserTypeInputDB}
}

func (r *UserOperationsDB) RolesRevoker() *RoleRevoker {
	return &RoleRevoker{connection: r.connection, userType: UserTypeInputDB}
}

func (r *UserOperationsDB) Creator() *UserDBCreator {
	return &UserDBCreator{connection: r.connection}
}

func (r *UserOperationsDB) Activator() *UserDBActivator {
	return &UserDBActivator{connection: r.connection}
}

func (r *UserOperationsDB) Deactivator() *UserDBDeactivator {
	return &UserDBDeactivator{connection: r.connection}
}

func (r *UserOperationsDB) Deleter() *UserDBDeleter {
	return &UserDBDeleter{connection: r.connection}
}

func (r *UserOperationsDB) Getter() *UserDBGetter {
	return &UserDBGetter{connection: r.connection}
}

func (r *UserOperationsDB) Lister() *UserDBLister {
	return &UserDBLister{connection: r.connection}
}

func (r *UserOperationsDB) KeyRotator() *UserDBKeyRotator {
	return &UserDBKeyRotator{connection: r.connection}
}
