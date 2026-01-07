package transport_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

func TestBaseEndoint(t *testing.T) {
	var endpoint transport.BaseEndpoint

	assert.Nil(t, endpoint.Query(), "query")
	assert.Nil(t, endpoint.Body(), "body")
}

func TestStaticEndpoint(t *testing.T) {
	static := transport.StaticEndpoint(http.MethodGet, "/live")

	assert.Equal(t, static.Method(), http.MethodGet, "method")
	assert.Equal(t, static.Path(), "/live", "path")
	assert.Nil(t, static.Query(), "query")
	assert.Nil(t, static.Body(), "body")
}

func TestIdentityEndpoint(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		id := "test-id"
		pathFmt := "/string/%s"
		wantPath := fmt.Sprintf(pathFmt, id)

		factory := transport.IdentityEndpoint[string](http.MethodGet, pathFmt)
		req := factory(id)

		if assert.Implements(t, (*transport.Endpoint)(nil), req, "factory returns valid requests") {
			ep := req.(transport.Endpoint)

			assert.Equal(t, ep.Method(), http.MethodGet, "method")
			assert.Equal(t, ep.Path(), wantPath, "path")
			assert.Nil(t, ep.Query(), "query")
			assert.Nil(t, ep.Body(), "body")
		}
	})

	t.Run("int", func(t *testing.T) {
		id := 123
		pathFmt := "/int/%d"
		wantPath := fmt.Sprintf(pathFmt, id)

		factory := transport.IdentityEndpoint[int](http.MethodGet, pathFmt)
		req := factory(id)

		if assert.Implements(t, (*transport.Endpoint)(nil), req, "factory returns valid requests") {
			ep := req.(transport.Endpoint)

			assert.Equal(t, ep.Method(), http.MethodGet, "method")
			assert.Equal(t, ep.Path(), wantPath, "path")
			assert.Nil(t, ep.Query(), "query")
			assert.Nil(t, ep.Body(), "body")
		}
	})

	t.Run("uuid.UUID", func(t *testing.T) {
		id := uuid.New()
		pathFmt := "/uuid/%s"
		wantPath := fmt.Sprintf(pathFmt, id)

		factory := transport.IdentityEndpoint[uuid.UUID](http.MethodGet, pathFmt)
		req := factory(id)

		if assert.Implements(t, (*transport.Endpoint)(nil), req, "factory returns valid requests") {
			ep := req.(transport.Endpoint)

			assert.Equal(t, ep.Method(), http.MethodGet, "method")
			assert.Equal(t, ep.Path(), wantPath, "path")
			assert.Nil(t, ep.Query(), "query")
			assert.Nil(t, ep.Body(), "body")
		}
	})

	t.Run("invalid pathFmt", func(t *testing.T) {
		pathFmt := "/first/%v/second/%d"
		require.Panics(t, func() {
			transport.IdentityEndpoint[any](http.MethodGet, pathFmt)
		}, "must validate pathFmt on creation (%q has %d formatting directives)",
			pathFmt, strings.Count(pathFmt, "%"),
		)
	})
}
