package users_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/users"
)

func TestUserInfo_UnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name string
		json string
		want users.UserInfo
	}{
		{
			name: "has username, no user_id",
			json: `{"username": "John Doe"}`,
			want: users.UserInfo{UserID: "John Doe"},
		},
		{
			name: "has user_id, no username",
			json: `{"user_id": "john_doe"}`,
			want: users.UserInfo{UserID: "john_doe"},
		},
		{
			name: "has user_id and username",
			json: `{"user_id": "john_doe", "username": "John Doe"}`,
			want: users.UserInfo{UserID: "john_doe"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var got users.UserInfo

			err := json.Unmarshal([]byte(tt.json), &got)

			require.NoError(t, err, "unmarshal user info")
			require.EqualValues(t, tt.want, got)
		})
	}
}
