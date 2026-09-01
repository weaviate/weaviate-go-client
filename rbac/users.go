package rbac

import (
	"context"
	"fmt"
	"slices"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type UsersClient struct {
	transport internal.Transport
	DB        *DBUsersClient
	OIDC      *OIDCUsersClient
}

func NewUsersClient(t internal.Transport) *UsersClient {
	dev.AssertNotNil(t, "transport")

	return &UsersClient{
		transport: t,
		DB: &DBUsersClient{
			kindClient: newKindClient[api.UserID](t, api.RBACKindDB),
		},
		OIDC: &OIDCUsersClient{
			kindClient: newKindClient[api.UserID](t, api.RBACKindOIDC),
		},
	}
}

type MyUserInfo struct {
	ID     string
	Roles  []Role
	Groups []string
}

func (c *UsersClient) MyUserInfo(ctx context.Context) (*MyUserInfo, error) {
	var resp api.GetOwnUserInfoResponse
	if err := c.transport.Do(ctx, api.GetOwnUserInfoRequest, &resp); err != nil {
		return nil, fmt.Errorf("get my user info: %w", err)
	}

	roles := slices.Grow([]Role(nil), len(resp.Roles))
	for i := range resp.Roles {
		roles = append(roles, unmarshalRole(&resp.Roles[i]))
	}

	return &MyUserInfo{
		ID:     resp.ID,
		Groups: resp.Groups,
		Roles:  roles,
	}, nil
}

func newKindClient[ID api.UserID | api.GroupID](t internal.Transport, kind string) *kindClient[ID] {
	return &kindClient[ID]{
		transport: t,
		kind:      kind,
	}
}

type kindClient[ID api.UserID | api.GroupID] struct {
	transport internal.Transport
	kind      string
}

type AssignedRolesOptions struct {
	ID                 string
	IncludePermissions bool
}

func (c *kindClient[ID]) AssignedRoles(ctx context.Context, options AssignedRolesOptions) ([]Role, error) {
	req := &api.GetAssignedRolesRequest{
		Kind:               c.kind,
		Entity:             api.RBACEntity(ID(options.ID)),
		IncludePermissions: options.IncludePermissions,
	}

	var resp []api.Role
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("get assigned roles: %w", err)
	}
	roles := slices.Grow([]Role(nil), len(resp))
	for i := range resp {
		roles = append(roles, unmarshalRole(&resp[i]))
	}
	return roles, nil
}

type AssignRolesOptions struct {
	ID    string
	Roles []string
}

func (c *kindClient[ID]) AssignRoles(ctx context.Context, options AssignRolesOptions) error {
	req := &api.ManageRolesRequest{
		Kind:   c.kind,
		Entity: api.RBACEntity(ID(options.ID)),
		Verb:   api.RoleVerbAssign,
		Roles:  options.Roles,
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("assign roles: %w", err)
	}
	return nil
}

type RevokeRolesOptions struct {
	ID    string
	Roles []string
}

func (c *kindClient[ID]) RevokeRoles(ctx context.Context, options RevokeRolesOptions) error {
	req := &api.ManageRolesRequest{
		Kind:   c.kind,
		Entity: api.RBACEntity(ID(options.ID)),
		Verb:   api.RoleVerbRevoke,
		Roles:  options.Roles,
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("revoke roles: %w", err)
	}
	return nil
}

type (
	DBUsersClient   struct{ *kindClient[api.UserID] }
	OIDCUsersClient struct{ *kindClient[api.UserID] }
)

func (c *DBUsersClient) Create(ctx context.Context, userID string) (string, error) {
	var resp api.CreateUserResponse
	if err := c.transport.Do(ctx, api.CreateUserRequest(userID), &resp); err != nil {
		return "", fmt.Errorf("create db user: %w", err)
	}
	return resp.APIKey, nil
}

func (c *DBUsersClient) Delete(ctx context.Context, userID string) error {
	if err := c.transport.Do(ctx, api.DeleteUserRequest(userID), nil); err != nil {
		return fmt.Errorf("delete db user: %w", err)
	}
	return nil
}

func (c *DBUsersClient) Activate(ctx context.Context, userID string) error {
	if err := c.transport.Do(ctx, api.ActivateUserRequest(userID), nil); err != nil {
		return fmt.Errorf("activate db user: %w", err)
	}
	return nil
}

type DeactivateUserOptions struct {
	ID        string
	RevokeKey bool
}

func (c *DBUsersClient) Deactivate(ctx context.Context, options DeactivateUserOptions) error {
	req := &api.DeactivateUserRequest{
		ID:        options.ID,
		RevokeKey: options.RevokeKey,
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("deactivate db user: %w", err)
	}
	return nil
}

func (c *DBUsersClient) RotateKey(ctx context.Context, userID string) (string, error) {
	var resp api.RotateAPIKeyResponse
	if err := c.transport.Do(ctx, api.RotateAPIKeyRequest(userID), &resp); err != nil {
		return "", fmt.Errorf("rotate db user api key: %w", err)
	}
	return resp.APIKey, nil
}

const (
	UserTypeDB    = api.UserTypeDB    // User created via REST API.
	UserTypeDBEnv = api.UserTypeDBEnv // User defined in the server environment.
	UserTypeOIDC  = api.UserTypeOIDC  // User managed by an OIDC provider.

	GroupTypeOIDC = api.GroupTypeOIDC // Group managed by an OIDC provider.
)

type UserInfoOptions struct {
	ID                string
	IncludeLastUsedAt bool
}

func (c *DBUsersClient) UserInfo(ctx context.Context, options UserInfoOptions) (*UserInfo, error) {
	req := &api.GetUserInfoRequest{
		ID:                options.ID,
		IncludeLastUsedAt: options.IncludeLastUsedAt,
	}

	var resp api.UserInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("get db user info: %w", err)
	}
	return &UserInfo{
		ID:                 resp.ID,
		Type:               resp.Type,
		Active:             resp.Active,
		Roles:              resp.Roles,
		Namespace:          resp.Namespace,
		APIKeyFirstLetters: resp.APIKeyFirstLetters,
		CreatedAt:          resp.CreatedAt,
		LastUsedAt:         resp.LastUsedAt,
	}, nil
}

type ListUsersOptions struct {
	IncludeLastUsedAt bool
}

func (c *DBUsersClient) List(ctx context.Context, options ListUsersOptions) ([]UserInfo, error) {
	req := &api.ListUsersRequest{
		IncludeLastUsedAt: options.IncludeLastUsedAt,
	}

	var resp []api.UserInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("list db users: %w", err)
	}

	out := make([]UserInfo, len(resp))
	for i, ui := range resp {
		out[i] = UserInfo{
			ID:                 ui.ID,
			Type:               ui.Type,
			Active:             ui.Active,
			Roles:              ui.Roles,
			Namespace:          ui.Namespace,
			APIKeyFirstLetters: ui.APIKeyFirstLetters,
			CreatedAt:          ui.CreatedAt,
			LastUsedAt:         ui.LastUsedAt,
		}
	}
	return out, nil
}
