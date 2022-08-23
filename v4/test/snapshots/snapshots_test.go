package schema

import (
	"context"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/snapshots"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshots_integration(t *testing.T) {
	if err := testenv.SetupLocalWeaviate(); err != nil {
		t.Fatalf("failed to setup weaviate: %s", err)
	}
	defer func() {
		if err := testenv.TearDownLocalWeaviate(); err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	}()

	client := testsuit.CreateTestClient()
	testsuit.CreateTestSchemaAndData(t, client)

	snapshotID := "1"
	t.Run("create snapshot", func(t *testing.T) {
		t.Run("get all", func(t *testing.T) {
			resultSet, err := client.GraphQL().
				Get().
				WithClassName("Pizza").
				WithFields(graphql.Field{Name: "name"}).
				Do(context.Background())
			assert.Nil(t, err)

			get := resultSet.Data["Get"].(map[string]interface{})
			pizza := get["Pizza"].([]interface{})
			assert.Equal(t, 4, len(pizza))
		})

		t.Run("run the snapshot process", func(t *testing.T) {
			meta, err := client.Snapshots().
				Creator().
				WithClassName("Pizza").
				WithStorageProvider(snapshots.STORAGE_PROVIDER_FILESYSTEM).
				WithSnapshotID(snapshotID).
				WithWaitForCompletion().
				Do(context.Background())
			require.Nil(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, models.SnapshotMetaStatusSUCCESS, *meta.Status)
		})

		t.Run("get the created snapshot", func(t *testing.T) {
			meta, err := client.Snapshots().
				StatusCreateSnapshot().
				WithClassName("Pizza").
				WithStorageProvider(snapshots.STORAGE_PROVIDER_FILESYSTEM).
				WithSnapshotID(snapshotID).
				Do(context.Background())
			require.Nil(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, models.SnapshotMetaStatusSUCCESS, *meta.Status)
		})
	})

	t.Run("restore snapshot", func(t *testing.T) {
		t.Run("remove Pizza class", func(t *testing.T) {
			err := client.Schema().
				ClassDeleter().
				WithClassName("Pizza").
				Do(context.Background())
			assert.Nil(t, err)
		})

		t.Run("run the restore snapshot process", func(t *testing.T) {
			meta, err := client.Snapshots().
				Restorer().
				WithClassName("Pizza").
				WithStorageProvider(snapshots.STORAGE_PROVIDER_FILESYSTEM).
				WithSnapshotID(snapshotID).
				WithWaitForCompletion().
				Do(context.Background())
			require.Nil(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, models.SnapshotRestoreMetaStatusSUCCESS, *meta.Status)
		})

		t.Run("check that Pizza class was restored", func(t *testing.T) {
			pizzas, err := client.Data().
				ObjectsGetter().
				WithClassName("Pizza").
				Do(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, 4, len(pizzas))
		})
	})
}
