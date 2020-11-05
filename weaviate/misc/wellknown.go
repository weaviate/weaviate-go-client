package misc

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/models"
	"net/http"
)

// ReadyChecker builder to check if weaviate is ready
type ReadyChecker struct {
	connection *connection.Connection
}

// Do the ready request
func (rc *ReadyChecker) Do(ctx context.Context) (bool, error) {
	response, err := rc.connection.RunREST(ctx, "/.well-known/ready", http.MethodGet, nil)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

// LiveChecker builder to check if weaviate is live
type LiveChecker struct {
	connection *connection.Connection
}

// Do the LiveChecker request
func (lc *LiveChecker) Do(ctx context.Context) (bool, error) {
	response, err := lc.connection.RunREST(ctx, "/.well-known/live", http.MethodGet, nil)
	if err != nil {
		return false, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

// OpenIDConfigGetter builder to retrieve the openID configuration
type OpenIDConfigGetter struct {
	connection *connection.Connection
}

// Do the open ID config request
func (oidcg *OpenIDConfigGetter) Do(ctx context.Context) (*models.OpenIDConfiguration, error) {
	response, err := oidcg.connection.RunREST(ctx, "/.well-known/openid-configuration", http.MethodGet, nil)
	if err != nil {
		return nil, except.NewDerivedWeaviateClientError(err)
	}
	if response.StatusCode == 404 {
		return nil, nil
	}
	if response.StatusCode == 200 {
		var openIDConfig models.OpenIDConfiguration
		decodeErr := response.DecodeBodyIntoTarget(&openIDConfig)
		return &openIDConfig, decodeErr
	}

	return nil, except.NewWeaviateClientError(response.StatusCode, string(response.Body))
}
