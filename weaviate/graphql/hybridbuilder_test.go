package graphql

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHybridBuilder_build(t *testing.T) {
	for _, tt := range []struct {
		name  string
		apply func(h *HybridArgumentBuilder)
		want  string
	}{
		{
			name: "all parameters",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("query").
					WithVector([]float32{1, 2, 3}).
					WithAlpha(0.6).
					WithProperties([]string{"prop1", "prop2"})
			},
			want: `hybrid:{query: "query", vector: [1,2,3], alpha: 0.6, properties: ["prop1","prop2"]}`,
		},
		{
			name: "only query",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("query")
			},
			want: `hybrid:{query: "query"}`,
		},
		{
			name: "query and vector",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("query").
					WithVector([]float32{1, 2, 3})
			},
			want: `hybrid:{query: "query", vector: [1,2,3]}`,
		},
		{
			name: "query and alpha",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("query").
					WithAlpha(0.6)
			},
			want: `hybrid:{query: "query", alpha: 0.6}`,
		},
		{
			name: "query with escaping and alpha",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("\"I'm a complex string\" says the string").
					WithAlpha(0.6)
			},
			want: `hybrid:{query: "\"I'm a complex string\" says the string", alpha: 0.6}`,
		},
		{
			name: "query with fusion type Ranked",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("some query").
					WithFusionType(Ranked)
			},
			want: `hybrid:{query: "some query", fusionType: rankedFusion}`,
		},
		{
			name: "query with fusion type Relative Score",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("some query").
					WithFusionType(RelativeScore)
			},
			want: `hybrid:{query: "some query", fusionType: relativeScoreFusion}`,
		},
		{
			name: "query and alpha and targetVectors",
			apply: func(h *HybridArgumentBuilder) {
				h.WithQuery("query").
					WithAlpha(0.6).
					WithTargetVectors("t1")
			},
			want: `hybrid:{query: "query", alpha: 0.6, targetVectors: ["t1"]}`,
		},
		{
			name: "query and nearText search",
			apply: func(h *HybridArgumentBuilder) {
				var (
					text     NearTextArgumentBuilder
					searches HybridSearchesArgumentBuilder
				)
				text.WithConcepts([]string{"concept"}).WithCertainty(0.9)
				searches.WithNearText(&text)
				h.WithQuery("I'm a simple string").
					WithSearches(&searches)
			},
			want: `hybrid:{query: "I'm a simple string", searches:{nearText:{concepts: ["concept"] certainty: 0.9}}}`,
		},
		{
			name: "query and nearVector search",
			apply: func(h *HybridArgumentBuilder) {
				var (
					vector   NearVectorArgumentBuilder
					searches HybridSearchesArgumentBuilder
				)
				vector.WithVector([]float32{0.1, 0.2, 0.3}).WithCertainty(0.9)
				searches.WithNearVector(&vector)
				h.WithQuery("I'm a simple string").WithSearches(&searches)
			},
			want: `hybrid:{query: "I'm a simple string", searches:{nearVector:{certainty: 0.9 vector: [0.1,0.2,0.3]}}}`,
		},
		{
			name: "nearVector with maxVectorDistance",
			apply: func(h *HybridArgumentBuilder) {
				var (
					vector   NearVectorArgumentBuilder
					searches HybridSearchesArgumentBuilder
				)
				vector.WithVector([]float32{0.1, 0.2, 0.3})
				searches.WithNearVector(&vector)
				h.WithQuery("I'm a simple string").
					WithSearches(&searches).
					WithMaxVectorDistance(0.8)
			},
			want: `hybrid:{query: "I'm a simple string", maxVectorDistance: 0.8, searches:{nearVector:{vector: [0.1,0.2,0.3]}}}`,
		},
		{
			name: "hybrid with bm25SearchOperator (OR)",
			apply: func(h *HybridArgumentBuilder) {
				var bm25 BM25SearchOperatorBuilder
				bm25.WithOperator(BM25SearchOperatorOr).WithMinimumMatch(4)
				h.WithQuery("hello").WithBM25SearchOperator(bm25)
			},
			want: `hybrid:{query: "hello", bm25SearchOperator:{operator:Or minimumOrTokensMatch:4}}`,
		},
		{
			name: "hybrid with bm25SearchOperator (AND)",
			apply: func(h *HybridArgumentBuilder) {
				var bm25 BM25SearchOperatorBuilder
				bm25.WithOperator(BM25SearchOperatorAnd)
				h.WithQuery("hello").WithBM25SearchOperator(bm25)
			},
			want: `hybrid:{query: "hello", bm25SearchOperator:{operator:And minimumOrTokensMatch:0}}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var hybrid HybridArgumentBuilder
			tt.apply(&hybrid)

			got := hybrid.build()

			require.Equal(t, tt.want, got)
		})
	}
}
