package data_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

// Test data.Client.Exists:
//   - Passes the right UUID and request defaults
//   - Uses api.ResourceExistsResponse as dest
//   - Propagates any error from the transport
func TestDataClient_Exists(t *testing.T) {
	id := uuid.New()
	rd := api.RequestDefaults{
		CollectionName:   "Exists",
		ConsistencyLevel: api.ConsistencyLevelOne,
		Tenant:           "john_doe",
	}

	for _, tt := range []struct {
		exists bool
		err    error
	}{
		{false, nil},
		{true, nil},
		{false, errors.New("whaaam!")},
	} {
		t.Run(fmt.Sprintf("exists=%t error=%v", tt.exists, tt.err), func(t *testing.T) {
			var called bool
			transport := testkit.TransportFunc(func(ctx context.Context, req, dest any) error {
				called = true
				assert.Equal(t, ctx, t.Context(), "bad contenxt")

				if assert.IsType(t, (*api.GetObjectRequest)(nil), req, "bad request") {
					get := req.(*api.GetObjectRequest)
					assert.Equal(t, id, get.UUID, "object uuid")
					assert.Equal(t, rd, get.RequestDefaults, "request defaults")
				}

				// Imitate ResourceExistsResponse.UnmarshalMessage.
				if assert.IsType(t, (*api.ResourceExistsResponse)(nil), dest, "bad response") {
					exists := (dest).(*api.ResourceExistsResponse)
					*exists = api.ResourceExistsResponse(tt.exists)
				}

				return tt.err
			})

			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "NewClient returned nil client")

			exists, err := c.Exists(t.Context(), id)
			require.True(t, called, "must call transport.Do")

			if tt.err == nil {
				assert.NoError(t, err, "returned error")
			} else {
				assert.ErrorIs(t, err, tt.err, "returned error")
			}
			assert.Equal(t, tt.exists, exists, "return value")
		})
	}
}
