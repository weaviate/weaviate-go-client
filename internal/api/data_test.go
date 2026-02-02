package api_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestObjectReference(t *testing.T) {
	t.Run("rountrip", func(t *testing.T) {
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

				var ref api.ObjectReference
				err = json.Unmarshal(beacon, &ref)
				require.NoError(t, err, "unmarshal beacon")
				require.Equal(t, tt.ref, &ref, "same as input")
			})
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		for _, tt := range []struct {
			name   string
			beacon string
		}{
			{
				name:   "not a beacon",
				beacon: "http://Songs/" + testkit.UUID.String(),
			},
			{
				name:   "too many parts",
				beacon: "weaviate://localhost/Songs/uuid/a/b/c",
			},
			{
				name:   "wrong order",
				beacon: "weaviate://" + testkit.UUID.String() + "/" + "Songs",
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				var ref api.ObjectReference
				err := json.Unmarshal([]byte(tt.beacon), &ref)
				require.Error(t, err)
			})
		}
	})
}
