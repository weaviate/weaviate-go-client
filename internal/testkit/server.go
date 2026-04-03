package testkit

import (
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// SchemeHostPort is a helper to obtain schema, hostname, and port,
// as expected by the client internals from a test server.
func SchemeHostPort(t *testing.T, srv *httptest.Server) (schema string, host string, port int) {
	t.Helper()

	url, err := url.Parse(srv.URL)
	require.NoError(t, err, "parse test server url")

	schema, host = url.Scheme, url.Hostname()
	port, err = strconv.Atoi(url.Port())
	require.NoError(t, err, "convert server port")

	return
}
