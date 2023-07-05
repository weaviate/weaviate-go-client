package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// TenantCreator builder object to create tenants
type TenantCreator struct {
	connection *connection.Connection
	className  string
	tenants    []models.Tenant
}

// WithClassName specifies the class that tenants will be added to
func (tc *TenantCreator) WithClassName(className string) *TenantCreator {
	tc.className = className
	return tc
}

// WithTenants specifies tenants that will be added to the class
func (tc *TenantCreator) WithTenants(tenants ...models.Tenant) *TenantCreator {
	tc.tenants = tenants
	return tc
}

// Add tenants to the class specified in the builder
func (tc *TenantCreator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/tenants", tc.className)
	responseData, err := tc.connection.RunREST(ctx, path, http.MethodPost, tc.tenants)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
