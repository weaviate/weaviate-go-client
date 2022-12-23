package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const oidcConfigURL = "/.well-known/openid-configuration"

type Config interface {
	GetAuthClient(con *connection.Connection) (*http.Client, error)
}

type authBase struct{}

func (ab authBase) getIdAndTokenEndpoint(con *connection.Connection) (string, string, error) {
	rest, err := con.RunREST(context.TODO(), oidcConfigURL, http.MethodGet, nil)
	if err != nil {
		return "", "", err
	}

	switch status := rest.StatusCode; status {
	case 404:
		log.Println("The client was configured to use authentication, but weaviate is configured without authentication. Are you sure this is correct?")
		return "", "", nil
	case 200: // status code is ok
	default:
		return "", "", fmt.Errorf("OIDC configuration url "+oidcConfigURL+"returned status code %v", fmt.Sprint(rest.StatusCode))
	}

	cfg := struct {
		Href     string `json:"href"`
		ClientID string `json:"clientId"`
	}{}
	decodeErr := rest.DecodeBodyIntoTarget(&cfg)
	if decodeErr != nil {
		return "", "", decodeErr
	}

	endpoints, err := con.RunRESTExternal(context.TODO(), cfg.Href, http.MethodGet, nil)
	if err != nil {
		return "", "", err
	}
	var resultEndpoints map[string]interface{}
	decodeEndpointErr := endpoints.DecodeBodyIntoTarget(&resultEndpoints)
	if decodeEndpointErr != nil {
		return "", "", err
	}
	return cfg.ClientID, resultEndpoints["token_endpoint"].(string), nil
}

type ClientCredentials struct {
	ClientSecret string
	Scopes       []string
	authBase
}

func (cc ClientCredentials) GetAuthClient(con *connection.Connection) (*http.Client, error) {
	clientId, tokenEndpoint, err := cc.getIdAndTokenEndpoint(con)
	if err != nil {
		return nil, err
	} else if clientId == "" && tokenEndpoint == "" {
		return nil, nil // not configured with authentication
	}

	if cc.Scopes == nil {
		if strings.HasPrefix(tokenEndpoint, "https://login.microsoftonline.com") {
			cc.Scopes = []string{clientId + "/.default"}
		}
	}

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
	clientId, tokenEndpoint, err := ro.getIdAndTokenEndpoint(con)
	if err != nil {
		return nil, err
	} else if clientId == "" && tokenEndpoint == "" {
		return nil, nil // not configured with authentication
	}

	if ro.Scopes == nil || len(ro.Scopes) == 0 {
		ro.Scopes = []string{"offline_access"}
	}

	conf := oauth2.Config{ClientID: clientId, Endpoint: oauth2.Endpoint{TokenURL: tokenEndpoint}}
	token, err := conf.PasswordCredentialsToken(context.TODO(), ro.Username, ro.Password)
	if err != nil {
		return nil, err
	}
	// username + password are not saved by the client, so there is no possibility of refreshing the token with a
	// refresh_token.
	if token.RefreshToken == "" {
		log.Printf("Auth001: Your access token is valid for %v and no refresh token was provided.", token.Expiry.Sub(time.Now()))
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
	ExpiresIn    int
	authBase
}

func (bt BearerToken) GetAuthClient(con *connection.Connection) (*http.Client, error) {
	// we don't need these values, but we can check if weaviate is configured with authentication enabled
	clientId, tokenEndpoint, err := bt.getIdAndTokenEndpoint(con)
	if err == nil && clientId == "" && tokenEndpoint == "" {
		return nil, nil
	}

	// there is no possibility of refreshing the token without a refresh_token.
	if bt.RefreshToken == "" {
		log.Printf("Auth001: Your access token is valid for %v and no refresh token was provided.", time.Now().Add(time.Second*time.Duration(bt.ExpiresIn)))
		return oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: bt.AccessToken})), nil
	}
	conf := oauth2.Config{ClientID: clientId, Endpoint: oauth2.Endpoint{TokenURL: tokenEndpoint}}
	tokenSource := conf.TokenSource(context.TODO(), &oauth2.Token{
		AccessToken: bt.AccessToken, RefreshToken: bt.RefreshToken, Expiry: time.Now().Add(time.Second * time.Duration(bt.ExpiresIn)),
	})

	return oauth2.NewClient(context.TODO(), tokenSource), nil
}
