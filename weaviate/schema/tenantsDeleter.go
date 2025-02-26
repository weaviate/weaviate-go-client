package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// TenantsDeleter builder object to delete tenants
type TenantsDeleter struct {
	connection *connection.Connection
	className  string
	tenants    []string
}

// WithClassName specifies the class that tenants will be deleted from
func (td *TenantsDeleter) WithClassName(className string) *TenantsDeleter {
	td.className = className
	return td
}

// WithTenants specifies tenants that will be deleted from the class
func (td *TenantsDeleter) WithTenants(tenants ...string) *TenantsDeleter {
	td.tenants = tenants
	return td
}

// Deletes tenants from the class specified in the builder
func (td *TenantsDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/tenants", td.className)
	responseData, err := td.connection.RunREST(ctx, path, http.MethodDelete, td.tenants)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
