package data_test

import (
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
			transport := testkit.NewTransport(t, []testkit.Stub[api.GetObjectRequest, api.ResourceExistsResponse]{
				{
					Request:  &api.GetObjectRequest{UUID: id, RequestDefaults: rd},
					Response: api.ResourceExistsResponse(tt.exists),
					Err:      tt.err,
				},
			})
			c := data.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			exists, err := c.Exists(t.Context(), id)

			if tt.err == nil {
				assert.NoError(t, err, "returned error")
			} else {
				assert.ErrorIs(t, err, tt.err, "returned error")
			}
			assert.Equal(t, tt.exists, exists, "object exists")
		})
	}
}
