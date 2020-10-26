package misc

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
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
		return false, err
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
		return false, err
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
func (oidcg *OpenIDConfigGetter) Do(ctx context.Context) (*paragons.OpenIDConfiguration, error) {
	response, err := oidcg.connection.RunREST(ctx, "/.well-known/openid-configuration", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, nil
	}
	if response.StatusCode == 200 {
		var openIDConfig paragons.OpenIDConfiguration
		decodeErr := response.DecodeBodyIntoTarget(&openIDConfig)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return &openIDConfig, nil
	}

	return nil, clienterrors.NewUnexpectedStatusCodeError(response.StatusCode, string(response.Body))
}
