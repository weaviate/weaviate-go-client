package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func TestMTBatching_stress(t *testing.T) {
	className := "Collection"
	client := testsuit.CreateTestClient(true)
	ctx := context.Background()

	defer func() {
		err := client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		require.Nil(t, err)
	}()

	t.Run("Create multi-tenancy collection", func(t *testing.T) {
		err := client.Schema().ClassCreator().WithClass(&models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					DataType: []string{"string"},
					Name:     "name",
				},
				{
					DataType: []string{"int"},
					Name:     "count",
				},
				{
					DataType: []string{"number"},
					Name:     "price",
				},
			},
			MultiTenancyConfig: &models.MultiTenancyConfig{
				AutoTenantActivation: true,
				AutoTenantCreation:   true,
				Enabled:              true,
			},
			Vectorizer: "none",
		}).Do(ctx)
		if err != nil {
			t.Fatalf("Create collection: %v", err)
		}
	})

	t.Run("Create 200 objects per 1000 tenants", func(t *testing.T) {
		tenantsN := 1000
		objectsN := 200
		for i := 0; i < tenantsN; i++ {
			batcher := client.Batch().ObjectsBatcher()
			tenant := fmt.Sprintf("tenant-%v", i)
			for j := 0; j < objectsN; j++ {
				vector := make([]float32, 1536)
				for i := range vector {
					vector[i] = float32(i) / 1536
				}
				obj := &models.Object{
					Class: "Collection",
					Properties: map[string]interface{}{
						"name":  fmt.Sprintf("object-%v", j),
						"count": j,
						"price": 42.0 * float64(j),
					},
					Vector: vector,
					Tenant: tenant,
				}
				batcher = batcher.WithObjects(obj)
			}
			_, err := batcher.Do(ctx)
			if err != nil {
				t.Fatalf("Batch create objects: %v", err)
			}
			t.Logf("Batch created %v objects in %s", objectsN, tenant)
		}
	})

	t.Run("Search for objects", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			tenant := fmt.Sprintf("tenant-%v", i)
			where := &filters.WhereBuilder{}
			res, err := client.GraphQL().Get().
				WithClassName(className).
				WithWhere(where.WithOperator(filters.LessThan).
					WithPath([]string{"price"}).WithValueNumber(float64(100))).
				WithTenant(tenant).
				WithFields(graphql.Field{Name: "name"}).
				Do(ctx)
			if err != nil {
				t.Fatalf("Search objects usage error: %v", err)
			}
			if res.Errors != nil {
				t.Fatalf("Search objects GQL error: %+v", res.Errors[0])
			}
			objects := res.Data["Get"].(map[string]interface{})[className].([]interface{})
			require.Greater(t, len(objects), 0)
		}
	})
}
