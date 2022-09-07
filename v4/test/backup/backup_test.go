package backup

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/semi-technologies/weaviate-go-client/v4/test/testsuit"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/backup"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/testenv"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackups_integration(t *testing.T) {
	if err := testenv.SetupLocalWeaviate(); err != nil {
		t.Fatalf("failed to setup weaviate: %s", err)
	}
	defer func() {
		if err := testenv.TearDownLocalWeaviate(); err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	}()

	const dockerComposeBackupDir = "/tmp/backups"
	client := testsuit.CreateTestClient()
	testsuit.CleanUpWeaviate(t, client)

	t.Run("create and restore backup with waiting", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "Pizza"

		t.Run("check data exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createResponse)
			assert.Equal(t, backupID, createResponse.ID)
			assert.Len(t, createResponse.Classes, 1)
			assert.Contains(t, createResponse.Classes, className)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, createResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, createResponse.Backend)
			assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
			assert.Empty(t, createResponse.Error)
		})

		t.Run("check data still exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})

		t.Run("check create status", func(t *testing.T) {
			createStatusResponse, err := client.Backup().CreateStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createStatusResponse)
			assert.Equal(t, backupID, createStatusResponse.ID)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, createStatusResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, createStatusResponse.Backend)
			assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *createStatusResponse.Status)
			assert.Empty(t, createStatusResponse.Error)
		})

		t.Run("remove existing class", func(t *testing.T) {
			err := client.Schema().ClassDeleter().
				WithClassName(className).
				Do(context.Background())

			assert.Nil(t, err)
		})

		t.Run("restore backup", func(t *testing.T) {
			restoreResponse, err := client.Backup().Restorer().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, restoreResponse)
			assert.Equal(t, backupID, restoreResponse.ID)
			assert.Len(t, restoreResponse.Classes, 1)
			assert.Contains(t, restoreResponse.Classes, className)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, restoreResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, restoreResponse.Backend)
			assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
			assert.Empty(t, restoreResponse.Error)
		})

		t.Run("check data again exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})

		t.Run("check restore status", func(t *testing.T) {
			restoreStatusResponse, err := client.Backup().RestoreStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, restoreStatusResponse)

			assert.Equal(t, backupID, restoreStatusResponse.ID)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, restoreStatusResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, restoreStatusResponse.Backend)
			assert.Equal(t, models.BackupRestoreStatusResponseStatusSUCCESS, *restoreStatusResponse.Status)
			assert.Empty(t, restoreStatusResponse.Error)
		})
	})

	t.Run("create and restore backup without waiting", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "Pizza"

		t.Run("check data exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createResponse)
			assert.Equal(t, backupID, createResponse.ID)
			assert.Len(t, createResponse.Classes, 1)
			assert.Contains(t, createResponse.Classes, className)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, createResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, createResponse.Backend)
			assert.Equal(t, models.BackupCreateResponseStatusSTARTED, *createResponse.Status)
			assert.Empty(t, createResponse.Error)
		})

		t.Run("wait until created", func(t *testing.T) {
			statusGetter := client.Backup().CreateStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID)

			for {
				createStatusResponse, err := statusGetter.Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, createStatusResponse)
				assert.Equal(t, backupID, createStatusResponse.ID)
				assert.Equal(t, dockerComposeBackupDir+"/"+backupID, createStatusResponse.Path)
				assert.Equal(t, backup.BACKEND_FILESYSTEM, createStatusResponse.Backend)
				assert.Empty(t, createStatusResponse.Error)
				assert.Contains(t, []string{models.BackupCreateStatusResponseStatusSTARTED,
					models.BackupCreateStatusResponseStatusTRANSFERRING,
					models.BackupCreateStatusResponseStatusTRANSFERRED,
					models.BackupCreateStatusResponseStatusSUCCESS}, *createStatusResponse.Status)

				if models.BackupCreateStatusResponseStatusSUCCESS == *createStatusResponse.Status {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		})

		t.Run("check data still exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})

		t.Run("remove existing class", func(t *testing.T) {
			err := client.Schema().ClassDeleter().
				WithClassName(className).
				Do(context.Background())

			assert.Nil(t, err)
		})

		t.Run("restore backup", func(t *testing.T) {
			restoreResponse, err := client.Backup().Restorer().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, restoreResponse)
			assert.Equal(t, backupID, restoreResponse.ID)
			assert.Len(t, restoreResponse.Classes, 1)
			assert.Contains(t, restoreResponse.Classes, className)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, restoreResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, restoreResponse.Backend)
			assert.Equal(t, models.BackupRestoreResponseStatusSTARTED, *restoreResponse.Status)
			assert.Empty(t, restoreResponse.Error)
		})

		t.Run("wait until restored", func(t *testing.T) {
			statusGetter := client.Backup().RestoreStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID)

			for {
				restoreStatusResponse, err := statusGetter.Do(context.Background())

				require.Nil(t, err)
				require.NotNil(t, restoreStatusResponse)
				assert.Equal(t, backupID, restoreStatusResponse.ID)
				assert.Equal(t, dockerComposeBackupDir+"/"+backupID, restoreStatusResponse.Path)
				assert.Equal(t, backup.BACKEND_FILESYSTEM, restoreStatusResponse.Backend)
				assert.Empty(t, restoreStatusResponse.Error)
				assert.Contains(t, []string{models.BackupRestoreStatusResponseStatusSTARTED,
					models.BackupRestoreStatusResponseStatusTRANSFERRING,
					models.BackupRestoreStatusResponseStatusTRANSFERRED,
					models.BackupRestoreStatusResponseStatusSUCCESS}, *restoreStatusResponse.Status)

				if models.BackupRestoreStatusResponseStatusSUCCESS == *restoreStatusResponse.Status {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		})

		t.Run("check data again exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
		})
	})

	t.Run("create and restore 1 of 2 classes", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		pizzaClassName := "Pizza"
		soupClassName := "Soup"

		t.Run("check data exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
			assertAllSoupsExist(t, client)
		})

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createResponse)
			assert.Equal(t, backupID, createResponse.ID)
			assert.Len(t, createResponse.Classes, 2)
			assert.Contains(t, createResponse.Classes, pizzaClassName)
			assert.Contains(t, createResponse.Classes, soupClassName)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, createResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, createResponse.Backend)
			assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
			assert.Empty(t, createResponse.Error)
		})

		t.Run("check data still exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
			assertAllSoupsExist(t, client)
		})

		t.Run("check create status", func(t *testing.T) {
			createStatusResponse, err := client.Backup().CreateStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createStatusResponse)
			assert.Equal(t, backupID, createStatusResponse.ID)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, createStatusResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, createStatusResponse.Backend)
			assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *createStatusResponse.Status)
			assert.Empty(t, createStatusResponse.Error)
		})

		t.Run("remove existing class", func(t *testing.T) {
			err := client.Schema().ClassDeleter().
				WithClassName(pizzaClassName).
				Do(context.Background())

			assert.Nil(t, err)
		})

		t.Run("restore backup", func(t *testing.T) {
			restoreResponse, err := client.Backup().Restorer().
				WithIncludeClassNames(pizzaClassName).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, restoreResponse)
			assert.Equal(t, backupID, restoreResponse.ID)
			assert.Len(t, restoreResponse.Classes, 1)
			assert.Contains(t, restoreResponse.Classes, pizzaClassName)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, restoreResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, restoreResponse.Backend)
			assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
			assert.Empty(t, restoreResponse.Error)
		})

		t.Run("check data again exist", func(t *testing.T) {
			assertAllPizzasExist(t, client)
			assertAllSoupsExist(t, client)
		})

		t.Run("check restore status", func(t *testing.T) {
			restoreStatusResponse, err := client.Backup().RestoreStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, restoreStatusResponse)

			assert.Equal(t, backupID, restoreStatusResponse.ID)
			assert.Equal(t, dockerComposeBackupDir+"/"+backupID, restoreStatusResponse.Path)
			assert.Equal(t, backup.BACKEND_FILESYSTEM, restoreStatusResponse.Backend)
			assert.Equal(t, models.BackupRestoreStatusResponseStatusSUCCESS, *restoreStatusResponse.Status)
			assert.Empty(t, restoreStatusResponse.Error)
		})
	})

	t.Run("fail creating backup on not existing backend", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := "not-existing-backend"
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "Pizza"

		t.Run("fail create", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, createResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), backend)
		})
	})

	t.Run("fail checking create status on not existing backend", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := "not-existing-backend"
		backupID := uuid.Must(uuid.NewRandom()).String()

		t.Run("fail check status", func(t *testing.T) {
			createStatusResponse, err := client.Backup().CreateStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, createStatusResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), backend)
		})
	})

	t.Run("fail restoring backup on not existing backend", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := "non-existing-backend"
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "not-existing-class"

		t.Run("fails restore", func(t *testing.T) {
			restoreResponse, err := client.Backup().Restorer().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, restoreResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), backend)
		})
	})

	t.Run("fail creating backup for not existing class", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "not-existing-class"

		t.Run("fail create", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, createResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), className)
		})
	})

	t.Run("fail restoring backup for existing class", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "Pizza"

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createResponse)
		})

		t.Run("fail restore", func(t *testing.T) {
			restoreResponse, err := client.Backup().Restorer().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, restoreResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), className)
		})
	})

	t.Run("fail creating existing backup", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "Pizza"

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createResponse)
		})

		t.Run("fail create", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, createResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), backupID)
		})
	})

	t.Run("fail checking create status for not existing backup", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()

		t.Run("fail check status", func(t *testing.T) {
			createStatusResponse, err := client.Backup().CreateStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, createStatusResponse)
			assert.Contains(t, err.Error(), "404")
			assert.Contains(t, err.Error(), backupID)
		})
	})

	t.Run("fail restoring not existing backup", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "Pizza"

		t.Run("fail restore", func(t *testing.T) {
			restoreResponse, err := client.Backup().Restorer().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, restoreResponse)
			assert.Contains(t, err.Error(), "404")
			assert.Contains(t, err.Error(), backupID)
		})
	})

	t.Run("fail checking restore status for not started restore", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		className := "Pizza"

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(className).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createResponse)
		})

		t.Run("fail restore", func(t *testing.T) {
			restoreStatusResponse, err := client.Backup().RestoreStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, restoreStatusResponse)
			assert.Contains(t, err.Error(), "404")
			assert.Contains(t, err.Error(), backend)
			assert.Contains(t, err.Error(), backupID)
		})
	})

	t.Run("fail creating backup for both include and exclude classes", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		pizzaClassName := "Pizza"
		soupClassName := "Soup"

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(pizzaClassName).
				WithExcludeClassNames(soupClassName).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, createResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), "include")
			assert.Contains(t, err.Error(), "exclude")
		})
	})

	t.Run("fail restoring backup for both include and exclude classes", func(t *testing.T) {
		testsuit.CreateTestSchemaAndData(t, client)
		defer testsuit.CleanUpWeaviate(t, client)

		backend := backup.BACKEND_FILESYSTEM
		backupID := uuid.Must(uuid.NewRandom()).String()
		pizzaClassName := "Pizza"
		soupClassName := "Soup"

		t.Run("create backup", func(t *testing.T) {
			createResponse, err := client.Backup().Creator().
				WithIncludeClassNames(pizzaClassName, soupClassName).
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(context.Background())

			require.Nil(t, err)
			require.NotNil(t, createResponse)
		})

		t.Run("remove existing class", func(t *testing.T) {
			err := client.Schema().ClassDeleter().
				WithClassName(pizzaClassName).
				Do(context.Background())

			assert.Nil(t, err)
		})

		t.Run("fail restore", func(t *testing.T) {
			restoreResponse, err := client.Backup().Restorer().
				WithIncludeClassNames(pizzaClassName).
				WithExcludeClassNames(soupClassName).
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())

			require.NotNil(t, err)
			require.Nil(t, restoreResponse)
			assert.Contains(t, err.Error(), "422")
			assert.Contains(t, err.Error(), "include")
			assert.Contains(t, err.Error(), "exclude")
		})
	})
}

func assertAllPizzasExist(t *testing.T, client *weaviate.Client) {
	assertAllFoodObjectsExist(t, client, "Pizza", 4)
}

func assertAllSoupsExist(t *testing.T, client *weaviate.Client) {
	assertAllFoodObjectsExist(t, client, "Soup", 2)
}

func assertAllFoodObjectsExist(t *testing.T, client *weaviate.Client, className string, count int) {
	resultSet, err := client.GraphQL().
		Get().
		WithClassName(className).
		WithFields(graphql.Field{Name: "name"}).
		Do(context.Background())
	assert.Nil(t, err)

	get := resultSet.Data["Get"].(map[string]interface{})
	objects := get[className].([]interface{})
	assert.Len(t, objects, count)
}
