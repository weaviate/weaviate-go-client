package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const OidcConfigURL = "/.well-known/openid-configuration"

type Config interface {
	GetAuthClient(con *connection.Connection) (*http.Client, error)
}

type authBase struct{}

func (ab authBase) getIdAndTokenEndpoint(con *connection.Connection) (string, string, error) {
	rest, err := con.RunREST(context.TODO(), "/.well-known/openid-configuration", http.MethodGet, nil)
	if err != nil {
		return "", "", err
	}
	if rest.StatusCode == 404 {
		log.Println("The client was configured to use authentication, but weaviate is configured without authentication. Are you sure this is correct?")
		return "", "", nil
	} else if rest.StatusCode != 200 {
		return "", "", errors.New("OIDC configuration url " + OidcConfigURL + "returned status code " + fmt.Sprint(rest.StatusCode))
	}
	var result map[string]interface{}
	decodeErr := rest.DecodeBodyIntoTarget(&result)
	if decodeErr != nil {
		return "", "", decodeErr
	}

	clientId := result["clientId"].(string)
	href := result["href"].(string)
	endpoints, err := con.RunRESTExternal(context.TODO(), href, http.MethodGet, nil)
	if err != nil {
		return "", "", err
	}
	var resultEndpoints map[string]interface{}
	decodeEndpointErr := endpoints.DecodeBodyIntoTarget(&resultEndpoints)
	if decodeEndpointErr != nil {
		return "", "", err
	}
	return clientId, resultEndpoints["token_endpoint"].(string), nil
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
	} else if err == nil && clientId == "" && tokenEndpoint == "" {
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
	authBase
}

func (ro ResourceOwnerPasswordFlow) GetAuthClient(con *connection.Connection) (*http.Client, error) {
	clientId, tokenEndpoint, err := ro.getIdAndTokenEndpoint(con)
	if err != nil {
		return nil, err
	} else if err == nil && clientId == "" && tokenEndpoint == "" {
		return nil, nil // not configured with authentication
	}

	conf := oauth2.Config{ClientID: clientId, Endpoint: oauth2.Endpoint{TokenURL: tokenEndpoint}}
	token, err := conf.PasswordCredentialsToken(context.TODO(), ro.Username, ro.Password)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(token)), nil
}

type BearerToken struct {
	Token string
	authBase
}

func (bt BearerToken) GetAuthClient(con *connection.Connection) (*http.Client, error) {
	// we don't need these values, but we can check if weaviate is configured with authentication enabled
	clientId, tokenEndpoint, err := bt.getIdAndTokenEndpoint(con)
	if err == nil && clientId == "" && tokenEndpoint == "" {
		return nil, nil
	}

	return oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: bt.Token})), nil
}
