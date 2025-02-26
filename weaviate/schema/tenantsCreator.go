package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// TenantsCreator builder object to create tenants
type TenantsCreator struct {
	connection *connection.Connection
	className  string
	tenants    []models.Tenant
}

// WithClassName specifies the class that tenants will be added to
func (tc *TenantsCreator) WithClassName(className string) *TenantsCreator {
	tc.className = className
	return tc
}

// WithTenants specifies tenants that will be added to the class
func (tc *TenantsCreator) WithTenants(tenants ...models.Tenant) *TenantsCreator {
	tc.tenants = tenants
	return tc
}

// Add tenants to the class specified in the builder
func (tc *TenantsCreator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/tenants", tc.className)
	responseData, err := tc.connection.RunREST(ctx, path, http.MethodPost, tc.tenants)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
