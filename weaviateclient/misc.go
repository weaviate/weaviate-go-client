package weaviateclient

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
	"net/http"
)


// Misc collection of endpoints that don't fit in other categories
type Misc struct {
	connection *connection.Connection
}


// ReadyChecker retrieves weaviate ready status
func (misc *Misc) ReadyChecker() *readyChecker {
	return &readyChecker{connection: misc.connection}
}

// LiveChecker retrieves weaviate live status
func (misc *Misc) LiveChecker() *liveChecker {
	return &liveChecker{connection: misc.connection}
}

// OpenIDConfigurationGetter retrieves the Open ID configuration
// may be nil
func (misc *Misc) OpenIDConfigurationGetter() *openIDConfigGetter {
	return &openIDConfigGetter{connection: misc.connection}
}

type readyChecker struct {
	connection *connection.Connection
}

// Do the ready request
func (rc *readyChecker) Do(ctx context.Context) (bool, error) {
	response, err := rc.connection.RunREST(ctx, "/.well-known/ready", http.MethodGet)
	if err != nil {
		return false, err
	}
	if response.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

type liveChecker struct {
	connection *connection.Connection
}

// Do the liveChecker request
func (lc *liveChecker) Do(ctx context.Context) (bool, error) {
	response, err := lc.connection.RunREST(ctx, "/.well-known/live", http.MethodGet)
	if err != nil {
		return false, err
	}
	if response.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

type openIDConfigGetter struct {
	connection *connection.Connection
}

// Do the open ID config request
func (oidcg *openIDConfigGetter) Do(ctx context.Context) (*models.OpenIDConfiguration, error) {
	response, err := oidcg.connection.RunREST(ctx, "/.well-known/openid-configuration", http.MethodGet)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, nil
	}
	if response.StatusCode == 200 {
		var openIDConfig models.OpenIDConfiguration
		decodeErr := response.DecodeBodyIntoTarget(&openIDConfig)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return &openIDConfig, nil
	}

	return nil, NewUnexpectedStatusCodeError(response.StatusCode, string(response.Body))
}
