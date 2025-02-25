package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// TenantsGetter builder object to get class tenants
type TenantsGetter struct {
	connection *connection.Connection
	className  string
}

// WithClassName specifies the class tenants will be fetched from
func (tg *TenantsGetter) WithClassName(className string) *TenantsGetter {
	tg.className = className
	return tg
}

// Do gets tenants of given class
func (tg *TenantsGetter) Do(ctx context.Context) ([]models.Tenant, error) {
	responseData, err := tg.connection.RunREST(ctx, fmt.Sprintf("/schema/%s/tenants", tg.className), http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if responseData.StatusCode == 200 {
		var tenants []models.Tenant
		decodeErr := responseData.DecodeBodyIntoTarget(&tenants)
		return tenants, decodeErr
	}
	return nil, except.NewWeaviateClientError(responseData.StatusCode, string(responseData.Body))
}
