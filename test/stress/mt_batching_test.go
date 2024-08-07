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
	require.Nil(t, err)

	t.Run("Create 1000 objects per 1000 tenants", func(t *testing.T) {
		tenantsN := 1000
		objectsN := 1000
		count := 0
		batcher := client.Batch().ObjectsBatcher()
		for i := 0; i < tenantsN; i++ {
			tenant := fmt.Sprintf("tenant-%v", i)
			for j := 0; j < objectsN; j++ {
				vector := make([]float32, 16)
				for k := range vector {
					vector[k] = float32(i+j+k) / float32(16+i+j)
				}
				obj := &models.Object{
					Class: "Collection",
					Properties: map[string]interface{}{
						"name":  fmt.Sprintf("object-%v", j),
						"count": j,
						"price": 42.0 * float64(j),
					},
					Tenant: tenant,
					Vector: vector,
				}
				batcher = batcher.WithObjects(obj)

				if count == 1000 {
					_, err := batcher.Do(ctx)
					require.Nil(t, err)
					batcher = client.Batch().ObjectsBatcher()
					count = 0
				} else {
					count++
				}
			}
		}
		_, err := batcher.Do(ctx)
		require.Nil(t, err)
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
			require.Nil(t, err)
			if res.Errors != nil {
				t.Fatalf("Search objects GQL error: %+v", res.Errors[0])
			}
			objects := res.Data["Get"].(map[string]interface{})[className].([]interface{})
			require.Greater(t, len(objects), 0)
		}
	})
}
