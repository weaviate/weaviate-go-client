package auth

import (
	"context"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// RefreshToken implements [transport.Exchanger] for the Refresh Token Grant.
type RefreshToken oauth2.Token

func (rt RefreshToken) Exchange(ctx context.Context, cfg oauth2.Config) (oauth2.TokenSource, error) {
	if rt.Expiry.IsZero() {
		rt.Expiry = time.Now().Add(time.Duration(rt.ExpiresIn) * time.Second)
	}
	return cfg.TokenSource(context.Background(), (*oauth2.Token)(&rt)), nil
}

// ClientCredentials implements [transport.Exchanger] for the Client Credentials Grant.
type ClientCredentials struct {
	ClientSecret string
	Scopes       []string
}

func (cc ClientCredentials) Exchange(ctx context.Context, cfg oauth2.Config) (oauth2.TokenSource, error) {
	ccc := clientcredentials.Config{
		TokenURL:     cfg.Endpoint.TokenURL,
		ClientID:     cfg.ClientID,
		ClientSecret: cc.ClientSecret,
		Scopes:       withDefaultScopes(cfg, cc.Scopes),
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	return ccc.TokenSource(context.Background()), nil
}

// RsourceOwnerPasswordCredentials implements [transport.Exchanger] for the Resource Owner Password Credentials Grant.
//
// The ROPC grant is considered a legacy pattern. We support it for parity with other clients.
type ResourceOwnerPasswordCredentials struct {
	Username     string
	Password     string
	ClientSecret string
	Scopes       []string
}

func (ropc ResourceOwnerPasswordCredentials) Exchange(ctx context.Context, cfg oauth2.Config) (oauth2.TokenSource, error) {
	cfg.ClientSecret = ropc.ClientSecret
	cfg.Scopes = withDefaultScopes(cfg, ropc.Scopes)
	t, err := cfg.PasswordCredentialsToken(ctx, ropc.Username, ropc.Password)
	if err != nil {
		return nil, err
	}
	return cfg.TokenSource(context.Background(), t), nil
}

const (
	defaultScope = "offline_access"
	microsoftURL = "https://login.microsoftonline.com"
)

func withDefaultScopes(conf oauth2.Config, scopes []string) (s []string) {
	s = append(s, defaultScope)
	s = append(s, scopes...)
	s = append(s, conf.Scopes...)

	if strings.HasPrefix(conf.Endpoint.TokenURL, microsoftURL) {
		s = append(s, conf.ClientID+"/.default")
	}
	return
}
