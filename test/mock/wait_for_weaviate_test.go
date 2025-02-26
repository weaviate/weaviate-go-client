package connection

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

func TestWaitForWeaviate(t *testing.T) {
	// Tests the WaitForWeaviate function if a connection can be established
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http"}
	client := weaviate.New(cfg)
	err := client.WaitForWeavaite(60 * time.Second)
	assert.Nil(t, err)
}

func TestWaitForWeaviate_NoConnection(t *testing.T) {
	// Tests the WaitForWeaviate function if no connection can be established
	mux := http.NewServeMux()
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http"}
	client := weaviate.New(cfg)
	err := client.WaitForWeavaite(5 * time.Second)
	assert.NotNil(t, err)
}

func TestWaitForWeaviate_longTimeforResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/.well-known/ready", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 2) // we repeat calls
		w.Write([]byte(`{}`))
	})
	s := httptest.NewServer(mux)
	defer s.Close()

	cfg := weaviate.Config{Host: strings.TrimPrefix(s.URL, "http://"), Scheme: "http"}
	client := weaviate.New(cfg)
	start := time.Now()
	err := client.WaitForWeavaite(5 * time.Second)
	assert.Nil(t, err)
	assert.Less(t, time.Since(start).Seconds(), 2.5) // allow for some overhead
}
