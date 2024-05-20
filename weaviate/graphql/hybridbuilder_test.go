package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHybridBuilder_build(t *testing.T) {
	t.Run("all parameters", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}
		str := hybrid.WithQuery("query").WithVector([]float32{1, 2, 3}).WithAlpha(0.6).WithProperties([]string{"prop1", "prop2"}).build()
		expected := `hybrid:{query: "query", vector: [1,2,3], alpha: 0.6, properties: ["prop1","prop2"]}`
		require.Equal(t, expected, str)
	})

	t.Run("only query", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}
		str := hybrid.WithQuery("query").build()
		expected := `hybrid:{query: "query"}`
		require.Equal(t, expected, str)
	})

	t.Run("query and vector", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}
		str := hybrid.WithQuery("query").WithVector([]float32{1, 2, 3}).build()
		expected := `hybrid:{query: "query", vector: [1,2,3]}`
		require.Equal(t, expected, str)
	})

	t.Run("query and alpha", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}
		str := hybrid.WithQuery("query").WithAlpha(0.6).build()
		expected := `hybrid:{query: "query", alpha: 0.6}`
		require.Equal(t, expected, str)
	})

	t.Run("query with escaping and alpha", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}

		str := hybrid.WithQuery("\"I'm a complex string\" says the string").WithAlpha(0.6).build()
		expected := `hybrid:{query: "\"I'm a complex string\" says the string", alpha: 0.6}`
		require.Equal(t, expected, str)
	})

	t.Run("query with fusion type Ranked", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}

		str := hybrid.WithQuery("some query").WithFusionType(Ranked).build()
		expected := `hybrid:{query: "some query", fusionType: rankedFusion}`
		require.Equal(t, expected, str)
	})

	t.Run("query with fusion type Relative Score", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}

		str := hybrid.WithQuery("some query").WithFusionType(RelativeScore).build()
		expected := `hybrid:{query: "some query", fusionType: relativeScoreFusion}`
		require.Equal(t, expected, str)
	})

	t.Run("query and alpha and targetVectors", func(t *testing.T) {
		hybrid := HybridArgumentBuilder{}
		str := hybrid.WithQuery("query").WithAlpha(0.6).WithTargetVectors("t1").build()
		expected := `hybrid:{query: "query", alpha: 0.6, targetVectors: ["t1"]}`
		require.Equal(t, expected, str)
	})

	t.Run("query and nearText search", func(t *testing.T) {
		neartText := &NearTextArgumentBuilder{}
		neartText.WithConcepts([]string{"concept"}).WithCertainty(0.9)
		searches := &HybridSearchesArgumentBuilder{}
		searches.WithNearText(neartText)
		hybrid := HybridArgumentBuilder{}
		str := hybrid.WithQuery("I'm a simple string").WithSearches(searches).build()
		expected := `hybrid:{query: "I'm a simple string", searches:{nearText:{concepts: ["concept"] certainty: 0.9}}}`
		require.Equal(t, expected, str)
	})

	t.Run("query and nearVector search", func(t *testing.T) {
		neartVector := &NearVectorArgumentBuilder{}
		neartVector.WithVector([]float32{0.1, 0.2, 0.3}).WithCertainty(0.9)
		searches := &HybridSearchesArgumentBuilder{}
		searches.WithNearVector(neartVector)
		hybrid := HybridArgumentBuilder{}
		str := hybrid.WithQuery("I'm a simple string").WithSearches(searches).build()
		expected := `hybrid:{query: "I'm a simple string", searches:{nearVector:{certainty: 0.9 vector: [0.1,0.2,0.3]}}}`
		require.Equal(t, expected, str)
	})
}
