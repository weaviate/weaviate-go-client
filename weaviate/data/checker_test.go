package data

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

func TestChecker_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		handler        http.HandlerFunc
		closeServer    bool
		expectedExists bool
		expectErr      bool
		expectedStatus int
	}{
		{
			name: "object exists (204)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodHead, r.Method)
				w.WriteHeader(http.StatusNoContent)
			},
			expectedExists: true,
			expectErr:      false,
		},
		{
			name: "object does not exist (404)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodHead, r.Method)
				w.WriteHeader(http.StatusNotFound)
			},
			expectedExists: false,
			expectErr:      false,
		},
		{
			name: "gateway timeout (504)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodHead, r.Method)
				w.WriteHeader(http.StatusGatewayTimeout)
			},
			expectedExists: false,
			expectErr:      true,
			expectedStatus: 504,
		},
		{
			name: "network error (server closed)",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			closeServer:    true,
			expectedExists: false,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tt.handler)
			if tt.closeServer {
				server.Close()
			} else {
				defer server.Close()
			}

			serverURL, err := url.Parse(server.URL)
			require.NoError(t, err)

			con := connection.NewConnection(serverURL.Scheme, serverURL.Host, server.Client(), time.Second, nil)
			checker := &Checker{
				connection:       con,
				id:               "test-id",
				className:        "TestClass",
				dbVersionSupport: newDBVersionSupportForTests("1.14.0"),
			}

			exists, err := checker.Do(context.Background())
			if tt.expectErr {
				assert.Error(t, err)
				assert.False(t, exists)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
			}
		})
	}
}
