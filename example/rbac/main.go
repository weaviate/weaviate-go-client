package main

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/example"
	"github.com/weaviate/weaviate-go-client/v6/rbac"
)

const (
	EnvHost   = "WEAVIATE_HOST"
	EnvAPIKey = "WEAVIATE_API_KEY"
)

func main() {
	ctx := context.Background()
	host, apiKey := example.ConnectionParams()

	// Connect to a WCD cluster.
	c, err := weaviate.NewWeaviateCloud(ctx, host, apiKey)
	example.Catch(err)
	defer c.Close()

	// All emails are stored in the Emails collection.
	// Each user has a dedicated tenant (not created here),
	// and can read/write the emails.
	// Admins are an external OIDC group; an admin can modify
	// the Emails collection.
	exists, err := c.Collections.Exists(ctx, "Emails")
	example.Catch(err)
	if exists {
		example.Catch(c.Collections.Delete(ctx, "Emails"))
	}

	_, err = c.Collections.Create(ctx, collections.Collection{
		Name: "Emails",
		Properties: []collections.Property{
			{Name: "title", DataType: collections.DataTypeText},
			{Name: "body", DataType: collections.DataTypeText},
		},
	})
	defer c.Collections.Delete(ctx, "Emails")
	example.Catch(err)

	// Create admin role and assign it to "mailbox-admins" group.
	err = c.Roles.Create(ctx, rbac.Role{
		ID: "custom-admin-role",
		Permissions: rbac.Permissions{
			Collections: []rbac.CollectionPermission{
				{Collection: "Emails", Read: true, Update: true},
			},
		},
	})
	defer c.Roles.Delete(ctx, "custom-admin-role")
	example.Catch(err)

	err = c.Groups.AssignRoles(ctx, rbac.AssignRolesOptions{
		ID:    "mailbox-admins",
		Roles: []string{"custom-admin-role"},
	})
	example.Catch(err)

	// Create user accounts and obtain API keys for them.
	apiKeys := map[string]string{
		"alex":  "",
		"becca": "",
		"cindy": "",
	}

	for username := range apiKeys {
		apiKey, err = c.Users.DB.Create(ctx, username)
		defer c.Users.DB.Delete(ctx, username)
		example.Catch(err)

		apiKeys[username] = apiKey

		err = c.Roles.Create(ctx, rbac.Role{
			ID: "user-" + username,
			Permissions: rbac.Permissions{
				Data: []rbac.DataPermission{
					{Collection: "Emails", Tenant: username, Create: true, Read: true},
				},
			},
		})
		example.Catch(err)
		defer c.Roles.Delete(ctx, "user-"+username)

		err = c.Users.DB.AssignRoles(ctx, rbac.AssignRolesOptions{
			ID:    username,
			Roles: []string{"user" + username},
		})
	}

	// Verify all roles exist and permissions are assigned correctly
	roles, err := c.Roles.List(ctx)
	example.Catch(err)
	log.Printf("Roles: %+v", roles)

	users, err := c.Users.DB.List(ctx, rbac.ListUsersOptions{})
	example.Catch(err)
	log.Printf("Users: %+v", users)

	ok, err := c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "user-alex",
		Data:   rbac.DataPermission{Collection: "Emails", Tenant: "user-alex", Read: true},
	})
	example.Catch(err)
	log.Printf("user-alex has read_data permission: %t", ok)

	ok, err = c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "custom-admin-role",
		Data:   rbac.DataPermission{Collection: "Emails", Tenant: "user-alex", Read: true},
	})
	example.Catch(err)
	log.Printf("custom-admin-role has read_data permission: %t", ok)

	ok, err = c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "custom-admin-role",
		Collection: rbac.CollectionPermission{
			Collection: "Emails", Update: true,
		},
	})
	example.Catch(err)
	log.Printf("custom-admin-role has update_collections permission for Emails: %t", ok)

	// Revoke role from cindy.
	err = c.Users.DB.RevokeRoles(ctx, rbac.RevokeRolesOptions{
		ID:    "cindy",
		Roles: []string{"user-cindy"},
	})
	example.Catch(err)

	assignments, err := c.Roles.UserAssignments(ctx, "user-cindy")
	log.Printf("Users assigned to user-cindy role: %v", assignments)

	// Add delete_collections to admin role
	err = c.Roles.AddPermissions(ctx, rbac.AddPermissions{
		RoleID: "custom-admin-role",
		Permissions: rbac.Permissions{
			Collections: []rbac.CollectionPermission{
				{Collection: "Emails", Delete: true},
			},
		},
	})
	example.Catch(err)

	ok, err = c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "custom-admin-role",
		Collection: rbac.CollectionPermission{
			Collection: "Emails", Delete: true,
		},
	})
	example.Catch(err)
	log.Printf("custom-admin-role has delete_collections permission for Emails: %t", ok)
}
