package backup

import (
	"context"
	"testing"

	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/backup"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuckups_integration(t *testing.T) {
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
	defer testsuit.CleanUpWeaviate(t, client)

	backupID := "1"
	className := "Pizza"

	t.Run("create backup", func(t *testing.T) {
		t.Run("check data exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})

		t.Run("run backup process", func(t *testing.T) {
			meta, err := client.Backup().Creator().
				WithClassName(className).
				WithStorageName(backup.STORAGE_FILESYSTEM).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())
			
			require.Nil(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, models.SnapshotMetaStatusSUCCESS, *meta.Status)
		})

		t.Run("check backup status", func(t *testing.T) {
			meta, err := client.Backup().CreateStatusGetter().
				WithClassName(className).
				WithStorageName(backup.STORAGE_FILESYSTEM).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, models.SnapshotMetaStatusSUCCESS, *meta.Status)
		})

		t.Run("check data still exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})
	})

	t.Run("restore backup", func(t *testing.T) {
		t.Run("remove Pizza class", func(t *testing.T) {
			err := client.Schema().ClassDeleter().
				WithClassName(className).
				Do(context.Background())

			assert.Nil(t, err)
		})

		t.Run("run restore process", func(t *testing.T) {
			meta, err := client.Backup().Restorer().
				WithClassName(className).
				WithStorageName(backup.STORAGE_FILESYSTEM).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, models.SnapshotRestoreMetaStatusSUCCESS, *meta.Status)
		})

		t.Run("check restore status", func(t *testing.T) {
			meta, err := client.Backup().RestoreStatusGetter().
				WithClassName(className).
				WithStorageName(backup.STORAGE_FILESYSTEM).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, meta)
			assert.Equal(t, models.SnapshotRestoreMetaStatusSUCCESS, *meta.Status)
		})

		t.Run("check data again exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})
	})
}

func assertAllPizzasExist(t *testing.T, client *weaviate.Client) {
	resultSet, err := client.GraphQL().
		Get().
		WithClassName("Pizza").
		WithFields(graphql.Field{Name: "name"}).
		Do(context.Background())
	assert.Nil(t, err)

	get := resultSet.Data["Get"].(map[string]interface{})
	pizzas := get["Pizza"].([]interface{})
	assert.Len(t, pizzas, 4)
}
