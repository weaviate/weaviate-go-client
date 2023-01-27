package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/connection"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const oidcConfigURL = "/.well-known/openid-configuration"

type Config interface {
	GetAuthClient(con *connection.Connection) (*http.Client, error)
}

type authBase struct{}

func (ab authBase) getIdAndTokenEndpoint(ctx context.Context, con *connection.Connection) (string, []string, string, error) {
	rest, err := con.RunREST(ctx, oidcConfigURL, http.MethodGet, nil)
	if err != nil {
		return "", []string{}, "", err
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
		return "", []string{}, "", nil
	case 200: // status code is ok
		decodeErr := rest.DecodeBodyIntoTarget(&cfg)
		if decodeErr != nil {
			// Some setups are behind proxies that return some default page - for example a login - for all requests.
			// If the response is not json, we assume that this is the case and try unauthenticated access.
			log.Printf("Auth005: Could not parse Weaviates OIDC configuration, using unauthenticated access. If "+
				"you added an authorization header yourself it will be unaffected. This can happen if weaviate is "+
				"miss-configured or you have a proxy inbetween the client and weaviate. You can test this by visiting %v.",
				oidcConfigURL)

			return "", []string{}, "", nil
		}
	default:
		return "", []string{}, "", fmt.Errorf("OIDC configuration url %s returned status code %v", oidcConfigURL, rest.StatusCode)
	}

	endpoints, err := con.RunRESTExternal(context.TODO(), cfg.Href, http.MethodGet, nil)
	if err != nil {
		return "", []string{}, "", err
	}
	var resultEndpoints map[string]interface{}
	decodeEndpointErr := endpoints.DecodeBodyIntoTarget(&resultEndpoints)
	if decodeEndpointErr != nil {
		return "", []string{}, "", err
	}

	tokenEndpoint, ok := resultEndpoints["token_endpoint"].(string)
	if !ok {
		return "", []string{}, "", errors.New("could not parse token_endpoint from OIDC response")
	}

	return cfg.ClientID, cfg.Scopes, tokenEndpoint, nil
}

type ClientCredentials struct {
	ClientSecret string
	Scopes       []string
	authBase
}

func (cc ClientCredentials) GetAuthClient(con *connection.Connection) (*http.Client, error) {
	clientId, weaviateScopes, tokenEndpoint, err := cc.getIdAndTokenEndpoint(context.Background(), con)
	if err != nil {
		return nil, err
	} else if clientId == "" && tokenEndpoint == "" {
		return nil, nil // not configured with authentication
	}

	// remove openid scopes from the scopes returned by weaviate (these are returned by default). These are not accepted
	// by some providers for client credentials
	for j := len(weaviateScopes) - 1; j >= 0; j-- {
		if weaviateScopes[j] == "openid" || weaviateScopes[j] == "email" {
			if j != len(weaviateScopes) {
				weaviateScopes[j] = weaviateScopes[len(weaviateScopes)-1]
			}
			weaviateScopes = weaviateScopes[:len(weaviateScopes)-1]
		}
	}

	if cc.Scopes == nil {
		if strings.HasPrefix(tokenEndpoint, "https://login.microsoftonline.com") {
			cc.Scopes = []string{clientId + "/.default"}
		}
	}
	cc.Scopes = append(cc.Scopes, weaviateScopes...)

	config := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: cc.ClientSecret,
		TokenURL:     tokenEndpoint,
		Scopes:       cc.Scopes,
	}
	return config.Client(context.TODO()), nil
}

type ResourceOwnerPasswordFlow struct {
	Username string
	Password string
	Scopes   []string
	authBase
}

func (ro ResourceOwnerPasswordFlow) GetAuthClient(con *connection.Connection) (*http.Client, error) {
	clientId, weaviateScopes, tokenEndpoint, err := ro.getIdAndTokenEndpoint(context.Background(), con)
	if err != nil {
		return nil, err
	} else if clientId == "" && tokenEndpoint == "" {
		return nil, nil // not configured with authentication
	}

	if ro.Scopes == nil || len(ro.Scopes) == 0 {
		ro.Scopes = []string{"offline_access"}
	}
	ro.Scopes = append(ro.Scopes, weaviateScopes...)

	conf := oauth2.Config{ClientID: clientId, Endpoint: oauth2.Endpoint{TokenURL: tokenEndpoint}}
	token, err := conf.PasswordCredentialsToken(context.TODO(), ro.Username, ro.Password)
	if err != nil {
		return nil, err
	}
	// username + password are not saved by the client, so there is no possibility of refreshing the token with a
	// refresh_token.
	if token.RefreshToken == "" {
		log.Printf("Auth002: Your access token is valid for %v and no refresh token was provided.", token.Expiry.Sub(time.Now()))
		return oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(token)), nil
	}

	// creat oauth configuration that includes the endpoint and client id as a token source with a refresh token
	// (if available), then the client can auto refresh the token
	confRefresh := oauth2.Config{ClientID: clientId, Endpoint: oauth2.Endpoint{TokenURL: tokenEndpoint}}
	tokenSource := confRefresh.TokenSource(context.TODO(), &oauth2.Token{
		AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, Expiry: token.Expiry,
	})

	return oauth2.NewClient(context.TODO(), tokenSource), nil
}

type BearerToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    uint
	authBase
}

func (bt BearerToken) GetAuthClient(con *connection.Connection) (*http.Client, error) {
	// we don't need these values, but we can check if weaviate is configured with authentication enabled
	clientId, weaviateScopes, tokenEndpoint, err := bt.getIdAndTokenEndpoint(context.Background(), con)
	if err == nil && clientId == "" && len(weaviateScopes) == 0 && tokenEndpoint == "" {
		return nil, nil
	}

	// there is no possibility of refreshing the token without a refresh_token.
	if bt.RefreshToken == "" {
		log.Printf("Auth002: Your access token is valid for %v and no refresh token was provided.", time.Now().Add(time.Second*time.Duration(bt.ExpiresIn)))
		return oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: bt.AccessToken})), nil
	}
	conf := oauth2.Config{ClientID: clientId, Endpoint: oauth2.Endpoint{TokenURL: tokenEndpoint}}
	tokenSource := conf.TokenSource(context.TODO(), &oauth2.Token{
		AccessToken: bt.AccessToken, RefreshToken: bt.RefreshToken, Expiry: time.Now().Add(time.Second * time.Duration(bt.ExpiresIn)),
	})

	return oauth2.NewClient(context.TODO(), tokenSource), nil
}
