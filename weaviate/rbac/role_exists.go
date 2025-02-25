package rbac

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

type RoleExists struct {
	connection *connection.Connection
	getter     *RoleGetter
}

func (re *RoleExists) WithName(name string) *RoleExists {
	re.getter.WithName(name)
	return re
}

func (re *RoleExists) Do(ctx context.Context) (bool, error) {
	_, err := re.getter.Do(ctx)
	return err == nil, err
}
