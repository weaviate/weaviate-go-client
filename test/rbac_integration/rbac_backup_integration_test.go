package rbac_integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestRBACBackupBasicCreation(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	// Create test class first (required for backup to work)
	testClass := &models.Class{
		Class:       "TestBackupClass",
		Description: "Test class for backup",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("TestBackupClass").Do(ctx)
	})

	backend := backup.BACKEND_FILESYSTEM
	backupID := "test1-backup"

	// Create backup
	createResponse, err := client.Backup().Creator().
		WithBackend(backend).
		WithBackupID(backupID).
		WithWaitForCompletion(true).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, createResponse)
	assert.Equal(t, backupID, createResponse.ID)
	assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
}

func TestRBACBackupCreationAndStatusCheck(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	// Create test class first (required for backup to work)
	testClass := &models.Class{
		Class:       "TestBackupClass",
		Description: "Test class for backup",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("TestBackupClass").Do(ctx)
	})

	backend := backup.BACKEND_FILESYSTEM
	backupID := "test2-backup"

	// Create backup
	createResponse, err := client.Backup().Creator().
		WithBackend(backend).
		WithBackupID(backupID).
		WithWaitForCompletion(true).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, createResponse)
	assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)

	// Check backup exists
	statusResponse, err := client.Backup().CreateStatusGetter().
		WithBackend(backend).
		WithBackupID(backupID).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, statusResponse)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *statusResponse.Status)
}

func TestRBACBackupWithClassCreationAndDeletion(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	backend := backup.BACKEND_FILESYSTEM
	backupID := "test3-backup"

	// Create test class first
	testClass := &models.Class{
		Class:       "TestClass",
		Description: "Test class",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err)

	// Create backup
	createResponse, err := client.Backup().Creator().
		WithBackend(backend).
		WithBackupID(backupID).
		WithWaitForCompletion(true).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, createResponse)
	assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)

	// Check backup exists
	statusResponse, err := client.Backup().CreateStatusGetter().
		WithBackend(backend).
		WithBackupID(backupID).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, statusResponse)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *statusResponse.Status)

	// Delete everything in database
	err = client.Schema().ClassDeleter().WithClassName("TestClass").Do(ctx)
	require.NoError(t, err)

	// Check delete worked
	schema, err := client.Schema().Getter().Do(ctx)
	require.NoError(t, err)

	found := false
	for _, class := range schema.Classes {
		if class.Class == "TestClass" {
			found = true
			break
		}
	}
	assert.False(t, found)
}

func TestRBACBackupCreateDeleteRestore(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	backend := backup.BACKEND_FILESYSTEM
	backupID := "test4-backup"

	// Create test class first
	testClass := &models.Class{
		Class:       "TestClass",
		Description: "Test class",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err)

	// Create backup
	createResponse, err := client.Backup().Creator().
		WithBackend(backend).
		WithBackupID(backupID).
		WithWaitForCompletion(true).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, createResponse)
	assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)

	// Check backup exists
	statusResponse, err := client.Backup().CreateStatusGetter().
		WithBackend(backend).
		WithBackupID(backupID).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, statusResponse)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *statusResponse.Status)

	// Delete everything in database
	err = client.Schema().ClassDeleter().WithClassName("TestClass").Do(ctx)
	require.NoError(t, err)

	// Check delete worked
	schema, err := client.Schema().Getter().Do(ctx)
	require.NoError(t, err)

	found := false
	for _, class := range schema.Classes {
		if class.Class == "TestClass" {
			found = true
			break
		}
	}
	assert.False(t, found)

	// Restore backup
	restoreResponse, err := client.Backup().Restorer().
		WithBackend(backend).
		WithBackupID(backupID).
		WithRBACAndUsers().
		WithWaitForCompletion(true).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, restoreResponse)
	assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
}

func TestRBACBackupFullCycleWithValidation(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	backend := backup.BACKEND_FILESYSTEM
	backupID := "test5-backup"

	// Create test class first
	testClass := &models.Class{
		Class:       "TestClass",
		Description: "Test class",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err)

	// Create backup
	createResponse, err := client.Backup().Creator().
		WithBackend(backend).
		WithBackupID(backupID).
		WithWaitForCompletion(true).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, createResponse)
	assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)

	// Check backup exists
	statusResponse, err := client.Backup().CreateStatusGetter().
		WithBackend(backend).
		WithBackupID(backupID).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, statusResponse)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *statusResponse.Status)

	// Delete everything in database
	err = client.Schema().ClassDeleter().WithClassName("TestClass").Do(ctx)
	require.NoError(t, err)

	// Check delete worked
	schema, err := client.Schema().Getter().Do(ctx)
	require.NoError(t, err)

	found := false
	for _, class := range schema.Classes {
		if class.Class == "TestClass" {
			found = true
			break
		}
	}
	assert.False(t, found)

	// Restore backup
	restoreResponse, err := client.Backup().Restorer().
		WithBackend(backend).
		WithBackupID(backupID).
		WithRBACAndUsers().
		WithWaitForCompletion(true).
		Do(ctx)

	require.NoError(t, err)
	require.NotNil(t, restoreResponse)
	assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)

	// Check restore worked
	schema, err = client.Schema().Getter().Do(ctx)
	require.NoError(t, err)

	found = false
	for _, class := range schema.Classes {
		if class.Class == "TestClass" {
			found = true
			break
		}
	}
	assert.True(t, found)
}
