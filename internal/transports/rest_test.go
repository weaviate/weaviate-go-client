package transports_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

func TestBaseEndoint(t *testing.T) {
	var endpoint transports.BaseEndpoint

	assert.Nil(t, endpoint.Query(), "query")
	assert.Nil(t, endpoint.Body(), "body")
}
