package users_test

import (
	"encoding/json"
	"testing"
	"time"

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
		{
			name: "is active",
			json: `{"active": true, "username": "John Doe"}`,
			want: users.UserInfo{Active: true, UserID: "John Doe"},
		},
		{
			name: "is deactivated",
			json: `{"active": false, "username": "John Doe"}`,
			want: users.UserInfo{Active: false, UserID: "John Doe"},
		},
		{
			name: "was created at specific time with timezone",
			json: `{"createdAt": "2025-04-08T09:33:00+02:00", "username": "John Doe"}`,
			want: users.UserInfo{CreatedAt: time.Date(2025, time.April, 8, 9, 33, 0, 0,
				time.FixedZone("", 2*60*60)), UserID: "John Doe"},
		},
		{
			name: "was created at specific time without timezone",
			json: `{"createdAt": "2025-04-08T09:33:00+00:00", "username": "John Doe"}`,
			want: users.UserInfo{CreatedAt: time.Date(2025, time.April, 8, 9, 33, 0, 0,
				time.UTC), UserID: "John Doe"},
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
