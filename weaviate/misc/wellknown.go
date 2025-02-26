package misc

import (
	"context"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
)

// OpenIDConfiguration of weaviate
type OpenIDConfiguration struct {
	// The Location to redirect to
	Href string `json:"href,omitempty"`
	// OAuth Client ID
	ClientID string `json:"clientId,omitempty"`
}

// ReadyChecker builder to check if weaviate is ready
type ReadyChecker struct {
	connection        *connection.Connection
	dbVersionProvider *db.VersionProvider
}

// Do the ready request
func (rc *ReadyChecker) Do(ctx context.Context) (bool, error) {
	response, err := rc.connection.RunREST(ctx, "/.well-known/ready", http.MethodGet, nil)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == 200 {
		rc.dbVersionProvider.Refresh()
		return true, nil
	}
	return false, nil
}

// LiveChecker builder to check if weaviate is live
type LiveChecker struct {
	connection        *connection.Connection
	dbVersionProvider *db.VersionProvider
}

// Do the LiveChecker request
func (lc *LiveChecker) Do(ctx context.Context) (bool, error) {
	response, err := lc.connection.RunREST(ctx, "/.well-known/live", http.MethodGet, nil)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == 200 {
		lc.dbVersionProvider.Refresh()
		return true, nil
	}
	return false, nil
}

// OpenIDConfigGetter builder to retrieve the openID configuration
type OpenIDConfigGetter struct {
	connection *connection.Connection
}

// Do the open ID config request
func (oidcg *OpenIDConfigGetter) Do(ctx context.Context) (*OpenIDConfiguration, error) {
	response, err := oidcg.connection.RunREST(ctx, "/.well-known/openid-configuration", http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == 404 {
		return nil, nil
	}
	if response.StatusCode == 200 {
		var openIDConfig OpenIDConfiguration
		decodeErr := response.DecodeBodyIntoTarget(&openIDConfig)
		return &openIDConfig, decodeErr
	}

	return nil, except.NewWeaviateClientError(response.StatusCode, string(response.Body))
}
