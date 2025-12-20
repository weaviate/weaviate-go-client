package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

func TestResourceExistsResponse(t *testing.T) {
	var exists ResourceExistsResponse

	assert.False(t, exists.Bool(), "Bool() return on zero value")
	assert.Implements(t, (*json.Unmarshaler)(nil), &exists)

	if assert.Implements(t, (*transport.StatusAccepter)(nil), exists) {
		acc := (any(exists)).(transport.StatusAccepter)
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

	for _, body := range []string{
		`{}`, `[]`, `[{}, {}]`,
		`{"accept": "any", "valid": "json"}`,
	} {
		require.NoError(t,
			json.Unmarshal([]byte(body), &exists),
			"headResponse.UnmarshalJSON must accept any input",
		)
	}

	assert.True(t, exists.Bool(), "Bool() return after UnmarshalJSON")
}
