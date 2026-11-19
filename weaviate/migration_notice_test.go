package weaviate_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
)

func TestNewClientMigrationNotice(t *testing.T) {
	// A subprocess gives the process-wide notice fresh state independently of
	// other tests constructing clients. It does not reset production state.
	if os.Getenv("WEAVIATE_TEST_MIGRATION_NOTICE") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestNewClientMigrationNotice$")
		cmd.Env = append(os.Environ(), "WEAVIATE_TEST_MIGRATION_NOTICE=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("notice subprocess: %v\n%s", err, output)
		}
		return
	}

	var output bytes.Buffer
	log.SetOutput(&output)
	log.SetFlags(0)

	_, err := weaviate.NewClient(weaviate.Config{
		AuthConfig:       auth.ApiKey{Value: "unused"},
		ConnectionClient: &http.Client{},
	})
	if err == nil || err.Error() != "only AuthConfig or ConnectionClient can be given in the config" {
		t.Fatalf("constructor error changed: %v", err)
	}
	if output.Len() != 0 {
		t.Fatalf("failed constructor logged a migration notice: %s", output.String())
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.39.5"}`))
	}))
	defer server.Close()
	endpoint, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		client, err := weaviate.NewClient(weaviate.Config{
			Host: endpoint.Host, Scheme: endpoint.Scheme, ConnectionClient: server.Client(),
		})
		if err != nil || client == nil {
			t.Fatalf("construct client: %v", err)
		}
	}

	notice := output.String()
	if count := strings.Count(notice, "We strongly recommend upgrading the whole application"); count != 1 {
		t.Fatalf("expected one notice after two successful constructors, got %d: %s", count, notice)
	}
	for _, text := range []string{
		"Weaviate's GraphQL API is deprecated",
		"v5 client's GraphQL() queries",
		"github.com/weaviate/weaviate-go-client/v6",
		"ingestion and collection-management code",
		"MIGRATING_TO_V6.md",
		"https://github.com/weaviate/weaviate-go-client/blob/ax/202611-v5-to-v6-upgrade-nudges/MIGRATING_TO_V6.md",
	} {
		if !strings.Contains(notice, text) {
			t.Errorf("notice missing %q: %s", text, notice)
		}
	}
}
