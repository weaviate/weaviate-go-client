package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

// Test how api.ResourceExistsResponse unmarshals supported HTTP and gRPC reponses.
func TestResourceExistsResponse(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		var exists api.ResourceExistsResponse
		assert.False(t, exists.Bool(), "Bool() of zero value")

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

		assert.Implements(t, (*json.Unmarshaler)(nil), &exists)
		for _, body := range []string{
			`{}`, `[]`, `[{}, {}]`,
			`{"accept": "any", "valid": "json"}`,
		} {
			require.NoError(t,
				json.Unmarshal([]byte(body), &exists),
				"headResponse.UnmarshalJSON must accept any input",
			)
		}

		assert.True(t, exists.Bool(), "Bool() after UnmarshalJSON")
	})

	t.Run("grpc search", func(t *testing.T) {
		var zero api.ResourceExistsResponse
		assert.False(t, zero.Bool(), "Bool() of zero value")

		require.Implements(t, (*transport.MessageUnmarshaler[proto.SearchReply])(nil), &zero)

		for _, tt := range []struct {
			reply *proto.SearchReply // Cannot be nil, but we mustn't copy the lock in proto.SearchReply.
			want  bool
		}{
			{reply: &proto.SearchReply{}, want: false},
			{reply: &proto.SearchReply{Results: make([]*proto.SearchResult, 1)}, want: true},
			{reply: &proto.SearchReply{Results: make([]*proto.SearchResult, 5)}, want: true},
		} {
			t.Run(fmt.Sprintf("reply with %d results", len(tt.reply.Results)), func(t *testing.T) {
				var exists api.ResourceExistsResponse

				u := (any(&exists)).(transport.MessageUnmarshaler[proto.SearchReply])
				err := u.UnmarshalMessage(tt.reply)
				require.NoError(t, err, "unmashal reply message")

				require.Equal(t, tt.want, exists.Bool(), "Bool() after UnmarshalMessage")
			})
		}
	})
}

// Test that we have an assertions guarding against nil input.
func TestResourceExistsResponse_UnmarshalMessage(t *testing.T) {
	var exists api.ResourceExistsResponse
	require.PanicsWithValue(t, "search reply is nil", func() {
		exists.UnmarshalMessage(nil)
	})
}
