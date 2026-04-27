package api_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestObjectReference(t *testing.T) {
	for _, tt := range []struct {
		ref  *api.ObjectReference
		want string
	}{
		{
			ref:  &api.ObjectReference{UUID: testkit.UUID},
			want: "weaviate://localhost/" + testkit.UUID.String(),
		},
		{
			ref:  &api.ObjectReference{Collection: "Songs", UUID: testkit.UUID},
			want: "weaviate://localhost/Songs/" + testkit.UUID.String(),
		},
	} {
		t.Run(tt.want, func(t *testing.T) {
			beacon, err := json.Marshal(tt.ref)
			require.NoError(t, err, "marshal text")

			got, err := strconv.Unquote(string(beacon))
			require.NoError(t, err, "unquote beacon")
			require.Equal(t, tt.want, got, "beacon")
		})
	}
}
