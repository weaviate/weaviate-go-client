package batch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/helpers"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestBatchDelete_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})

	t.Run("batch delete dry run", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		where := filters.Where().
			WithOperator(filters.Equal).
			WithPath([]string{"id"}).
			WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

		expectedResults := &models.BatchDeleteResponseResults{
			Matches:    1,
			Failed:     0,
			Successful: 0,
			Limit:      10000,
			Objects: []*models.BatchDeleteResponseResultsObjectsItems0{
				{
					ID:     "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
					Status: helpers.StringPointer("DRYRUN"),
				},
			},
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithClassName("Pizza").
			WithDryRun(true).
			WithOutput("verbose").
			WithWhere(where).
			Do(context.Background())
		require.Nil(t, err)
		require.NotNil(t, resp.Match)
		assert.Equal(t, "Pizza", resp.Match.Class)
		assert.Equal(t, where.Build(), resp.Match.Where)
		assert.Equal(t, expectedResults, resp.Results)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("batch delete dry run and no output setting", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		where := filters.Where().
			WithOperator(filters.Equal).
			WithPath([]string{"id"}).
			WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

		expectedResults := &models.BatchDeleteResponseResults{
			Matches:    1,
			Failed:     0,
			Successful: 0,
			Limit:      10000,
			Objects:    nil,
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithClassName("Pizza").
			WithDryRun(true).
			WithWhere(where).
			Do(context.Background())
		require.Nil(t, err)
		require.NotNil(t, resp.Match)
		assert.Equal(t, "Pizza", resp.Match.Class)
		assert.Equal(t, where.Build(), resp.Match.Where)
		assert.Equal(t, expectedResults, resp.Results)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("batch delete dry run and no output setting and no dry run", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		where := filters.Where().
			WithOperator(filters.Equal).
			WithPath([]string{"id"}).
			WithValueString("5b6a08ba-1d46-43aa-89cc-8b070790c6f2")

		expectedResults := &models.BatchDeleteResponseResults{
			Matches:    1,
			Failed:     0,
			Successful: 1,
			Limit:      10000,
			Objects:    nil,
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithClassName("Pizza").
			WithWhere(where).
			Do(context.Background())
		require.Nil(t, err)
		require.NotNil(t, resp.Match)
		assert.Equal(t, "Pizza", resp.Match.Class)
		assert.Equal(t, where.Build(), resp.Match.Where)
		assert.Equal(t, expectedResults, resp.Results)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("batch delete", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		nowString := fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))

		where := filters.Where().
			WithOperator(filters.LessThan).
			WithPath([]string{"_creationTimeUnix"}).
			WithValueString(nowString)

		expectedResults := &models.BatchDeleteResponseResults{
			Matches:    4,
			Failed:     0,
			Successful: 4,
			Limit:      10000,
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithClassName("Pizza").
			WithOutput("verbose").
			WithWhere(where).
			Do(context.Background())
		require.Nil(t, err)
		require.NotNil(t, resp.Match)
		assert.Equal(t, "Pizza", resp.Match.Class)
		assert.Equal(t, where.Build(), resp.Match.Where)
		assert.Equal(t, expectedResults.Matches, resp.Results.Matches)
		assert.Equal(t, expectedResults.Failed, resp.Results.Failed)
		assert.Equal(t, expectedResults.Successful, resp.Results.Successful)
		assert.Len(t, resp.Results.Objects, 4)

		for _, obj := range resp.Results.Objects {
			require.NotNil(t, obj.Status)
			require.NotNil(t, obj.Status)
			assert.Equal(t, helpers.StringPointer("SUCCESS"), obj.Status)
			assert.Nil(t, obj.Errors)
		}

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("batch delete with consistency level", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		nowString := fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))

		where := filters.Where().
			WithOperator(filters.LessThan).
			WithPath([]string{"_creationTimeUnix"}).
			WithValueString(nowString)

		expectedResults := &models.BatchDeleteResponseResults{
			Matches:    4,
			Failed:     0,
			Successful: 4,
			Limit:      10000,
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithClassName("Pizza").
			WithOutput("verbose").
			WithWhere(where).
			WithConsistencyLevel(replication.ConsistencyLevel.ONE).
			Do(context.Background())
		require.Nil(t, err)
		require.NotNil(t, resp.Match)
		assert.Equal(t, "Pizza", resp.Match.Class)
		assert.Equal(t, where.Build(), resp.Match.Where)
		assert.Equal(t, expectedResults.Matches, resp.Results.Matches)
		assert.Equal(t, expectedResults.Failed, resp.Results.Failed)
		assert.Equal(t, expectedResults.Successful, resp.Results.Successful)
		assert.Len(t, resp.Results.Objects, 4)

		for _, obj := range resp.Results.Objects {
			require.NotNil(t, obj.Status)
			require.NotNil(t, obj.Status)
			assert.Equal(t, helpers.StringPointer("SUCCESS"), obj.Status)
			assert.Nil(t, obj.Errors)
		}

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("batch delete no matches", func(t *testing.T) {
		client := testsuit.CreateTestClient(false)
		testsuit.CreateTestSchemaAndData(t, client)

		where := filters.Where().
			WithOperator(filters.Equal).
			WithPath([]string{"id"}).
			WithValueString("267f5125-c9fd-4ca6-9134-f383ff5f0cb6")

		expectedResults := &models.BatchDeleteResponseResults{
			Matches:    0,
			Failed:     0,
			Successful: 0,
			Limit:      10000,
			Objects:    nil,
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithClassName("Pizza").
			WithDryRun(true).
			WithOutput("verbose").
			WithWhere(where).
			Do(context.Background())
		assert.Nil(t, err)
		assert.Equal(t, expectedResults, resp.Results)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatal(err.Error())
		}
	})
}

func TestBatchDelete_MultiTenancy(t *testing.T) {
	cleanup := func() {
		client := testsuit.CreateTestClient(false)
		err := client.Schema().AllDeleter().Do(context.Background())
		require.Nil(t, err)
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("deletes objects from MT class", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		for _, tenant := range tenants {
			for className, ids := range testsuit.IdsByClass {
				resp, err := client.Batch().ObjectsBatchDeleter().
					WithClassName(className).
					WithWhere(filters.Where().
						WithOperator(filters.Like).
						WithPath([]string{"name"}).
						WithValueText("*")).
					WithOutput("minimal").
					WithTenant(tenant.Name).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.Results)
				assert.Equal(t, int64(len(ids)), resp.Results.Matches)
				assert.Equal(t, int64(len(ids)), resp.Results.Successful)
			}
		}

		t.Run("verify deleted", func(t *testing.T) {
			for _, tenant := range tenants {
				for className, ids := range testsuit.IdsByClass {
					for _, id := range ids {
						exists, err := client.Data().Checker().
							WithID(id).
							WithClassName(className).
							WithTenant(tenant.Name).
							Do(context.Background())

						require.Nil(t, err)
						require.False(t, exists)
					}
				}
			}
		})
	})

	t.Run("fails deleting objects from MT class without tenant", func(t *testing.T) {
		defer cleanup()

		tenants := testsuit.Tenants{
			{Name: "tenantNo1"},
			{Name: "tenantNo2"},
		}

		client := testsuit.CreateTestClient(false)
		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants.Names()...)

		for className := range testsuit.IdsByClass {
			resp, err := client.Batch().ObjectsBatchDeleter().
				WithClassName(className).
				WithWhere(filters.Where().
					WithOperator(filters.Like).
					WithPath([]string{"name"}).
					WithValueText("*")).
				WithOutput("minimal").
				Do(context.Background())

			require.NotNil(t, err)
			clientErr := err.(*fault.WeaviateClientError)
			assert.Equal(t, 422, clientErr.StatusCode)
			assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled, but request was without tenant")
			require.Nil(t, resp)
		}

		t.Run("verify not deleted", func(t *testing.T) {
			for _, tenant := range tenants {
				for className, ids := range testsuit.IdsByClass {
					for _, id := range ids {
						exists, err := client.Data().Checker().
							WithID(id).
							WithClassName(className).
							WithTenant(tenant.Name).
							Do(context.Background())

						require.Nil(t, err)
						require.True(t, exists)
					}
				}
			}
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
