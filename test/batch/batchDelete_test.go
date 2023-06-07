package batch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/helpers"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
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
		client := testsuit.CreateTestClient()
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

	t.Run("batch delete", func(t *testing.T) {
		client := testsuit.CreateTestClient()
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
		client := testsuit.CreateTestClient()
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
		client := testsuit.CreateTestClient()
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
	idsByClass := map[string][]string{
		"Pizza": {
			"10523cdd-15a2-42f4-81fa-267fe92f7cd6",
			"927dd3ac-e012-4093-8007-7799cc7e81e4",
			"00000000-0000-0000-0000-000000000000",
			"5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
		},
		"Soup": {
			"8c156d37-81aa-4ce9-a811-621e2702b825",
			"27351361-2898-4d1a-aad7-1ca48253eb0b",
		},
		"Risotto": {
			"da751a25-f573-4715-a893-e607b2de0ba4",
			"10c2ee44-7d58-42be-9d64-5766883ca8cb",
			"696bf381-7f98-40a4-bcad-841780e00e0e",
		},
	}

	t.Run("setup weaviate", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("deletes objects from multi tenant class", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		tenants := []string{"tenantNo1", "tenantNo2"}

		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants...)

		for _, tenant := range tenants {
			for className, ids := range idsByClass {
				resp, err := client.Batch().ObjectsBatchDeleter().
					WithClassName(className).
					WithWhere(filters.Where().
						WithOperator(filters.Like).
						WithPath([]string{"name"}).
						WithValueText("*")).
					WithOutput("minimal").
					WithTenantKey(tenant).
					Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, resp)
				require.NotNil(t, resp.Results)
				assert.Equal(t, int64(len(ids)), resp.Results.Matches)
				assert.Equal(t, int64(len(ids)), resp.Results.Successful)
			}
		}

		t.Run("check deleted", func(t *testing.T) {
			for _, tenant := range tenants {
				for className, ids := range idsByClass {
					for _, id := range ids {
						exists, err := client.Data().Checker().
							WithID(id).
							WithClassName(className).
							WithTenantKey(tenant).
							Do(context.Background())

						require.Nil(t, err)
						require.False(t, exists)
					}
				}
			}
		})

		t.Run("clean up classes", func(t *testing.T) {
			client := testsuit.CreateTestClient()
			err := client.Schema().AllDeleter().Do(context.Background())
			require.Nil(t, err)
		})
	})

	t.Run("fails deleting objects from multi tenant class without tenant key", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		tenants := []string{"tenantNo1", "tenantNo2"}

		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants...)

		for className := range idsByClass {
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
			assert.Contains(t, clientErr.Msg, "has multi-tenancy enabled")
			require.Nil(t, resp)
		}

		t.Run("check not deleted", func(t *testing.T) {
			for _, tenant := range tenants {
				for className, ids := range idsByClass {
					for _, id := range ids {
						exists, err := client.Data().Checker().
							WithID(id).
							WithClassName(className).
							WithTenantKey(tenant).
							Do(context.Background())

						require.Nil(t, err)
						require.True(t, exists)
					}
				}
			}
		})

		t.Run("clean up classes", func(t *testing.T) {
			client := testsuit.CreateTestClient()
			err := client.Schema().AllDeleter().Do(context.Background())
			require.Nil(t, err)
		})
	})

	t.Run("does not delete objects from multi tenant class with different tenant key than filter", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		tenants := []string{"tenantNo1", "tenantNo2"}

		testsuit.CreateSchemaFoodForTenants(t, client)
		testsuit.CreateTenantsFood(t, client, tenants...)
		testsuit.CreateDataFoodForTenants(t, client, tenants...)

		for className := range idsByClass {
			resp, err := client.Batch().ObjectsBatchDeleter().
				WithClassName(className).
				WithWhere(filters.Where().
					WithOperator(filters.Equal).
					WithPath([]string{testsuit.TenantKey}).
					WithValueText(tenants[1])).
				WithOutput("minimal").
				WithTenantKey(tenants[0]).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Results)
			assert.Equal(t, int64(0), resp.Results.Matches)
			assert.Equal(t, int64(0), resp.Results.Successful)
		}

		t.Run("check not deleted", func(t *testing.T) {
			for _, tenant := range tenants {
				for className, ids := range idsByClass {
					for _, id := range ids {
						exists, err := client.Data().Checker().
							WithID(id).
							WithClassName(className).
							WithTenantKey(tenant).
							Do(context.Background())

						require.Nil(t, err)
						require.True(t, exists)
					}
				}
			}
		})

		t.Run("clean up classes", func(t *testing.T) {
			client := testsuit.CreateTestClient()
			err := client.Schema().AllDeleter().Do(context.Background())
			require.Nil(t, err)
		})
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
