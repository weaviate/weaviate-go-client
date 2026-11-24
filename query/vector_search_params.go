package query

import "github.com/weaviate/weaviate-go-client/v6/internal/api"

// VectorSearchParams tune the vector index for a single request.
//
// Zero fields are not sent, so the collection's vectorIndexConfig applies
// to them. Each field targets one index type, and an index ignores the
// fields it does not know:
//
//   - Ef is the HNSW layer-0 beam width. It also applies to the HNSW stage
//     of a dynamic index.
//   - SearchProbe is the number of centroid postings an HFresh index visits.
//   - FilterStrategy selects the HNSW strategy for the filtered graph walk.
//     It has no effect without a filter.
//
// The server rejects values above its HNSW_MAX_SEARCH_EF and
// HFRESH_MAX_SEARCH_PROBE limits. Servers older than v1.42 ignore these
// parameters. Not applied by Aggregate.
//
//	query.VectorSearchParams{Ef: 256, FilterStrategy: query.FilterStrategyACORN}
type VectorSearchParams api.VectorSearchParams

// FilterStrategy selects the HNSW strategy for a filtered graph walk.
type FilterStrategy = api.FilterStrategy

const (
	// FilterStrategySweeping walks the graph and discards nodes that fail the filter.
	FilterStrategySweeping = api.FilterStrategySweeping
	// FilterStrategyACORN expands the neighbourhood of filtered-out nodes
	// to keep the walk connected on restrictive filters.
	FilterStrategyACORN = api.FilterStrategyACORN
	// FilterStrategyPathseer follows filtered-out nodes along a bounded path.
	FilterStrategyPathseer = api.FilterStrategyPathseer
)

// marshalTargetVectorSearchParams converts per-target params for [api.SearchRequest].
// It returns nil for an empty map so that nothing is sent.
func marshalTargetVectorSearchParams(m map[string]VectorSearchParams) map[string]api.VectorSearchParams {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]api.VectorSearchParams, len(m))
	for target, p := range m {
		out[target] = api.VectorSearchParams(p)
	}
	return out
}
