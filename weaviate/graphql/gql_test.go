package graphql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
)

type fakeREST struct {
	response *connection.ResponseData
}

func (f fakeREST) RunREST(context.Context, string, string, interface{}) (*connection.ResponseData, error) {
	return f.response, nil
}

func TestRunGraphQLQueryPointsToV6WhenGraphQLIsDisabled(t *testing.T) {
	disabled := fakeREST{&connection.ResponseData{
		StatusCode: 422,
		Body:       []byte(`{"error":[{"message":"graphql api is disabled"}]}`),
	}}

	_, err := runGraphQLQuery(context.Background(), disabled, "{Get{Product{sku}}}")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status code: 422")
	assert.Contains(t, err.Error(), "graphql api is disabled")
	assert.Contains(t, err.Error(), "github.com/weaviate/weaviate-go-client/v6")
}

func TestRunGraphQLQueryKeepsOtherErrors(t *testing.T) {
	failing := fakeREST{&connection.ResponseData{StatusCode: 500, Body: []byte("boom")}}

	_, err := runGraphQLQuery(context.Background(), failing, "{Get{Product{sku}}}")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	assert.NotContains(t, err.Error(), "weaviate-go-client/v6")
}
