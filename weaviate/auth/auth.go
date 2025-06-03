package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const oidcConfigURL = "/.well-known/openid-configuration"

type Config interface {
	GetAuthInfo(con *connection.Connection) (*http.Client, map[string]string, error)
	// Returns a nil pointer if the authentication method does not use an api key.
	ApiKey() *string
}

type authBase struct {
	ClientId       string
	WeaviateScopes []string
	TokenEndpoint  string
	Config
}

func (ab *authBase) getIdAndTokenEndpoint(ctx context.Context, con *connection.Connection) error {
	rest, err := con.RunREST(ctx, oidcConfigURL, http.MethodGet, nil)
	if err != nil {
		return err
	}
	cfg := struct {
		Href     string   `json:"href"`
		ClientID string   `json:"clientId"`
		Scopes   []string `json:"scopes"`
	}{}

	switch status := rest.StatusCode; status {
	case 404:
		log.Println("Auth001: The client was configured to use authentication, but weaviate is configured without" +
			"authentication. Are you sure this is correct?")
		return nil
	case 200: // status code is ok
		decodeErr := rest.DecodeBodyIntoTarget(&cfg)
		if decodeErr != nil {
			// Some setups are behind proxies that return some default page - for example a login - for all requests.
			// If the response is not json, we assume that this is the case and try unauthenticated access.
			log.Printf("Auth005: Could not parse Weaviates OIDC configuration, using unauthenticated access. If "+
				"you added an authorization header yourself it will be unaffected. This can happen if weaviate is "+
				"miss-configured or you have a proxy in between the client and weaviate. You can test this by visiting %v.",
				oidcConfigURL)

			return nil
		}
	default:
		return fmt.Errorf("OIDC configuration url %s returned status code %v", oidcConfigURL, rest.StatusCode)
	}

	endpoints, err := con.RunRESTExternal(context.TODO(), cfg.Href, http.MethodGet, nil)
	if err != nil {
		return err
	}
	var resultEndpoints map[string]interface{}
	decodeEndpointErr := endpoints.DecodeBodyIntoTarget(&resultEndpoints)
	if decodeEndpointErr != nil {
		return err
	}

	tokenEndpoint, ok := resultEndpoints["token_endpoint"].(string)
	if !ok {
		return errors.New("could not parse token_endpoint from OIDC response")
	}
	ab.ClientId = cfg.ClientID
	ab.WeaviateScopes = cfg.Scopes
	ab.TokenEndpoint = tokenEndpoint
	return nil
}

type ClientCredentials struct {
	ClientSecret string
	Scopes       []string
	authBase
}

func (cc ClientCredentials) GetAuthInfo(con *connection.Connection) (*http.Client, map[string]string, error) {
	err := cc.getIdAndTokenEndpoint(context.Background(), con)
	if err != nil {
		return nil, nil, err
	} else if cc.ClientId == "" && cc.TokenEndpoint == "" {
		return nil, nil, nil // not configured with authentication
	}

	if cc.Scopes == nil {
		if strings.HasPrefix(cc.TokenEndpoint, "https://login.microsoftonline.com") {
			cc.Scopes = []string{cc.ClientId + "/.default"}
		}
	}
	cc.Scopes = append(cc.Scopes, cc.WeaviateScopes...)

	config := clientcredentials.Config{
		ClientID:     cc.ClientId,
		ClientSecret: cc.ClientSecret,
		TokenURL:     cc.TokenEndpoint,
		Scopes:       cc.Scopes,
	}
	return config.Client(context.TODO()), nil, nil
}

func (cc ClientCredentials) ApiKey() *string {
	return nil
}

type ResourceOwnerPasswordFlow struct {
	Username string
	Password string
	Scopes   []string
	authBase
}

func (ro ResourceOwnerPasswordFlow) GetAuthInfo(con *connection.Connection) (*http.Client, map[string]string, error) {
	err := ro.getIdAndTokenEndpoint(context.Background(), con)
	if err != nil {
		return nil, nil, err
	} else if ro.ClientId == "" && ro.TokenEndpoint == "" {
		return nil, nil, nil // not configured with authentication
	}

	if len(ro.Scopes) == 0 {
		ro.Scopes = []string{"offline_access"}
	}
	ro.Scopes = append(ro.Scopes, ro.WeaviateScopes...)

	conf := oauth2.Config{ClientID: ro.ClientId, Endpoint: oauth2.Endpoint{TokenURL: ro.TokenEndpoint}}
	token, err := conf.PasswordCredentialsToken(context.TODO(), ro.Username, ro.Password)
	if err != nil {
		return nil, nil, err
	}
	// username + password are not saved by the client, so there is no possibility of refreshing the token with a
	// refresh_token.
	if token.RefreshToken == "" {
		log.Printf("Auth002: Your access token is valid for %v and no refresh token was provided.", time.Until(token.Expiry))
		return oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(token)), nil, nil
	}

	// creat oauth configuration that includes the endpoint and client id as a token source with a refresh token
	// (if available), then the client can auto refresh the token
	confRefresh := oauth2.Config{ClientID: ro.ClientId, Endpoint: oauth2.Endpoint{TokenURL: ro.TokenEndpoint}}
	tokenSource := confRefresh.TokenSource(context.TODO(), &oauth2.Token{
		AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, Expiry: token.Expiry,
	})

	return oauth2.NewClient(context.TODO(), tokenSource), nil, nil
}

func (ro ResourceOwnerPasswordFlow) ApiKey() *string {
	return nil
}

type BearerToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    uint
	authBase
}

func (bt BearerToken) GetAuthInfo(con *connection.Connection) (*http.Client, map[string]string, error) {
	// we don't need these values, but we can check if weaviate is configured with authentication enabled
	err := bt.getIdAndTokenEndpoint(context.Background(), con)
	if err == nil && bt.ClientId == "" && len(bt.WeaviateScopes) == 0 && bt.TokenEndpoint == "" {
		return nil, nil, nil
	}

	// there is no possibility of refreshing the token without a refresh_token.
	if bt.RefreshToken == "" {
		log.Printf("Auth002: Your access token is valid for %v and no refresh token was provided.", time.Now().Add(time.Second*time.Duration(bt.ExpiresIn)))
		return oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: bt.AccessToken})), nil, nil
	}
	conf := oauth2.Config{ClientID: bt.ClientId, Endpoint: oauth2.Endpoint{TokenURL: bt.TokenEndpoint}}
	tokenSource := conf.TokenSource(context.TODO(), &oauth2.Token{
		AccessToken: bt.AccessToken, RefreshToken: bt.RefreshToken, Expiry: time.Now().Add(time.Second * time.Duration(bt.ExpiresIn)),
	})

	return oauth2.NewClient(context.TODO(), tokenSource), nil, nil
}

func (bt BearerToken) ApiKey() *string {
	return nil
}

type ApiKey struct {
	Value string
}

// Returns the header used for the authentification.
func (api ApiKey) GetAuthInfo(con *connection.Connection) (*http.Client, map[string]string, error) {
	additional_headers := make(map[string]string)
	additional_headers["authorization"] = "Bearer " + api.Value
	return nil, additional_headers, nil
}

func (api ApiKey) ApiKey() *string {
	return &api.Value
}
