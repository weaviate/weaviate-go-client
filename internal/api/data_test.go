package api_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestReference_MarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		ref  api.Reference
		want string
	}{
		{
			ref:  api.Reference{Target: api.ObjectPath{UUID: testkit.UUID}},
			want: "weaviate://localhost/" + testkit.UUID.String(),
		},
		{
			ref:  api.Reference{Target: api.ObjectPath{Collection: "Songs", UUID: testkit.UUID}},
			want: "weaviate://localhost/Songs/" + testkit.UUID.String(),
		},
	} {
		t.Run(tt.want, func(t *testing.T) {
			want, err := json.Marshal(map[string]string{
				"beacon": tt.want,
			})
			require.NoError(t, err)

			got, err := json.Marshal(tt.ref)
			require.NoError(t, err)

			require.Equal(t, want, got, "marshaled reference")
		})
	}
}
