package main

import (
	"context"
	"log"
	"os"

	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/rbac"
)

const (
	EnvHost   = "WEAVIATE_HOST"
	EnvAPIKey = "WEAVIATE_API_KEY"
)

func main() {
	host, hasHost := os.LookupEnv(EnvHost)
	apiKey, hasKey := os.LookupEnv(EnvAPIKey)
	if !hasHost || !hasKey {
		log.Printf("%q and %q must be defined. Skipping example.", EnvHost, EnvAPIKey)
		return
	}

	ctx := context.Background()

	// Connect to a WCD cluster.
	c, err := weaviate.NewWeaviateCloud(ctx, host, apiKey)
	catch(err)

	// All emails are stored in the Emails collection.
	// Each user has a dedicated tenant (not created here),
	// and can read/write the emails.
	// Admins are an external OIDC group; an admin can modify
	// the Emails collection.
	_, err = c.Collections.Create(ctx, collections.Collection{
		Name: "Emails",
		Properties: []collections.Property{
			{Name: "title", DataType: collections.DataTypeText},
			{Name: "body", DataType: collections.DataTypeText},
		},
	})
	catch(err)
	defer c.Collections.Delete(ctx, "Emails")

	// Create admin role and assign it to "mailbox-admins" group.
	err = c.Roles.Create(ctx, rbac.Role{
		ID: "admin",
		Permissions: rbac.Permissions{
			Collections: []rbac.CollectionPermission{
				{Collection: "Emails", Read: true, Update: true},
			},
		},
	})
	catch(err)
	defer c.Roles.Delete(ctx, "admin")

	err = c.Groups.AssignRoles(ctx, rbac.AssignRolesOptions{
		ID:    "mailbox-admins",
		Roles: []string{"admin"},
	})
	catch(err)

	// Create user accounts and obtain API keys for them.
	apiKeys := map[string]string{
		"alex":  "",
		"becca": "",
		"cindy": "",
	}

	for username := range apiKeys {
		apiKey, err = c.Users.DB.Create(ctx, username)
		catch(err)
		defer c.Users.DB.Delete(ctx, "username")

		apiKeys[username] = apiKey

		err = c.Roles.Create(ctx, rbac.Role{
			ID: "user-" + username,
			Permissions: rbac.Permissions{
				Data: []rbac.DataPermission{
					{Collection: "Emails", Tenant: username, Create: true, Read: true},
				},
			},
		})
		catch(err)
		defer c.Roles.Delete(ctx, "user-"+username)

		err = c.Users.DB.AssignRoles(ctx, rbac.AssignRolesOptions{
			ID:    username,
			Roles: []string{"user" + username},
		})
	}

	// Verify all roles exist and permissions are assigned correctly
	roles, err := c.Roles.List(ctx)
	catch(err)
	log.Printf("Roles: %v", roles)

	users, err := c.Users.DB.List(ctx, rbac.ListUsersOptions{})
	catch(err)
	log.Printf("Users: %v", users)

	ok, err := c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "user-alex",
		Data:   rbac.DataPermission{Collection: "Emails", Tenant: "user-alex", Read: true},
	})
	catch(err)
	log.Printf("user-alex has read_data permission: %t", ok)

	ok, err = c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "admin",
		Data:   rbac.DataPermission{Collection: "Emails", Tenant: "user-alex", Read: true},
	})
	catch(err)
	log.Printf("admin has read_data permission: %t", ok)

	ok, err = c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "admin",
		Collection: rbac.CollectionPermission{
			Collection: "Emails", Update: true,
		},
	})
	catch(err)
	log.Printf("admin has update_collections permission for Emails: %t", ok)

	// Revoke role from cindy.
	err = c.Users.DB.RevokeRoles(ctx, rbac.RevokeRolesOptions{
		ID:    "cindy",
		Roles: []string{"user-cindy"},
	})
	catch(err)

	assignments, err := c.Roles.UserAssignments(ctx, "user-cindy")
	log.Printf("Users assigned to user-cindy role: %v", assignments)

	// Add delete_collections to admin role
	err = c.Roles.AddPermissions(ctx, rbac.AddPermissions{
		RoleID: "admin",
		Permissions: rbac.Permissions{
			Collections: []rbac.CollectionPermission{
				{Collection: "Emails", Delete: true},
			},
		},
	})
	catch(err)

	ok, err = c.Roles.HasPermission(ctx, rbac.HasPermission{
		RoleID: "admin",
		Collection: rbac.CollectionPermission{
			Collection: "Emails", Delete: true,
		},
	})
	catch(err)
	log.Printf("admin has delete_collections permission for Emails: %t", ok)
}

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(err)
}
