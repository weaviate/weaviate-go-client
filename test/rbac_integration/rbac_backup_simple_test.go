package rbac_integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup/rbac"
	weaviateRbac "github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestBackupAllRBAC(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	defer stop()

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	// Create test class
	class := &models.Class{
		Class: "TestClass",
		Properties: []*models.Property{
			{Name: "title", DataType: []string{"text"}},
		},
	}
	err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
	require.NoError(t, err)

	// Create test role
	role := weaviateRbac.NewRole("test-role", weaviateRbac.DataPermission{
		Actions:    []string{models.PermissionActionReadData},
		Collection: "*",
	})
	err = client.Roles().Creator().WithRole(role).Do(ctx)
	require.NoError(t, err)

	// Create test user
	_, err = client.Users().DB().Creator().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	// Assign role to user
	err = client.Users().DB().RolesAssigner().WithUserID("test-user").WithRoles("test-role").Do(ctx)
	require.NoError(t, err)

	// Create backup with all RBAC
	_, err = client.Backup().Creator().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-all").
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check backup exists
	status, err := client.Backup().CreateStatusGetter().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-all").
		Do(ctx)
	require.NoError(t, err)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *status.Status)

	// Delete everything in database
	err = client.Schema().ClassDeleter().WithClassName("TestClass").Do(ctx)
	require.NoError(t, err)

	_, err = client.Users().DB().Deleter().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	err = client.Roles().Deleter().WithName("test-role").Do(ctx)
	require.NoError(t, err)

	// Check delete worked
	schema, err := client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Empty(t, schema.Classes)

	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	assert.Error(t, err)

	_, err = client.Users().DB().Lister().Do(ctx)
	require.NoError(t, err)

	// Restore backup
	_, err = client.Backup().Restorer().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-all").
		WithRBACAndUsers().
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check restore worked
	schema, err = client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Len(t, schema.Classes, 1)
	assert.Equal(t, "TestClass", schema.Classes[0].Class)

	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	require.NoError(t, err)

	users, err := client.Users().DB().Lister().Do(ctx)
	require.NoError(t, err)
	found := false
	for _, user := range users {
		if user.UserID == "test-user" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestBackupNoRBAC(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	defer stop()

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	// Create test class
	class := &models.Class{
		Class: "TestClass",
		Properties: []*models.Property{
			{Name: "title", DataType: []string{"text"}},
		},
	}
	err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
	require.NoError(t, err)

	// Create test role
	role := weaviateRbac.NewRole("test-role", weaviateRbac.DataPermission{
		Actions:    []string{models.PermissionActionReadData},
		Collection: "*",
	})
	err = client.Roles().Creator().WithRole(role).Do(ctx)
	require.NoError(t, err)

	// Create test user
	_, err = client.Users().DB().Creator().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	// Create backup with no RBAC
	_, err = client.Backup().Creator().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-none").
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check backup exists
	status, err := client.Backup().CreateStatusGetter().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-none").
		Do(ctx)
	require.NoError(t, err)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *status.Status)

	// Delete everything in database
	err = client.Schema().ClassDeleter().WithClassName("TestClass").Do(ctx)
	require.NoError(t, err)

	_, err = client.Users().DB().Deleter().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	err = client.Roles().Deleter().WithName("test-role").Do(ctx)
	require.NoError(t, err)

	// Check delete worked
	schema, err := client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Empty(t, schema.Classes)
	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	assert.Error(t, err)
	_, err = client.Users().DB().Lister().Do(ctx)
	require.NoError(t, err)

	// Restore backup
	_, err = client.Backup().Restorer().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-none").
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check restore worked - class should be restored but no RBAC
	schema, err = client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Len(t, schema.Classes, 1)
	assert.Equal(t, "TestClass", schema.Classes[0].Class)

	// RBAC should NOT be restored
	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	assert.Error(t, err)
}

func TestBackupRolesOnly(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	defer stop()

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	// Create test class
	class := &models.Class{
		Class: "TestClass",
		Properties: []*models.Property{
			{Name: "title", DataType: []string{"text"}},
		},
	}
	err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
	require.NoError(t, err)

	// Create test role
	role := weaviateRbac.NewRole("test-role", weaviateRbac.DataPermission{
		Actions:    []string{models.PermissionActionReadData},
		Collection: "*",
	})
	err = client.Roles().Creator().WithRole(role).Do(ctx)
	require.NoError(t, err)

	// Create test user
	_, err = client.Users().DB().Creator().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	// Assign role to user
	err = client.Users().DB().RolesAssigner().WithUserID("test-user").WithRoles("test-role").Do(ctx)
	require.NoError(t, err)

	// Create backup with roles only
	_, err = client.Backup().Creator().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-roles").
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check backup exists
	status, err := client.Backup().CreateStatusGetter().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-roles").
		Do(ctx)
	require.NoError(t, err)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *status.Status)

	// Delete everything in database
	err = client.Schema().ClassDeleter().WithClassName("TestClass").Do(ctx)
	require.NoError(t, err)

	_, err = client.Users().DB().Deleter().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	err = client.Roles().Deleter().WithName("test-role").Do(ctx)
	require.NoError(t, err)

	// Check delete worked
	schema, err := client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Empty(t, schema.Classes)
	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	assert.Error(t, err)
	_, err = client.Users().DB().Lister().Do(ctx)
	require.NoError(t, err)

	// Restore backup
	_, err = client.Backup().Restorer().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-roles").
		WithRBACRoles(rbac.RBACAll).
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check restore worked - class and role should be restored
	schema, err = client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Len(t, schema.Classes, 1)
	assert.Equal(t, "TestClass", schema.Classes[0].Class)

	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	require.NoError(t, err)

	// Users should NOT be restored
	users, err := client.Users().DB().Lister().Do(ctx)
	require.NoError(t, err)
	found := false
	for _, user := range users {
		if user.UserID == "test-user" {
			found = true
			break
		}
	}
	assert.False(t, found)
}

func TestBackupUsersOnly(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	defer stop()

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	// Create test class
	class := &models.Class{
		Class: "TestClass",
		Properties: []*models.Property{
			{Name: "title", DataType: []string{"text"}},
		},
	}
	err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
	require.NoError(t, err)

	// Create test role
	role := weaviateRbac.NewRole("test-role", weaviateRbac.DataPermission{
		Actions:    []string{models.PermissionActionReadData},
		Collection: "*",
	})
	err = client.Roles().Creator().WithRole(role).Do(ctx)
	require.NoError(t, err)

	// Create test user
	_, err = client.Users().DB().Creator().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	// Assign role to user
	err = client.Users().DB().RolesAssigner().WithUserID("test-user").WithRoles("test-role").Do(ctx)
	require.NoError(t, err)

	// Create backup with users only
	_, err = client.Backup().Creator().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-users").
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check backup exists
	status, err := client.Backup().CreateStatusGetter().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-users").
		Do(ctx)
	require.NoError(t, err)
	assert.Equal(t, models.BackupCreateStatusResponseStatusSUCCESS, *status.Status)

	// Delete everything in database
	err = client.Schema().ClassDeleter().WithClassName("TestClass").Do(ctx)
	require.NoError(t, err)

	_, err = client.Users().DB().Deleter().WithUserID("test-user").Do(ctx)
	require.NoError(t, err)

	err = client.Roles().Deleter().WithName("test-role").Do(ctx)
	require.NoError(t, err)

	// Check delete worked
	schema, err := client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Empty(t, schema.Classes)
	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	assert.Error(t, err)
	_, err = client.Users().DB().Lister().Do(ctx)
	require.NoError(t, err)

	// Restore backup
	_, err = client.Backup().Restorer().
		WithBackend(backup.BACKEND_FILESYSTEM).
		WithBackupID("test-backup-users").
		WithRBACUsers(rbac.UserAll).
		WithWaitForCompletion(true).
		Do(ctx)
	require.NoError(t, err)

	// Check restore worked - class and users should be restored
	schema, err = client.Schema().Getter().Do(ctx)
	require.NoError(t, err)
	assert.Len(t, schema.Classes, 1)
	assert.Equal(t, "TestClass", schema.Classes[0].Class)

	users, err := client.Users().DB().Lister().Do(ctx)
	require.NoError(t, err)
	found := false
	for _, user := range users {
		if user.UserID == "test-user" {
			found = true
			break
		}
	}
	assert.True(t, found)

	// Roles should NOT be restored
	_, err = client.Roles().Getter().WithName("test-role").Do(ctx)
	assert.Error(t, err)
}
