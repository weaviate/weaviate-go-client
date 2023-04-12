package fault

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeaviateClientError_Error(t *testing.T) {
	tests := []struct {
		name        string
		weaviateErr *WeaviateClientError
		want        string
	}{
		{
			name: "Just message",
			weaviateErr: &WeaviateClientError{
				StatusCode: 404,
				Msg:        "page not found",
			},
			want: "status code: 404, error: page not found",
		},
		{
			name: "Derived error",
			weaviateErr: &WeaviateClientError{
				IsUnexpectedStatusCode: false,
				StatusCode:             -1,
				Msg:                    "Derived error",
				DerivedFromError:       fmt.Errorf("connection not ready"),
			},
			want: "status code: -1, error: Derived error: connection not ready",
		},
		{
			name: "Derived error in derived",
			weaviateErr: &WeaviateClientError{
				IsUnexpectedStatusCode: false,
				StatusCode:             -1,
				Msg:                    "Derived error",
				DerivedFromError: &WeaviateClientError{
					IsUnexpectedStatusCode: false,
					StatusCode:             -2,
					Msg:                    "Derived error2",
					DerivedFromError:       fmt.Errorf("connection not ready"),
				},
			},
			want: "status code: -1, error: Derived error: status code: -2, error: Derived error2: connection not ready",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualError(t, tt.weaviateErr, tt.want)
		})
	}
}
