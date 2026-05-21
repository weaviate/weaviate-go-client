package schema

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/docker"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

// imageTag returns the value of the named env var, falling back to defaultTag
// when the var is unset or empty.
func imageTag(envVar, defaultTag string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	return defaultTag
}

// serverVersion bundles a Weaviate image tag with optional server-side env var
// overrides that modify default behaviour (e.g. DEFAULT_VECTOR_INDEX=flat).
type serverVersion struct {
	name        string
	image       string
	extraEnv    map[string]string
	wantDefault string // expected VectorIndexType when the caller sends none
}

// TestDefaultVectorIndexType_Legacy137_4_InjectsHnsw proves that when the
// client connects to Weaviate 1.37.4 (which does not apply a server-side
// DEFAULT_VECTOR_INDEX), the client injects "hnsw" for every class whose
// VectorIndexType is empty.
func TestDefaultVectorIndexType_Legacy137_4_InjectsHnsw(t *testing.T) {
	ctx := context.Background()

	legacyTag := imageTag("WEAVIATE_VERSION_LEGACY", "1.37.4")
	container, err := docker.StartWeaviateWithEnv(ctx, "semitechnologies/weaviate:"+legacyTag, nil)
	require.NoError(t, err, "start Weaviate legacy container (%s)", legacyTag)
	defer func() {
		require.NoError(t, container.Terminate(ctx))
	}()

	client, err := weaviate.NewClient(weaviate.Config{
		Host:   container.Endpoint(docker.HTTP),
		Scheme: "http",
	})
	require.NoError(t, err)

	t.Run("legacy single-vector: empty type becomes hnsw", func(t *testing.T) {
		className := "LegacyNoType"
		defer client.Schema().ClassDeleter().WithClassName(className).Do(ctx) //nolint:errcheck

		class := &models.Class{Class: className}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, "hnsw", got.VectorIndexType,
			"client must inject 'hnsw' for legacy server when VectorIndexType is empty")
	})

	t.Run("named vector: empty type becomes hnsw", func(t *testing.T) {
		className := "LegacyNamedNoType"
		defer client.Schema().ClassDeleter().WithClassName(className).Do(ctx) //nolint:errcheck

		class := &models.Class{
			Class: className,
			VectorConfig: map[string]models.VectorConfig{
				"main": {
					Vectorizer:      map[string]interface{}{"none": map[string]interface{}{}},
					VectorIndexType: "", // intentionally empty
				},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		require.Contains(t, got.VectorConfig, "main")
		assert.Equal(t, "hnsw", got.VectorConfig["main"].VectorIndexType,
			"client must inject 'hnsw' in named vector config for legacy server")
	})
}

// TestDefaultVectorIndexType_New137_5_AppliesServerDefault proves that when
// the client connects to Weaviate 1.37.5 (which natively supports
// DEFAULT_VECTOR_INDEX), the client does NOT inject "hnsw" — instead the
// server applies whatever its DEFAULT_VECTOR_INDEX env var says (here "flat").
func TestDefaultVectorIndexType_New137_5_AppliesServerDefault(t *testing.T) {
	ctx := context.Background()

	extraEnv := map[string]string{
		"DEFAULT_VECTOR_INDEX": "flat",
	}
	// The new-version image tag is read from WEAVIATE_VERSION (existing
	// repo-wide convention). Default is a locally-built 1.37.5 image that
	// includes the DEFAULT_VECTOR_INDEX feature. The version string reports
	// as "1.37.5-..." which the supportsServerSideDefaultVectorIndexType
	// helper treats as >= 1.37.5.
	newTag := imageTag("WEAVIATE_VERSION", "1.37.5-e0fe0d5.amd64")
	container, err := docker.StartWeaviateWithEnv(ctx, "semitechnologies/weaviate:"+newTag, extraEnv)
	require.NoError(t, err, "start Weaviate new container (%s)", newTag)
	defer func() {
		require.NoError(t, container.Terminate(ctx))
	}()

	client, err := weaviate.NewClient(weaviate.Config{
		Host:   container.Endpoint(docker.HTTP),
		Scheme: "http",
	})
	require.NoError(t, err)

	t.Run("legacy single-vector: server default applies (flat)", func(t *testing.T) {
		className := "NewServerNoType"
		defer client.Schema().ClassDeleter().WithClassName(className).Do(ctx) //nolint:errcheck

		class := &models.Class{Class: className}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, "flat", got.VectorIndexType,
			"server must apply DEFAULT_VECTOR_INDEX=flat; client must NOT override it with 'hnsw'")
	})

	t.Run("named vector: server default applies (flat)", func(t *testing.T) {
		className := "NewServerNamedNoType"
		defer client.Schema().ClassDeleter().WithClassName(className).Do(ctx) //nolint:errcheck

		class := &models.Class{
			Class: className,
			VectorConfig: map[string]models.VectorConfig{
				"main": {
					Vectorizer:      map[string]interface{}{"none": map[string]interface{}{}},
					VectorIndexType: "", // intentionally empty — server should fill in "flat"
				},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		require.Contains(t, got.VectorConfig, "main")
		assert.Equal(t, "flat", got.VectorConfig["main"].VectorIndexType,
			"server must apply DEFAULT_VECTOR_INDEX=flat in named vector config")
	})
}

// TestDefaultVectorIndexType_ExplicitFlat_BothVersions verifies that an
// explicit VectorIndexType:"flat" set by the caller is preserved unchanged on
// both the legacy (1.37.4) and new (1.37.5) servers. The inject helper must
// never overwrite a non-empty user value.
func TestDefaultVectorIndexType_ExplicitFlat_BothVersions(t *testing.T) {
	ctx := context.Background()

	versions := []serverVersion{
		{
			name:  "legacy " + imageTag("WEAVIATE_VERSION_LEGACY", "1.37.4"),
			image: "semitechnologies/weaviate:" + imageTag("WEAVIATE_VERSION_LEGACY", "1.37.4"),
		},
		{
			name:  "new " + imageTag("WEAVIATE_VERSION", "1.37.5-e0fe0d5.amd64") + " with DEFAULT_VECTOR_INDEX=flat",
			image: "semitechnologies/weaviate:" + imageTag("WEAVIATE_VERSION", "1.37.5-e0fe0d5.amd64"),
			extraEnv: map[string]string{
				"DEFAULT_VECTOR_INDEX": "flat",
			},
		},
	}

	for _, sv := range versions {
		sv := sv
		t.Run(sv.name, func(t *testing.T) {
			container, err := docker.StartWeaviateWithEnv(ctx, sv.image, sv.extraEnv)
			require.NoError(t, err, "start container")
			defer func() {
				require.NoError(t, container.Terminate(ctx))
			}()

			client, err := weaviate.NewClient(weaviate.Config{
				Host:   container.Endpoint(docker.HTTP),
				Scheme: "http",
			})
			require.NoError(t, err)

			t.Run("legacy single-vector explicit flat preserved", func(t *testing.T) {
				className := "ExplicitFlatSingle"
				defer client.Schema().ClassDeleter().WithClassName(className).Do(ctx) //nolint:errcheck

				class := &models.Class{
					Class:           className,
					VectorIndexType: "flat",
				}
				require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))

				got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
				require.NoError(t, err)
				assert.Equal(t, "flat", got.VectorIndexType,
					"explicit 'flat' must be preserved on %s", sv.name)
			})

			t.Run("named vector explicit flat preserved", func(t *testing.T) {
				className := "ExplicitFlatNamed"
				defer client.Schema().ClassDeleter().WithClassName(className).Do(ctx) //nolint:errcheck

				class := &models.Class{
					Class: className,
					VectorConfig: map[string]models.VectorConfig{
						"main": {
							Vectorizer:      map[string]interface{}{"none": map[string]interface{}{}},
							VectorIndexType: "flat",
						},
					},
				}
				require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))

				got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
				require.NoError(t, err)
				require.Contains(t, got.VectorConfig, "main")
				assert.Equal(t, "flat", got.VectorConfig["main"].VectorIndexType,
					"explicit 'flat' must be preserved in named vector config on %s", sv.name)
			})
		})
	}
}
