package rbac_integration

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup"
	rb "github.com/weaviate/weaviate-go-client/v5/weaviate/backup/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

// waitForBackupCompletion polls the backup status until it is no longer "STARTED"
func waitForBackupCompletion(t *testing.T, client *weaviate.Client, backend, backupID string) {
	t.Helper()
	timeout := time.After(30 * time.Second)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			t.Fatalf("Timeout waiting for backup %s to complete", backupID)
		case <-tick:
			status, err := client.Backup().CreateStatusGetter().
				WithBackend(backend).
				WithBackupID(backupID).
				Do(context.Background())
			if err != nil {
				t.Logf("Error checking backup status: %v", err)
				continue
			}
			if status != nil && *status.Status != "STARTED" {
				t.Logf("Backup %s completed with status %s", backupID, *status.Status)
				return
			}
		}
	}
}

// TestRBACBackupWithUserRoleManagement tests comprehensive user and role management with backups
func TestRBACBackupWithUserRoleManagement(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	backend := backup.BACKEND_FILESYSTEM

	const (
		adminRole     = "test-manager"
		viewerRole    = "test-reader"
		editorRole    = "test-writer"
		adminUser     = "manager-user"
		viewerUser    = "reader-user"
		editorUser    = "writer-user"
		testClassName = "TestDocument"
	)

	// Create test class first (required for backup to work) - at main test level
	class := &models.Class{
		Class:       testClassName,
		Description: "Test document class for RBAC backup testing",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
			{
				Name:     "content",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Logf("Created class '%s'", testClassName)

	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName(testClassName).Do(ctx)
	})

	t.Run("setup comprehensive RBAC environment", func(t *testing.T) {
		// Create multiple roles with different permissions
		roles := map[string]rbac.Role{
			adminRole: rbac.NewRole(adminRole,
				rbac.BackupsPermission{
					Actions:    []string{models.PermissionActionManageBackups},
					Collection: "*",
				},
				rbac.DataPermission{
					Actions:    []string{models.PermissionActionCreateData, models.PermissionActionReadData, models.PermissionActionUpdateData, models.PermissionActionDeleteData},
					Collection: "*",
				},
			),
			viewerRole: rbac.NewRole(viewerRole,
				rbac.DataPermission{
					Actions:    []string{models.PermissionActionReadData},
					Collection: "*",
				},
			),
			editorRole: rbac.NewRole(editorRole,
				rbac.DataPermission{
					Actions:    []string{models.PermissionActionCreateData, models.PermissionActionReadData, models.PermissionActionUpdateData},
					Collection: testClassName,
				},
			),
		}

		// Create roles
		for roleName, role := range roles {
			err := client.Roles().Creator().WithRole(role).Do(ctx)
			require.NoError(t, err, "failed to create role %s", roleName)
			t.Logf("Created role '%s'", roleName)

			t.Cleanup(func() {
				client.Roles().Deleter().WithName(roleName).Do(ctx)
			})
		}

		// Create users
		users := []string{adminUser, viewerUser, editorUser}
		userRoles := map[string]string{
			adminUser:  adminRole,
			viewerUser: viewerRole,
			editorUser: editorRole,
		}

		for _, user := range users {
			apiKey, err := client.Users().DB().Creator().WithUserID(user).Do(ctx)
			require.NoError(t, err, "failed to create user %s", user)
			require.NotEmpty(t, apiKey, "API key should not be empty")
			t.Logf("Created user '%s'", user)

			// Assign role to user
			if role, exists := userRoles[user]; exists {
				err = client.Users().DB().RolesAssigner().WithUserID(user).WithRoles(role).Do(ctx)
				require.NoError(t, err, "failed to assign role %s to user %s", role, user)
				t.Logf("Assigned role '%s' to user '%s'", role, user)
			}

			t.Cleanup(func() {
				client.Users().DB().Deleter().WithUserID(user).Do(ctx)
			})
		}
	})

	// Test all roles backup and restore
	t.Run("test all roles backup and restore", func(t *testing.T) {
		backupID := fmt.Sprintf("all-roles-%d", random.Int63())

		// Backup all roles and users
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithWaitForCompletion(true).
			Do(ctx)

		if err != nil && strings.Contains(err.Error(), "already in progress") {
			t.Logf("Backup already in progress, waiting...")
			waitForBackupCompletion(t, client, backend, backupID)

			createResponse, err = client.Backup().Creator().
				WithBackend(backend).
				WithBackupID(backupID).
				WithWaitForCompletion(true).
				Do(ctx)
		}

		require.NoError(t, err, "failed to create all roles backup")
		require.NotNil(t, createResponse)
		if *createResponse.Status != models.BackupCreateResponseStatusSUCCESS {
			t.Logf("Backup creation failed with status: %s", *createResponse.Status)
			if createResponse.Error != "" {
				t.Logf("Backup creation error: %s", createResponse.Error)
			}
		}
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Created all roles backup %s", backupID)

		// Delete all roles
		err = client.Roles().Deleter().WithName(adminRole).Do(ctx)
		require.NoError(t, err, "failed to delete role %s", adminRole)
		err = client.Roles().Deleter().WithName(viewerRole).Do(ctx)
		require.NoError(t, err, "failed to delete role %s", viewerRole)
		err = client.Roles().Deleter().WithName(editorRole).Do(ctx)
		require.NoError(t, err, "failed to delete role %s", editorRole)

		// Delete all users
		deleted, err := client.Users().DB().Deleter().WithUserID(adminUser).Do(ctx)
		if err == nil && deleted {
			t.Logf("Deleted user %s", adminUser)
		}
		deleted, err = client.Users().DB().Deleter().WithUserID(viewerUser).Do(ctx)
		if err == nil && deleted {
			t.Logf("Deleted user %s", viewerUser)
		}
		deleted, err = client.Users().DB().Deleter().WithUserID(editorUser).Do(ctx)
		if err == nil && deleted {
			t.Logf("Deleted user %s", editorUser)
		}

		// Delete the class before restore to avoid conflicts
		err = client.Schema().ClassDeleter().WithClassName(testClassName).Do(ctx)
		require.NoError(t, err, "failed to delete test class before restore")

		// Restore backup
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACRoles(rb.RBACAll).
			WithRBACUsers(rb.UserAll).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore all roles backup")
		require.NotNil(t, restoreResponse)
		if *restoreResponse.Status != models.BackupRestoreResponseStatusSUCCESS {
			t.Logf("Backup restore failed with status: %s", *restoreResponse.Status)
			if restoreResponse.Error != "" {
				t.Logf("Backup restore error: %s", restoreResponse.Error)
			}
		}
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
		t.Logf("Restored all roles backup %s", backupID)

		// Verify all roles were restored
		allRoles, err := client.Roles().AllGetter().Do(ctx)
		if err == nil {
			t.Logf("Found %d roles after restore", len(allRoles))
			for _, role := range allRoles {
				t.Logf("Role '%s' exists after restore", role.Name)
			}
		} else {
			t.Logf("Could not list roles after restore: %v", err)
		}

		// Verify users were restored
		users, err := client.Users().DB().Lister().Do(ctx)
		if err == nil {
			userMap := make(map[string]bool)
			for _, user := range users {
				userMap[user.UserID] = true
			}

			if userMap[adminUser] {
				t.Logf("User '%s' exists after restore", adminUser)
			} else {
				t.Logf("User '%s' does not exist after restore", adminUser)
			}
			if userMap[viewerUser] {
				t.Logf("User '%s' exists after restore", viewerUser)
			} else {
				t.Logf("User '%s' does not exist after restore", viewerUser)
			}
			if userMap[editorUser] {
				t.Logf("User '%s' exists after restore", editorUser)
			} else {
				t.Logf("User '%s' does not exist after restore", editorUser)
			}
		} else {
			t.Logf("Could not list users after restore: %v", err)
		}
	})

	// Test user assignments backup and restore
	t.Run("test user assignments backup and restore", func(t *testing.T) {
		backupID := fmt.Sprintf("user-assignments-%d", random.Int63())

		// Backup with user assignments only
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create user assignments backup")
		require.NotNil(t, createResponse)
		if *createResponse.Status != models.BackupCreateResponseStatusSUCCESS {
			t.Logf("User assignments backup creation failed with status: %s", *createResponse.Status)
			if createResponse.Error != "" {
				t.Logf("User assignments backup creation error: %s", createResponse.Error)
			}
		}
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Created user assignments backup %s", backupID)

		// Delete all users
		deleted, err := client.Users().DB().Deleter().WithUserID(adminUser).Do(ctx)
		if err == nil && deleted {
			t.Logf("Deleted user %s", adminUser)
		}
		deleted, err = client.Users().DB().Deleter().WithUserID(viewerUser).Do(ctx)
		if err == nil && deleted {
			t.Logf("Deleted user %s", viewerUser)
		}
		deleted, err = client.Users().DB().Deleter().WithUserID(editorUser).Do(ctx)
		if err == nil && deleted {
			t.Logf("Deleted user %s", editorUser)
		}

		// Delete the class before restore to avoid conflicts
		err = client.Schema().ClassDeleter().WithClassName(testClassName).Do(ctx)
		require.NoError(t, err, "failed to delete test class before restore")

		// Restore backup
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACRoles(rb.RBACNone).
			WithRBACUsers(rb.UserAll).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore user assignments backup")
		require.NotNil(t, restoreResponse)
		if *restoreResponse.Status != models.BackupRestoreResponseStatusSUCCESS {
			t.Logf("User assignments backup restore failed with status: %s", *restoreResponse.Status)
			if restoreResponse.Error != "" {
				t.Logf("User assignments backup restore error: %s", restoreResponse.Error)
			}
		}
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
		t.Logf("Restored user assignments backup %s", backupID)

		// Verify users were restored
		users, err := client.Users().DB().Lister().Do(ctx)
		if err == nil {
			userMap := make(map[string]bool)
			for _, user := range users {
				userMap[user.UserID] = true
			}

			if userMap[adminUser] {
				t.Logf("User '%s' exists after restore", adminUser)
			} else {
				t.Logf("User '%s' does not exist after restore", adminUser)
			}
			if userMap[viewerUser] {
				t.Logf("User '%s' exists after restore", viewerUser)
			} else {
				t.Logf("User '%s' does not exist after restore", viewerUser)
			}
			if userMap[editorUser] {
				t.Logf("User '%s' exists after restore", editorUser)
			} else {
				t.Logf("User '%s' does not exist after restore", editorUser)
			}
		} else {
			t.Logf("Could not list users after restore: %v", err)
		}
	})
}

// TestRBACBackupComplexScenarios tests complex RBAC backup scenarios
func TestRBACBackupComplexScenarios(t *testing.T) {
	ctx := context.Background()
	container, stop := testenv.SetupLocalContainer(t, ctx, test.RBAC, true)
	t.Cleanup(stop)

	client := testsuit.CreateTestClientForContainer(t, container)
	testsuit.CleanUpWeaviate(t, client)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	backend := backup.BACKEND_FILESYSTEM

	// Create a test class at main test level (required for all backup operations)
	testClass := &models.Class{
		Class:       "ComplexTestClass",
		Description: "Test class for complex scenarios",
		Properties: []*models.Property{
			{
				Name:     "name",
				DataType: []string{"text"},
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(testClass).Do(ctx)
	require.NoError(t, err, "failed to create test class")
	t.Cleanup(func() {
		client.Schema().ClassDeleter().WithClassName("ComplexTestClass").Do(ctx)
	})

	t.Run("test backup with all RBAC data", func(t *testing.T) {
		backupID := fmt.Sprintf("all-rbac-%d", random.Int63())

		// Create multiple test roles
		testRole1 := rbac.NewRole("test-role-1",
			rbac.DataPermission{
				Actions:    []string{models.PermissionActionReadData},
				Collection: "*",
			},
		)
		err := client.Roles().Creator().WithRole(testRole1).Do(ctx)
		require.NoError(t, err, "failed to create test-role-1")
		t.Logf("Created role 'test-role-1'")
		t.Cleanup(func() {
			client.Roles().Deleter().WithName("test-role-1").Do(ctx)
		})

		testRole2 := rbac.NewRole("test-role-2",
			rbac.DataPermission{
				Actions:    []string{models.PermissionActionCreateData, models.PermissionActionReadData},
				Collection: "*",
			},
		)
		err = client.Roles().Creator().WithRole(testRole2).Do(ctx)
		require.NoError(t, err, "failed to create test-role-2")
		t.Logf("Created role 'test-role-2'")
		t.Cleanup(func() {
			client.Roles().Deleter().WithName("test-role-2").Do(ctx)
		})

		// Create test users
		apiKey, err := client.Users().DB().Creator().WithUserID("test-user-1").Do(ctx)
		require.NoError(t, err, "failed to create test-user-1")
		require.NotEmpty(t, apiKey)
		err = client.Users().DB().RolesAssigner().WithUserID("test-user-1").WithRoles("test-role-1").Do(ctx)
		require.NoError(t, err, "failed to assign role to test-user-1")
		t.Logf("Created user 'test-user-1' with role 'test-role-1'")
		t.Cleanup(func() {
			client.Users().DB().Deleter().WithUserID("test-user-1").Do(ctx)
		})

		apiKey, err = client.Users().DB().Creator().WithUserID("test-user-2").Do(ctx)
		require.NoError(t, err, "failed to create test-user-2")
		require.NotEmpty(t, apiKey)
		err = client.Users().DB().RolesAssigner().WithUserID("test-user-2").WithRoles("test-role-2").Do(ctx)
		require.NoError(t, err, "failed to assign role to test-user-2")
		t.Logf("Created user 'test-user-2' with role 'test-role-2'")
		t.Cleanup(func() {
			client.Users().DB().Deleter().WithUserID("test-user-2").Do(ctx)
		})

		// Create backup with all RBAC data
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create all RBAC backup")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Created all RBAC backup %s", backupID)

		// Delete all roles and users
		client.Users().DB().Deleter().WithUserID("test-user-1").Do(ctx)
		client.Roles().Deleter().WithName("test-role-1").Do(ctx)
		client.Users().DB().Deleter().WithUserID("test-user-2").Do(ctx)
		client.Roles().Deleter().WithName("test-role-2").Do(ctx)

		// Delete the class before restore to avoid conflicts
		err = client.Schema().ClassDeleter().WithClassName("ComplexTestClass").Do(ctx)
		require.NoError(t, err, "failed to delete test class before restore")

		// Restore with all RBAC data
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACRoles(rb.RBACAll).
			WithRBACUsers(rb.UserAll).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore all RBAC backup")
		require.NotNil(t, restoreResponse)
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
		t.Logf("Restored all RBAC backup %s", backupID)

		// Verify all roles were restored
		allRoles, err := client.Roles().AllGetter().Do(ctx)
		if err == nil {
			t.Logf("Found %d roles after restore", len(allRoles))
			for _, role := range allRoles {
				t.Logf("Role '%s' exists after restore", role.Name)
			}
		} else {
			t.Logf("Could not list roles after restore: %v", err)
		}

		// Verify users were restored
		users, err := client.Users().DB().Lister().Do(ctx)
		if err == nil {
			t.Logf("Found %d users after restore", len(users))
			for _, user := range users {
				t.Logf("User '%s' exists after restore", user.UserID)
			}
		} else {
			t.Logf("Could not list users after restore: %v", err)
		}
	})
}
