package except

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

func TestNewWeaviateClientError(t *testing.T) {
	tests := []struct {
		name        string
		weaviateErr error
		want        string
	}{
		{
			name:        "error",
			weaviateErr: NewWeaviateClientErrorf(404, "problem with: %s %s %s", "connection", "to", "localhost"),
			want:        "status code: 404, error: problem with: connection to localhost",
		},
		{
			name:        "derived error",
			weaviateErr: NewDerivedWeaviateClientError(fmt.Errorf("connection closed")),
			want:        "status code: -1, error: check the DerivedFromError field for more information: connection closed",
		},
		{
			name: "unexpected error",
			weaviateErr: NewUnexpectedStatusCodeErrorFromRESTResponse(&connection.ResponseData{
				Body: []byte("body of an error"),
			}),
			want: "status code: 0, error: body of an error",
		},
		{
			name:        "check some error",
			weaviateErr: CheckResponseDataErrorAndStatusCode(nil, fmt.Errorf("some error")),
			want:        "status code: -1, error: check the DerivedFromError field for more information: some error",
		},
		{
			name: "check response error",
			weaviateErr: CheckResponseDataErrorAndStatusCode(&connection.ResponseData{
				StatusCode: 503,
				Body:       []byte("service unavailable"),
			}, nil),
			want: "status code: 503, error: service unavailable",
		},
		{
			name: "check response error that is expected",
			weaviateErr: CheckResponseDataErrorAndStatusCode(&connection.ResponseData{
				StatusCode: 200,
				Body:       []byte("OK"),
			}, nil, 200),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want == "" {
				assert.Nil(t, tt.weaviateErr)
			} else {
				assert.EqualError(t, tt.weaviateErr, tt.want)
			}
		})
	}
}
