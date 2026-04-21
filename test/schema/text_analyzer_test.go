package schema

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

const minTextAnalyzerVersion = "1.37.0"

// TestTextAnalyzer_integration covers Weaviate 1.37.0 schema features ported
// from weaviate-python-client PR #2006:
//   - Property.TextAnalyzer (asciiFold, asciiFoldIgnore, stopwordPreset)
//   - InvertedIndexConfig.StopwordPresets
//
// All tests are skipped when the server is older than 1.37.0.
func TestTextAnalyzer_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		require.NoError(t, err, "failed to setup weaviate")
	})

	client := testsuit.CreateTestClient(false)
	testsuit.AtLeastWeaviateVersion(t, client, minTextAnalyzerVersion,
		"text analyzer config requires Weaviate >= "+minTextAnalyzerVersion)

	ctx := context.Background()

	t.Run("StopwordPresets_AppliedAndRoundTripped", func(t *testing.T) {
		className := "TestStopwordPresets1"
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		class := &models.Class{
			Class: className,
			InvertedIndexConfig: &models.InvertedIndexConfig{
				StopwordPresets: map[string][]string{
					"fr": {"le", "la", "les"},
				},
			},
			Properties: []*models.Property{
				{Name: "title", DataType: []string{"text"}},
			},
		}

		err := client.Schema().ClassCreator().WithClass(class).Do(ctx)
		require.NoError(t, err)
		t.Cleanup(func() { _ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx) })

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		require.NotNil(t, got.InvertedIndexConfig)
		assert.Equal(t,
			map[string][]string{"fr": {"le", "la", "les"}},
			got.InvertedIndexConfig.StopwordPresets,
		)
	})

	t.Run("StopwordPresets_Update_ReplacesPreset", func(t *testing.T) {
		className := "TestStopwordPresetsUpdate"
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		class := &models.Class{
			Class: className,
			InvertedIndexConfig: &models.InvertedIndexConfig{
				StopwordPresets: map[string][]string{"fr": {"le"}},
			},
			Properties: []*models.Property{
				{Name: "title", DataType: []string{"text"}},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))
		t.Cleanup(func() { _ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx) })

		// Update: replace the preset map entirely.
		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		got.InvertedIndexConfig.StopwordPresets = map[string][]string{
			"fr": {"le", "la", "les", "un", "une"},
		}
		require.NoError(t, client.Schema().ClassUpdater().WithClass(got).Do(ctx))

		got, err = client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		assert.Equal(t,
			map[string][]string{"fr": {"le", "la", "les", "un", "une"}},
			got.InvertedIndexConfig.StopwordPresets,
		)
	})

	t.Run("StopwordPresets_RemoveInUse_RejectedByServer", func(t *testing.T) {
		className := "TestStopwordPresetsRemoveInUse"
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		class := &models.Class{
			Class: className,
			InvertedIndexConfig: &models.InvertedIndexConfig{
				StopwordPresets: map[string][]string{"fr": {"le"}},
			},
			Properties: []*models.Property{
				{
					Name:         "title",
					DataType:     []string{"text"},
					TextAnalyzer: &models.TextAnalyzerConfig{StopwordPreset: "fr"},
				},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))
		t.Cleanup(func() { _ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx) })

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		got.InvertedIndexConfig.StopwordPresets = map[string][]string{}

		err = client.Schema().ClassUpdater().WithClass(got).Do(ctx)
		assert.Error(t, err, "expected server to reject removing an in-use stopword preset")
	})

	t.Run("TextAnalyzer_CombinedASCIIFoldAndStopwordPreset", func(t *testing.T) {
		className := "TestTextAnalyzerCombined"
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		asciiFold := true
		class := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "title",
					DataType: []string{"text"},
					TextAnalyzer: &models.TextAnalyzerConfig{
						ASCIIFold:      asciiFold,
						StopwordPreset: "en",
					},
				},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))
		t.Cleanup(func() { _ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx) })

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		require.Len(t, got.Properties, 1)
		require.NotNil(t, got.Properties[0].TextAnalyzer)
		assert.True(t, got.Properties[0].TextAnalyzer.ASCIIFold)
		assert.Equal(t, "en", got.Properties[0].TextAnalyzer.StopwordPreset)
	})

	t.Run("TextAnalyzer_ASCIIFoldIgnore_RoundTrips", func(t *testing.T) {
		className := "TestTextAnalyzerASCIIFoldIgnore"
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		class := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "title",
					DataType: []string{"text"},
					TextAnalyzer: &models.TextAnalyzerConfig{
						ASCIIFold:       true,
						ASCIIFoldIgnore: []string{"é", "ñ"},
					},
				},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))
		t.Cleanup(func() { _ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx) })

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		require.Len(t, got.Properties, 1)
		require.NotNil(t, got.Properties[0].TextAnalyzer)
		assert.True(t, got.Properties[0].TextAnalyzer.ASCIIFold)
		assert.ElementsMatch(t, []string{"é", "ñ"}, got.Properties[0].TextAnalyzer.ASCIIFoldIgnore)
	})

	t.Run("NestedProperty_TextAnalyzer_RoundTrips", func(t *testing.T) {
		className := "TestNestedPropertyTextAnalyzer"
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		class := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "meta",
					DataType: []string{"object"},
					NestedProperties: []*models.NestedProperty{
						{
							Name:         "headline",
							DataType:     []string{"text"},
							TextAnalyzer: &models.TextAnalyzerConfig{ASCIIFold: true},
						},
					},
				},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))
		t.Cleanup(func() { _ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx) })

		got, err := client.Schema().ClassGetter().WithClassName(className).Do(ctx)
		require.NoError(t, err)
		require.Len(t, got.Properties, 1)
		require.Len(t, got.Properties[0].NestedProperties, 1)
		require.NotNil(t, got.Properties[0].NestedProperties[0].TextAnalyzer)
		assert.True(t, got.Properties[0].NestedProperties[0].TextAnalyzer.ASCIIFold)
	})
}
