package testkit

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

// SchemeHostPort is a helper to obtain schema, hostname, and port,
// as expected by the client internals from a test server.
func SchemeHostPort(t *testing.T, srv *httptest.Server) (schema string, host string, port string) {
	t.Helper()

	url, err := url.Parse(srv.URL)
	require.NoError(t, err, "parse test server url")

	return url.Scheme, url.Hostname(), url.Port()
}
