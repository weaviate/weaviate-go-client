package internal

import (
	"testing"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate/entities/models"
)

func newVersionProvider(version string) *db.VersionProvider {
	return db.NewVersionProvider(func() string { return version })
}

func TestSupportsServerSideDefaultVectorIndexType(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{"1.37.4 is below threshold", "1.37.4", false},
		{"1.37.5 is exactly threshold", "1.37.5", true},
		{"1.37.5 with build suffix", "1.37.5-e0fe0d5.amd64", true},
		{"1.38.0 is above minor", "1.38.0", true},
		{"2.0.0 is major bump", "2.0.0", true},
		{"empty string", "", false},
		{"malformed string", "malformed", false},
		{"only two parts", "1.37", false},
		{"1.37.3 is below patch", "1.37.3", false},
		{"1.36.99 is below minor", "1.36.99", false},
		{"0.37.5 is below major", "0.37.5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := supportsServerSideDefaultVectorIndexType(tt.version)
			if got != tt.want {
				t.Errorf("supportsServerSideDefaultVectorIndexType(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

func TestInjectLegacyDefaultVectorIndexType(t *testing.T) {
	t.Run("nil class does not panic", func(t *testing.T) {
		p := newVersionProvider("1.37.4")
		InjectLegacyDefaultVectorIndexType(p, nil)
	})

	t.Run("nil provider does not panic", func(t *testing.T) {
		class := &models.Class{}
		InjectLegacyDefaultVectorIndexType(nil, class)
	})

	t.Run("server 1.37.5 + empty type stays empty", func(t *testing.T) {
		p := newVersionProvider("1.37.5")
		class := &models.Class{VectorIndexType: ""}
		InjectLegacyDefaultVectorIndexType(p, class)
		if class.VectorIndexType != "" {
			t.Errorf("expected empty, got %q", class.VectorIndexType)
		}
	})

	t.Run("server 1.37.5 + hnsw stays hnsw", func(t *testing.T) {
		p := newVersionProvider("1.37.5")
		class := &models.Class{VectorIndexType: "hnsw"}
		InjectLegacyDefaultVectorIndexType(p, class)
		if class.VectorIndexType != "hnsw" {
			t.Errorf("expected hnsw, got %q", class.VectorIndexType)
		}
	})

	t.Run("server 1.37.4 + empty type becomes hnsw", func(t *testing.T) {
		p := newVersionProvider("1.37.4")
		class := &models.Class{VectorIndexType: ""}
		InjectLegacyDefaultVectorIndexType(p, class)
		if class.VectorIndexType != "hnsw" {
			t.Errorf("expected hnsw, got %q", class.VectorIndexType)
		}
	})

	t.Run("server 1.37.4 + flat stays flat", func(t *testing.T) {
		p := newVersionProvider("1.37.4")
		class := &models.Class{VectorIndexType: "flat"}
		InjectLegacyDefaultVectorIndexType(p, class)
		if class.VectorIndexType != "flat" {
			t.Errorf("expected flat, got %q", class.VectorIndexType)
		}
	})

	t.Run("server 1.37.4 + empty VectorConfig entry becomes hnsw", func(t *testing.T) {
		p := newVersionProvider("1.37.4")
		class := &models.Class{
			VectorIndexType: "hnsw",
			VectorConfig: map[string]models.VectorConfig{
				"custom": {VectorIndexType: ""},
			},
		}
		InjectLegacyDefaultVectorIndexType(p, class)
		if class.VectorConfig["custom"].VectorIndexType != "hnsw" {
			t.Errorf("expected hnsw in VectorConfig, got %q", class.VectorConfig["custom"].VectorIndexType)
		}
	})

	t.Run("server 1.37.4 + non-empty VectorConfig entry stays unchanged", func(t *testing.T) {
		p := newVersionProvider("1.37.4")
		class := &models.Class{
			VectorConfig: map[string]models.VectorConfig{
				"custom": {VectorIndexType: "flat"},
			},
		}
		InjectLegacyDefaultVectorIndexType(p, class)
		if class.VectorConfig["custom"].VectorIndexType != "flat" {
			t.Errorf("expected flat in VectorConfig, got %q", class.VectorConfig["custom"].VectorIndexType)
		}
	})

	t.Run("empty server version + empty type becomes hnsw (safe fallback)", func(t *testing.T) {
		p := newVersionProvider("")
		class := &models.Class{VectorIndexType: ""}
		InjectLegacyDefaultVectorIndexType(p, class)
		if class.VectorIndexType != "hnsw" {
			t.Errorf("expected hnsw, got %q", class.VectorIndexType)
		}
	})
}
