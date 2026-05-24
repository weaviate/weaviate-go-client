package tokenize

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/tokenize"
	"github.com/weaviate/weaviate/entities/models"
)

const minTokenizeVersion = "1.37.0"

func ptrBool(b bool) *bool { return &b }

func TestTokenize_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		require.NoError(t, err, "failed to setup weaviate")
	})

	client := testsuit.CreateTestClient(false)
	testsuit.AtLeastWeaviateVersion(t, client, minTokenizeVersion,
		"tokenize endpoint requires Weaviate >= "+minTokenizeVersion)

	ctx := context.Background()

	// -------- Serialization ------------------------------------------------

	t.Run("Tokenization_Enum", func(t *testing.T) {
		// Weaviate >=1.37.4 filters default stopwords on the query side for `word`;
		// older versions return query == indexed.
		queryStopwordsFiltered := testsuit.ServerAtLeast(t, client, "1.37.4")
		cases := []struct {
			name            string
			method          tokenize.Tokenization
			text            string
			expectedIndexed []string
			expectedQuery   []string // nil = same as expectedIndexed
		}{
			{"word", tokenize.Word, "The quick brown fox", []string{"the", "quick", "brown", "fox"}, []string{"quick", "brown", "fox"}},
			{"lowercase", tokenize.Lowercase, "Hello World Test", []string{"hello", "world", "test"}, nil},
			{"whitespace", tokenize.Whitespace, "Hello World Test", []string{"Hello", "World", "Test"}, nil},
			{"field", tokenize.Field, "  Hello World  ", []string{"Hello World"}, nil},
			{"trigram", tokenize.Trigram, "Hello", []string{"hel", "ell", "llo"}, nil},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := client.Tokenize().Text().
					WithText(tc.text).
					WithTokenization(tc.method).
					Do(ctx)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.method, result.Tokenization)
				assert.Equal(t, tc.expectedIndexed, result.Indexed)
				wantQuery := tc.expectedIndexed
				if queryStopwordsFiltered && tc.expectedQuery != nil {
					wantQuery = tc.expectedQuery
				}
				assert.Equal(t, wantQuery, result.Query)
			})
		}
	})

	t.Run("NoAnalyzerConfig", func(t *testing.T) {
		result, err := client.Tokenize().Text().
			WithText("hello world").
			WithTokenization(tokenize.Word).
			Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, tokenize.Word, result.Tokenization)
		assert.Equal(t, []string{"hello", "world"}, result.Indexed)
		assert.Nil(t, result.AnalyzerConfig)
	})

	t.Run("AsciiFold", func(t *testing.T) {
		cfg := &tokenize.AnalyzerConfig{AsciiFold: ptrBool(true)}
		result, err := client.Tokenize().Text().
			WithText("L'école est fermée").
			WithTokenization(tokenize.Word).
			WithAnalyzerConfig(cfg).
			Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"l", "ecole", "est", "fermee"}, result.Indexed)
	})

	t.Run("AsciiFold_WithIgnore", func(t *testing.T) {
		cfg := &tokenize.AnalyzerConfig{
			AsciiFold:       ptrBool(true),
			AsciiFoldIgnore: []string{"é"},
		}
		result, err := client.Tokenize().Text().
			WithText("L'école est fermée").
			WithTokenization(tokenize.Word).
			WithAnalyzerConfig(cfg).
			Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"l", "école", "est", "fermée"}, result.Indexed)
	})

	t.Run("StopwordPreset_String", func(t *testing.T) {
		cfg := &tokenize.AnalyzerConfig{StopwordPreset: "en"}
		result, err := client.Tokenize().Text().
			WithText("The quick brown fox").
			WithTokenization(tokenize.Word).
			WithAnalyzerConfig(cfg).
			Do(ctx)
		require.NoError(t, err)
		assert.NotContains(t, result.Query, "the")
		assert.Contains(t, result.Query, "quick")
	})

	t.Run("Combined_AsciiFold_Stopwords", func(t *testing.T) {
		cfg := &tokenize.AnalyzerConfig{
			AsciiFold:       ptrBool(true),
			AsciiFoldIgnore: []string{"é"},
			StopwordPreset:  "en",
		}
		result, err := client.Tokenize().Text().
			WithText("The école est fermée").
			WithTokenization(tokenize.Word).
			WithAnalyzerConfig(cfg).
			Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"the", "école", "est", "fermée"}, result.Indexed)
		assert.NotContains(t, result.Query, "the")
		assert.Contains(t, result.Query, "école")
	})

	t.Run("CustomPreset_Additions", func(t *testing.T) {
		// 1.37.4 changed `stopwordPresets` to `[]string`; SDK still emits the
		// legacy map. TODO: rework SDK API and re-enable (weaviate/weaviate#11079).
		if testsuit.ServerAtLeast(t, client, "1.37.4") {
			t.Skip("custom stopword presets need an SDK API rework for 1.37.4+ wire shape")
		}
		cfg := &tokenize.AnalyzerConfig{StopwordPreset: "custom"}
		presets := map[string]*tokenize.StopwordConfig{
			"custom": {Additions: []string{"test"}},
		}
		result, err := client.Tokenize().Text().
			WithText("hello world test").
			WithTokenization(tokenize.Word).
			WithAnalyzerConfig(cfg).
			WithStopwordPresets(presets).
			Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"hello", "world", "test"}, result.Indexed)
		assert.Equal(t, []string{"hello", "world"}, result.Query)
	})

	t.Run("CustomPreset_BaseAndRemovals", func(t *testing.T) {
		if testsuit.ServerAtLeast(t, client, "1.37.4") {
			t.Skip("custom stopword presets need an SDK API rework for 1.37.4+ wire shape")
		}
		cfg := &tokenize.AnalyzerConfig{StopwordPreset: "en-no-the"}
		presets := map[string]*tokenize.StopwordConfig{
			"en-no-the": {Preset: "en", Removals: []string{"the"}},
		}
		result, err := client.Tokenize().Text().
			WithText("the quick").
			WithTokenization(tokenize.Word).
			WithAnalyzerConfig(cfg).
			WithStopwordPresets(presets).
			Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, []string{"the", "quick"}, result.Indexed)
		assert.Equal(t, []string{"the", "quick"}, result.Query)
	})

	// -------- Deserialization ---------------------------------------------

	t.Run("Result_Types", func(t *testing.T) {
		result, err := client.Tokenize().Text().
			WithText("hello").
			WithTokenization(tokenize.Word).
			Do(ctx)
		require.NoError(t, err)
		assert.IsType(t, &tokenize.TokenizeResult{}, result)
		assert.IsType(t, []string{}, result.Indexed)
		assert.IsType(t, []string{}, result.Query)
	})

	t.Run("AnalyzerConfig_Echoed", func(t *testing.T) {
		// Behaviour still covered by AsciiFold/StopwordPreset_String etc.
		if testsuit.ServerAtLeast(t, client, "1.37.4") {
			t.Skip("server no longer echoes AnalyzerConfig on response (weaviate/weaviate#11079)")
		}
		cfg := &tokenize.AnalyzerConfig{
			AsciiFold:       ptrBool(true),
			AsciiFoldIgnore: []string{"é"},
			StopwordPreset:  "en",
		}
		result, err := client.Tokenize().Text().
			WithText("L'école").
			WithTokenization(tokenize.Word).
			WithAnalyzerConfig(cfg).
			Do(ctx)
		require.NoError(t, err)
		require.NotNil(t, result.AnalyzerConfig)
		require.NotNil(t, result.AnalyzerConfig.AsciiFold)
		assert.True(t, *result.AnalyzerConfig.AsciiFold)
		assert.Equal(t, []string{"é"}, result.AnalyzerConfig.AsciiFoldIgnore)
		assert.Equal(t, "en", result.AnalyzerConfig.StopwordPreset)
	})

	t.Run("AnalyzerConfig_None", func(t *testing.T) {
		result, err := client.Tokenize().Text().
			WithText("hello").
			WithTokenization(tokenize.Word).
			Do(ctx)
		require.NoError(t, err)
		assert.Nil(t, result.AnalyzerConfig)
	})

	// -------- Property-scoped tokenize -------------------------------------

	t.Run("PropertyTokenize_Field", func(t *testing.T) {
		className := "TestTokenizePropField"
		_ = client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		class := &models.Class{
			Class:      className,
			Vectorizer: "none",
			Properties: []*models.Property{
				{
					Name:         "tag",
					DataType:     []string{"text"},
					Tokenization: "field",
				},
			},
		}
		require.NoError(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))
		defer client.Schema().ClassDeleter().WithClassName(className).Do(ctx)

		result, err := client.Tokenize().Property().
			WithClassName(className).
			WithPropertyName("tag").
			WithText("  Hello World  ").
			Do(ctx)
		require.NoError(t, err)
		assert.Equal(t, tokenize.Field, result.Tokenization)
		assert.Equal(t, []string{"Hello World"}, result.Indexed)
	})
}
