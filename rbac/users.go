package rbac

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type UsersClient struct {
	transport internal.Transport
	DB        *DBUsersClient
	OIDC      *OIDCUsersClient
}

func NewUsersClient(t internal.Transport) *UsersClient {
	return &UsersClient{
		transport: t,
		DB: &DBUsersClient{
			kindClient: &kindClient{
				transport: t,
				kind:      api.RBACKindDB,
			},
		},
		OIDC: &OIDCUsersClient{
			kindClient: &kindClient{
				transport: t,
				kind:      api.RBACKindOIDC,
			},
		},
	}
}

type MyUser struct {
	ID     string
	Roles  []Role
	Groups []string
}

func (c *UsersClient) MyUser(ctx context.Context) (*MyUser, error) {
	return nil, nil
}

type kindClient struct {
	transport internal.Transport
	kind      string
}

type AssignedRoles struct {
	UserID             string
	IncludePermissions bool
}

func (c *kindClient) AssignedRoles(ctx context.Context, options AssignedRoles) ([]Role, error) {
	return nil, nil
}

func (c *kindClient) AssignRoles(ctx context.Context, userID string, roles ...string) error {
	return nil
}

func (c *kindClient) RevokeRoles(ctx context.Context, userID string, roles ...string) error {
	return nil
}

type (
	DBUsersClient   struct{ *kindClient }
	OIDCUsersClient struct{ *kindClient }
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

type DeactivateOptions struct {
	ID        string
	RevokeKey bool
}

func (c *DBUsersClient) Deactivate(ctx context.Context, options DeactivateOptions) error {
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

type ByNameOptions struct {
	ID                string
	IncludeLastUsedAt bool
}

func (c *DBUsersClient) ByName(ctx context.Context, options ByNameOptions) (*UserInfo, error) {
	req := &api.GetUserInfoRequest{
		ID:                options.ID,
		IncludeLastUsedAt: options.IncludeLastUsedAt,
	}

	var resp api.UserInfo
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("get db user info: %w", err)
	}
	return &UserInfo{
		ID:     resp.ID,
		Type:   resp.Type,
		Active: resp.Active,
		Roles:  resp.Roles,
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
			ID:     ui.ID,
			Type:   ui.Type,
			Active: ui.Active,
			Roles:  ui.Roles,
		}
	}
	return out, nil
}
