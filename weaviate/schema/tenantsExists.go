package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// TenantsExists builder object to check if tenant exists
type TenantsExists struct {
	connection *connection.Connection
	className  string
	tenant     string
}

// WithClassName specifies the class tenants will be fetched from
func (te *TenantsExists) WithClassName(className string) *TenantsExists {
	te.className = className
	return te
}

// WithTenant specifies the class tenants will be fetched from
func (te *TenantsExists) WithTenant(tenant string) *TenantsExists {
	te.tenant = tenant
	return te
}

// Do head tenant of given class
func (te *TenantsExists) Do(ctx context.Context) (bool, error) {
	responseData, err := te.connection.RunREST(ctx, fmt.Sprintf("/schema/%s/tenants/%s", te.className, te.tenant), http.MethodHead, nil)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	return responseData.StatusCode == 200, nil
}
