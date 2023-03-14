package connection

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func TestWaitForWeaviate(t *testing.T) {
	// Returns the address of the auth server
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http"}
	client := weaviate.New(cfg)
	err := client.WaitForWeavaite(5)
	assert.Nil(t, err)
}

func TestWaitForWeaviate_NoConnection(t *testing.T) {
	// Returns the address of the auth server
	mux := http.NewServeMux()
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http"}
	client := weaviate.New(cfg)
	err := client.WaitForWeavaite(5)
	assert.NotNil(t, err)
}
