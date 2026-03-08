package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

// TestVectorTarget tests existing implementations of [query.VectorTarget] and [query.VectorKind]
// in combination with [query.WeightedVector] and [query.MultiVectorTarget] to ensure any valid
// combination of target vectors can be expressed and produces valid api.SearchTarget.
func TestVectorTarget(t *testing.T) {
	type combinationMethoder interface {
		CombinationMethod() api.CombinationMethod
	}

	var (
		singleVector = []float32{1, 2, 3}
		multiVector  = [][]float32{{1, 2, 3}, {1, 2, 3}}
	)

	for _, tt := range []struct {
		name   string
		target query.VectorTarget
		want   api.SearchTarget
	}{
		{
			name:   "unnamed single vector",
			target: types.Vector{Single: singleVector},
			want: api.SearchTarget{Vectors: []api.TargetVector{
				{Vector: api.Vector{Single: singleVector}},
			}},
		},
		{
			name:   "unnamed multi-vector",
			target: types.Vector{Multi: multiVector},
			want: api.SearchTarget{Vectors: []api.TargetVector{
				{Vector: api.Vector{Multi: multiVector}},
			}},
		},
		{
			name:   "named single vector",
			target: types.Vector{Name: "v1", Single: singleVector},
			want: api.SearchTarget{Vectors: []api.TargetVector{
				{Vector: api.Vector{Name: "v1", Single: singleVector}},
			}},
		},
		{
			name:   "named multi-vector",
			target: types.Vector{Name: "v1", Multi: multiVector},
			want: api.SearchTarget{Vectors: []api.TargetVector{
				{Vector: api.Vector{Name: "v1", Multi: multiVector}},
			}},
		},
		{
			name:   "vector index name",
			target: query.VectorName("v1"),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1"}},
				},
			},
		},
		{
			name: "sum with vector embeddings",
			target: query.Sum([]types.Vector{
				{Name: "v1", Single: singleVector},
				{Name: "v2", Multi: multiVector},
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1", Single: singleVector}},
					{Vector: api.Vector{Name: "v2", Multi: multiVector}},
				}, CombinationMethod: api.CombinationMethodSum,
			},
		},
		{
			name: "sum with vector names",
			target: query.Sum([]query.VectorName{
				"v1", "v2",
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1"}},
					{Vector: api.Vector{Name: "v2"}},
				}, CombinationMethod: api.CombinationMethodSum,
			},
		},
		{
			name: "average with vector embeddings",
			target: query.Average([]types.Vector{
				{Name: "v1", Single: singleVector},
				{Name: "v2", Multi: multiVector},
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1", Single: singleVector}},
					{Vector: api.Vector{Name: "v2", Multi: multiVector}},
				}, CombinationMethod: api.CombinationMethodAverage,
			},
		},
		{
			name: "average with vector names",
			target: query.Average([]query.VectorName{
				"v1", "v2",
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1"}},
					{Vector: api.Vector{Name: "v2"}},
				}, CombinationMethod: api.CombinationMethodAverage,
			},
		},
		{
			name: "min with vector embeddings",
			target: query.Min([]types.Vector{
				{Name: "v1", Single: singleVector},
				{Name: "v2", Multi: multiVector},
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1", Single: singleVector}},
					{Vector: api.Vector{Name: "v2", Multi: multiVector}},
				}, CombinationMethod: api.CombinationMethodMin,
			},
		},
		{
			name: "min with vector names",
			target: query.Min([]query.VectorName{
				"v1", "v2",
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1"}},
					{Vector: api.Vector{Name: "v2"}},
				}, CombinationMethod: api.CombinationMethodMin,
			},
		},
		{
			name: "manual weights with vector embeddings",
			target: query.ManualWeights([]query.WeightedVector[types.Vector]{
				query.Weighted(types.Vector{Name: "v1", Single: singleVector}, 0.22),
				query.Weighted(types.Vector{Name: "v2", Multi: multiVector}, 0.88),
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1", Single: singleVector}, Weight: testkit.Ptr[float32](0.22)},
					{Vector: api.Vector{Name: "v2", Multi: multiVector}, Weight: testkit.Ptr[float32](0.88)},
				}, CombinationMethod: api.CombinationMethodManualWeights,
			},
		},
		{
			name: "manual weights with vector names",
			target: query.ManualWeights([]query.WeightedVector[query.VectorName]{
				query.Weighted(query.VectorName("v1"), 0.22),
				query.Weighted(query.VectorName("v2"), 0.88),
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1"}, Weight: testkit.Ptr[float32](0.22)},
					{Vector: api.Vector{Name: "v2"}, Weight: testkit.Ptr[float32](0.88)},
				}, CombinationMethod: api.CombinationMethodManualWeights,
			},
		},
		{
			name: "relative score with vector embeddings",
			target: query.RelativeScore([]query.WeightedVector[types.Vector]{
				query.Weighted(types.Vector{Name: "v1", Single: singleVector}, 0.22),
				query.Weighted(types.Vector{Name: "v2", Multi: multiVector}, 0.88),
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1", Single: singleVector}, Weight: testkit.Ptr[float32](0.22)},
					{Vector: api.Vector{Name: "v2", Multi: multiVector}, Weight: testkit.Ptr[float32](0.88)},
				}, CombinationMethod: api.CombinationMethodRelativeScore,
			},
		},
		{
			name: "relative score with vector names",
			target: query.RelativeScore([]query.WeightedVector[query.VectorName]{
				query.Weighted(query.VectorName("v1"), 0.22),
				query.Weighted(query.VectorName("v2"), 0.88),
			}),
			want: api.SearchTarget{
				Vectors: []api.TargetVector{
					{Vector: api.Vector{Name: "v1"}, Weight: testkit.Ptr[float32](0.22)},
					{Vector: api.Vector{Name: "v2"}, Weight: testkit.Ptr[float32](0.88)},
				},
				CombinationMethod: api.CombinationMethodRelativeScore,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want.CombinationMethod != "" && assert.Implements(t, (*combinationMethoder)(nil), tt.target) {
				cm := tt.target.(combinationMethoder)
				assert.Equal(t, tt.want.CombinationMethod, cm.CombinationMethod())
			}
			require.Equal(t, tt.want.Vectors, tt.target.Vectors(), "target vectors")
		})
	}
}
