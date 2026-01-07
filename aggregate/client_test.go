package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/aggregate"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestClient_NearVector(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Aggregations",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	for _, tt := range []struct {
		name  string
		nv    aggregate.NearVector
		stubs []testkit.Stub[api.AggregateRequest, api.AggregateResponse]
		want  aggregate.Result
	}{} {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			transport := testkit.NewTransport(t, tt.stubs)
			c := aggregate.NewClient(transport, rd)

			// Act
			res, err := c.NearVector(t.Context(), tt.nv)

			require.NoError(t, err, "near vector error")
			require.Equal(t, tt.want, res, "bad result")
		})
	}
}
