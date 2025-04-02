package users

import "github.com/weaviate/weaviate-go-client/v5/weaviate/connection"

type UserOperationsDB struct {
	connection *connection.Connection
}

func (r *UserOperationsDB) RolesGetter() *UserRolesGetter {
	return (&UserRolesGetter{connection: r.connection}).WithUserType(UserTypeDb)
}

func (r *UserOperationsDB) RolesAssigner() *RoleAssigner {
	return (&RoleAssigner{connection: r.connection}).WithUserType(UserTypeDb)
}

func (r *UserOperationsDB) RolesRevoker() *RoleRevoker {
	return (&RoleRevoker{connection: r.connection}).WithUserType(UserTypeDb)
}

func (r *UserOperationsDB) Creator() *UserDbCreator {
	return &UserDbCreator{connection: r.connection}
}

func (r *UserOperationsDB) Activator() *UserDbActivator {
	return &UserDbActivator{connection: r.connection}
}

func (r *UserOperationsDB) Deactivator() *UserDbDeactivator {
	return &UserDbDeactivator{connection: r.connection}
}

func (r *UserOperationsDB) Deleter() *UserDbDeleter {
	return &UserDbDeleter{connection: r.connection}
}

func (r *UserOperationsDB) Getter() *UserDbGetter {
	return &UserDbGetter{connection: r.connection}
}

func (r *UserOperationsDB) Lister() *UserDbLister {
	return &UserDbLister{connection: r.connection}
}
