package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// TenantsUpdater builder object to update tenants
type TenantsUpdater struct {
	connection *connection.Connection
	className  string
	tenants    []models.Tenant
}

// WithClassName specifies the class that tenants of will be updated
func (tu *TenantsUpdater) WithClassName(className string) *TenantsUpdater {
	tu.className = className
	return tu
}

// WithTenants specifies tenants of the class that will be updated
func (tu *TenantsUpdater) WithTenants(tenants ...models.Tenant) *TenantsUpdater {
	tu.tenants = tenants
	return tu
}

// Update tenants of the class specified in the builder
func (tu *TenantsUpdater) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/tenants", tu.className)
	responseData, err := tu.connection.RunREST(ctx, path, http.MethodPut, tu.tenants)
	return except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
}
