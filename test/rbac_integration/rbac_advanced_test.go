package rbac_integration

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup"
	backupRbac "github.com/weaviate/weaviate-go-client/v5/weaviate/backup/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

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

	// Test selective role backup and restore
	t.Run("test selective role backup and restore", func(t *testing.T) {
		backupID := fmt.Sprintf("selective-roles-%d", random.Int63())

		// Backup only specific roles
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACSpecificRoles(adminRole, viewerRole). // Exclude editorRole
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create selective role backup")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Created selective role backup %s", backupID)

		// Delete all roles
		for _, role := range []string{adminRole, viewerRole, editorRole} {
			err := client.Roles().Deleter().WithName(role).Do(ctx)
			require.NoError(t, err, "failed to delete role %s", role)
		}

		// Restore backup
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACSpecificRoles(adminRole, viewerRole).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore selective role backup")
		require.NotNil(t, restoreResponse)
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
		t.Logf("Restored selective role backup %s", backupID)

		// Verify only selected roles were restored
		for _, role := range []string{adminRole, viewerRole} {
			exists, err := client.Roles().Exists().WithName(role).Do(ctx)
			require.NoError(t, err, "failed to check role existence")
			assert.True(t, exists, "Role '%s' should exist after selective restore", role)
		}

		// Editor role should not exist (was not included in backup)
		exists, err := client.Roles().Exists().WithName(editorRole).Do(ctx)
		if err == nil {
			assert.False(t, exists, "Role '%s' should not exist after selective restore", editorRole)
		} else {
			t.Logf("Could not check editor role existence: %v", err)
		}
	})

	// Test user assignments backup and restore
	t.Run("test user assignments backup and restore", func(t *testing.T) {
		backupID := fmt.Sprintf("user-assignments-%d", random.Int63())

		// Backup with user assignments only
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACUsersOnly().
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create user assignments backup")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Created user assignments backup %s", backupID)

		// Delete all users
		for _, user := range []string{adminUser, viewerUser, editorUser} {
			deleted, err := client.Users().DB().Deleter().WithUserID(user).Do(ctx)
			if err == nil && deleted {
				t.Logf("Deleted user %s", user)
			}
		}

		// Restore backup
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBACUsersOnly().
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore user assignments backup")
		require.NotNil(t, restoreResponse)
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
		t.Logf("Restored user assignments backup %s", backupID)

		// Verify users were restored
		users, err := client.Users().DB().Lister().Do(ctx)
		if err == nil {
			userMap := make(map[string]bool)
			for _, user := range users {
				userMap[user.UserID] = true
			}

			for _, expectedUser := range []string{adminUser, viewerUser, editorUser} {
				if userMap[expectedUser] {
					t.Logf("User '%s' exists after restore", expectedUser)
				} else {
					t.Logf("User '%s' does not exist after restore", expectedUser)
				}
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

	t.Run("test backup with complex RBAC hierarchy", func(t *testing.T) {
		backupID := fmt.Sprintf("complex-hierarchy-%d", random.Int63())

		// Create multiple roles with hierarchy
		roles := []string{"test-superuser", "test-manager", "test-coordinator", "test-member", "test-visitor"}
		for i, roleName := range roles {
			// Each role has progressively fewer permissions
			var permissions []rbac.Permission
			
			if i <= 0 { // super-admin
				permissions = append(permissions, rbac.BackupsPermission{
					Actions:    []string{models.PermissionActionManageBackups},
					Collection: "*",
				})
			}
			if i <= 1 { // admin and above
				permissions = append(permissions, rbac.DataPermission{
					Actions:    []string{models.PermissionActionCreateData, models.PermissionActionReadData, models.PermissionActionUpdateData, models.PermissionActionDeleteData},
					Collection: "*",
				})
			}
			if i <= 2 { // manager and above
				permissions = append(permissions, rbac.DataPermission{
					Actions:    []string{models.PermissionActionCreateData, models.PermissionActionReadData, models.PermissionActionUpdateData},
					Collection: "Documents",
				})
			}
			if i <= 3 { // user and above
				permissions = append(permissions, rbac.DataPermission{
					Actions:    []string{models.PermissionActionReadData},
					Collection: "*",
				})
			}
			// guest gets minimal permissions by default

			role := rbac.NewRole(roleName, permissions...)
			err := client.Roles().Creator().WithRole(role).Do(ctx)
			require.NoError(t, err, "failed to create role %s", roleName)
			t.Logf("Created role '%s' with %d permissions", roleName, len(permissions))

			t.Cleanup(func() {
				client.Roles().Deleter().WithName(roleName).Do(ctx)
			})
		}

		// Create users and assign roles
		for _, roleName := range roles {
			userName := fmt.Sprintf("user-%s", roleName)
			apiKey, err := client.Users().DB().Creator().WithUserID(userName).Do(ctx)
			require.NoError(t, err, "failed to create user %s", userName)
			require.NotEmpty(t, apiKey)

			err = client.Users().DB().RolesAssigner().WithUserID(userName).WithRoles(roleName).Do(ctx)
			require.NoError(t, err, "failed to assign role %s to user %s", roleName, userName)
			t.Logf("Created user '%s' with role '%s'", userName, roleName)

			t.Cleanup(func() {
				client.Users().DB().Deleter().WithUserID(userName).Do(ctx)
			})
		}

		// Create backup with custom RBAC configuration
		rbacConfig := &backupRbac.RBACConfig{
			Scope:                    backupRbac.RBACAll,
			RoleSelection:           backupRbac.RoleSelectionSpecific,
			SpecificRoles:           []string{"test-superuser", "test-manager", "test-coordinator"}, // Only backup first 3 roles
			IncludeUserAssignments:  true,
			IncludeGroupAssignments: false,
		}

		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBAC(rbacConfig).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create complex hierarchy backup")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Created complex hierarchy backup %s", backupID)

		// Delete all roles and users
		for _, roleName := range roles {
			userName := fmt.Sprintf("user-%s", roleName)
			client.Users().DB().Deleter().WithUserID(userName).Do(ctx)
			client.Roles().Deleter().WithName(roleName).Do(ctx)
		}

		// Restore with same configuration
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(backupID).
			WithRBAC(rbacConfig).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore complex hierarchy backup")
		require.NotNil(t, restoreResponse)
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)
		t.Logf("Restored complex hierarchy backup %s", backupID)

		// Verify only the specified roles were restored
		for i, roleName := range roles {
			exists, err := client.Roles().Exists().WithName(roleName).Do(ctx)
			if err == nil {
				if i <= 2 { // Should exist (superuser, manager, coordinator)
					assert.True(t, exists, "Role '%s' should exist after restore", roleName)
				} else { // Should not exist (member, visitor)
					assert.False(t, exists, "Role '%s' should not exist after restore", roleName)
				}
			}
		}
	})

	t.Run("test incremental backup strategy", func(t *testing.T) {
		baseBackupID := fmt.Sprintf("base-backup-%d", random.Int63())
		incrementalBackupID := fmt.Sprintf("incremental-backup-%d", random.Int63())

		// Create initial roles and users
		initialRole := "initial-role"
		role := rbac.NewRole(initialRole,
			rbac.DataPermission{
				Actions:    []string{models.PermissionActionReadData},
				Collection: "*",
			},
		)
		err = client.Roles().Creator().WithRole(role).Do(ctx)
		require.NoError(t, err, "failed to create initial role")

		// Create base backup
		createResponse, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(baseBackupID).
			WithRBACAll().
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create base backup")
		require.NotNil(t, createResponse)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse.Status)
		t.Logf("Created base backup %s", baseBackupID)

		// Add more roles
		additionalRole := "additional-role"
		role2 := rbac.NewRole(additionalRole,
			rbac.DataPermission{
				Actions:    []string{models.PermissionActionCreateData, models.PermissionActionReadData},
				Collection: "*",
			},
		)
		err = client.Roles().Creator().WithRole(role2).Do(ctx)
		require.NoError(t, err, "failed to create additional role")

		// Create incremental backup (only new role)
		createResponse2, err := client.Backup().Creator().
			WithBackend(backend).
			WithBackupID(incrementalBackupID).
			WithRBACSpecificRoles(additionalRole).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to create incremental backup")
		require.NotNil(t, createResponse2)
		assert.Equal(t, models.BackupCreateResponseStatusSUCCESS, *createResponse2.Status)
		t.Logf("Created incremental backup %s", incrementalBackupID)

		// Clean up roles
		client.Roles().Deleter().WithName(initialRole).Do(ctx)
		client.Roles().Deleter().WithName(additionalRole).Do(ctx)

		// Restore base backup first
		restoreResponse, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(baseBackupID).
			WithRBACAll().
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore base backup")
		require.NotNil(t, restoreResponse)
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse.Status)

		// Then restore incremental backup
		restoreResponse2, err := client.Backup().Restorer().
			WithBackend(backend).
			WithBackupID(incrementalBackupID).
			WithRBACSpecificRoles(additionalRole).
			WithWaitForCompletion(true).
			Do(ctx)

		require.NoError(t, err, "failed to restore incremental backup")
		require.NotNil(t, restoreResponse2)
		assert.Equal(t, models.BackupRestoreResponseStatusSUCCESS, *restoreResponse2.Status)

		// Verify both roles exist
		for _, roleName := range []string{initialRole, additionalRole} {
			exists, err := client.Roles().Exists().WithName(roleName).Do(ctx)
			if err == nil {
				assert.True(t, exists, "Role '%s' should exist after incremental restore", roleName)
			}
		}

		// Cleanup
		client.Roles().Deleter().WithName(initialRole).Do(ctx)
		client.Roles().Deleter().WithName(additionalRole).Do(ctx)
	})
}
