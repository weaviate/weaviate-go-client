package batch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/batch"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		dryRun := true
		output := "verbose"
		id := "5b6a08ba-1d46-43aa-89cc-8b070790c6f2"

		filter := batch.BatchDeleteFilter{
			DryRun: &dryRun,
			Output: &output,
			Match: &batch.BatchDeleteMatch{
				Class: "Pizza",
				Where: &models.WhereFilter{
					Operator:    "Equal",
					Path:        []string{"id"},
					ValueString: &id,
				},
			},
		}

		expectedResults := batch.BatchDeleteResults{
			Matches:    1,
			Failed:     0,
			Successful: 0,
			Limit:      10000,
			Objects: []batch.BatchDeleteResultObject{
				{
					ID:     "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
					Status: "DRYRUN",
				},
			},
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithFilter(&filter).
			Do(context.Background())
		require.Nil(t, err)
		require.NotNil(t, resp.Match)
		assert.Equal(t, "Pizza", resp.Match.Class)
		assert.Equal(t, filter.Match.Where, resp.Match.Where)
		assert.Equal(t, expectedResults, resp.Results)

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("batch delete", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		dryRun := false
		output := "verbose"
		nowString := fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))

		filter := batch.BatchDeleteFilter{
			DryRun: &dryRun,
			Output: &output,
			Match: &batch.BatchDeleteMatch{
				Class: "Pizza",
				Where: &models.WhereFilter{
					Operator:    "LessThan",
					Path:        []string{"_creationTimeUnix"},
					ValueString: &nowString,
				},
			},
		}

		expectedResults := batch.BatchDeleteResults{
			Matches:    4,
			Failed:     0,
			Successful: 4,
			Limit:      10000,
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithFilter(&filter).
			Do(context.Background())
		require.Nil(t, err)
		require.NotNil(t, resp.Match)
		assert.Equal(t, "Pizza", resp.Match.Class)
		assert.Equal(t, filter.Match.Where, resp.Match.Where)
		assert.Equal(t, expectedResults.Matches, resp.Results.Matches)
		assert.Equal(t, expectedResults.Failed, resp.Results.Failed)
		assert.Equal(t, expectedResults.Successful, resp.Results.Successful)
		assert.Len(t, resp.Results.Objects, 4)

		for _, obj := range resp.Results.Objects {
			assert.Equal(t, "SUCCESS", obj.Status)
			assert.Nil(t, obj.Errors)
		}

		testsuit.CleanUpWeaviate(t, client)
	})

	t.Run("batch delete no matches", func(t *testing.T) {
		client := testsuit.CreateTestClient()
		testsuit.CreateTestSchemaAndData(t, client)

		dryRun := true
		output := "verbose"
		id := "267f5125-c9fd-4ca6-9134-f383ff5f0cb6"

		filter := batch.BatchDeleteFilter{
			DryRun: &dryRun,
			Output: &output,
			Match: &batch.BatchDeleteMatch{
				Class: "Pizza",
				Where: &models.WhereFilter{
					Operator:    "Equal",
					Path:        []string{"id"},
					ValueString: &id,
				},
			},
		}

		expectedResults := batch.BatchDeleteResults{
			Matches:    0,
			Failed:     0,
			Successful: 0,
			Limit:      10000,
			Objects:    nil,
		}

		resp, err := client.Batch().ObjectsBatchDeleter().
			WithFilter(&filter).
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
