package api_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/go-openapi/testify/v2/require"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

func TestObjectReference(t *testing.T) {
	t.Run("rountrip", func(t *testing.T) {
		for _, tt := range []struct {
			ref  *api.ObjectReference
			want string
		}{
			{
				ref:  &api.ObjectReference{UUID: uuid.Nil},
				want: "weaviate://localhost/" + uuid.Nil.String(),
			},
			{
				ref:  &api.ObjectReference{Collection: "Songs", UUID: uuid.Nil},
				want: "weaviate://localhost/Songs/" + uuid.Nil.String(),
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
	})

	t.Run("unmarshal", func(t *testing.T) {
		for _, tt := range []struct {
			name   string
			beacon string
		}{
			{
				name:   "not a beacon",
				beacon: "http://Songs/" + uuid.Nil.String(),
			},
			{
				name:   "too many parts",
				beacon: "weaviate://localhost/Songs/uuid/a/b/c",
			},
			{
				name:   "wrong order",
				beacon: "weaviate://" + uuid.Nil.String() + "/" + "Songs",
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
