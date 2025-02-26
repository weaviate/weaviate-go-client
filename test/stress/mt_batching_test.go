package test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func TestMTBatching_stress(t *testing.T) {
	className := "MtStressTestCollection"
	client := testsuit.CreateTestClient(true)
	ctx := context.Background()

	defer func() {
		err := client.Schema().ClassDeleter().WithClassName(className).Do(ctx)
		require.Nil(t, err)
	}()
	require.Nil(t, client.Schema().ClassDeleter().WithClassName(className).Do(ctx))
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

	tenantsN := 100
	objectsN := 25

	wg := sync.WaitGroup{}
	eg := errgroup.Group{}
	eg.SetLimit(20)
	for i := 0; i < objectsN; i++ {
		i := i
		wg.Add(1)
		eg.Go(func() error {
			batcher := client.Batch().ObjectsBatcher()

			for j := 0; j < tenantsN; j++ {
				tenant := fmt.Sprintf("tenant-%v", j)
				vector := make([]float32, 16)
				for k := range vector {
					vector[k] = float32(i+j+k) / float32(16+i+j)
				}
				obj := &models.Object{
					Class: className,
					Properties: map[string]interface{}{
						"name":  fmt.Sprintf("object-%v", i),
						"count": i,
						"price": 42.0 * float64(i),
					},
					Tenant: tenant,
					Vector: vector,
				}
				batcher = batcher.WithObjects(obj)
			}
			resp, err := batcher.Do(ctx)
			require.Nil(t, err)
			require.NotNil(t, resp)
			for _, r := range resp {
				require.Nil(t, r.Result.Errors)
				require.Equal(t, *r.Result.Status, "SUCCESS")
			}
			wg.Done()
			return nil
		})
	}
	wg.Wait()

	tenants, err := client.Schema().TenantsGetter().WithClassName(className).Do(ctx)
	require.Nil(t, err)
	require.Len(t, tenants, tenantsN)
	for i := 0; i < tenantsN; i++ {
		tenant := fmt.Sprintf("tenant-%v", i)
		res, err := client.GraphQL().Get().
			WithClassName(className).
			WithTenant(tenant).
			WithFields(graphql.Field{Name: "name"}).
			Do(ctx)
		require.Nil(t, err)
		require.Nil(t, res.Errors)

		objects := res.Data["Get"].(map[string]interface{})[className].([]interface{})
		require.Equal(t, len(objects), objectsN)
	}
}
