package api_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

func TestResourceExistsResponse(t *testing.T) {
	t.Run("http response", func(t *testing.T) {
		var exists api.ResourceExistsResponse
		assert.False(t, exists.Bool(), "Bool() of zero value")

		if assert.Implements(t, (*transports.StatusAccepter)(nil), exists) {
			acc := (any(exists)).(transports.StatusAccepter)
			assert.True(t, acc.AcceptStatus(http.StatusNotFound), "must accept 404 Not Found")

			// Check all other error status codes, fail fast.
			for status := 300; status < 600; status++ {
				if status == http.StatusNotFound {
					continue
				}
				ok := assert.Falsef(t, acc.AcceptStatus(status), "must not accept HTTP %d", status)
				if !ok {
					break
				}
			}
		}

		assert.Implements(t, (*json.Unmarshaler)(nil), &exists)
		for _, body := range []string{
			`{}`, `[]`, `[{}, {}]`,
			`{"accept": "any", "valid": "json"}`,
		} {
			require.NoError(t,
				json.Unmarshal([]byte(body), &exists),
				"ResourceExistsResponse.UnmarshalJSON must accept any valid JSON",
			)
		}

		assert.True(t, exists.Bool(), "Bool() after UnmarshalJSON")
	})

	t.Run("nil receiver", func(t *testing.T) {
		var exists *api.ResourceExistsResponse
		require.Panics(t, func() {
			exists.UnmarshalJSON(nil) //nolint:errcheck
		})
	})
}
