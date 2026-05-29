package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

type Hybrid struct {
	Limit                  int32            // Limit the number of results returned for the query.
	Offset                 int32            // Skip the first N objects in the collection.
	AutoLimit              int32            // Return objects in the first N similarity clusters.
	After                  uuid.UUID        // Skip all objects before the one with this ID.
	Filter                 filter.Expr      // Filter results based on their properties.
	ReturnMetadata         ReturnMetadata   // Select query and object metadata to return for each object.
	ReturnVectors          []string         // List vectors to return for each object in the result set.
	ReturnReferences       []Reference      // Select reference properties to return.
	ReturnNestedProperties []NestedProperty // Return object properties and a subset of their nested properties.

	// Select a subset of properties to return. By default, all properties are returned.
	// To not return any properties, initialize this value to an empty slice explicitly.
	ReturnProperties []string

	// Keyword-based search query. Required parameter.
	Query string

	// QueryProperties limits BM25 search to the selected properties.
	QueryProperties []string

	// NearVector parameters for vector search.
	// Only set ONE of [Hybrid.NearText] / [Hybrid.NearVector].
	NearVector *NearVector

	// NearText parameters for vector search.
	// Only set ONE of [Hybrid.NearText] / [Hybrid.NearVector].
	NearText *NearText

	// Alpha controls bias between BM25 and vector search result sets.
	// An alpha of 0 is pure vector search; an Alpha of 1 is pure BM25 search.
	Alpha *float32

	// Fusion algorithm for combining BM25 and vector search result sets.
	Fusion api.HybridFusion

	// Similarity threshold for BM25 search component.
	KeywordSimilarity KeywordSimilarity

	// groupBy can only be set by [NearTextFunc.GroupBy], as it changes the shape of the response.
	groupBy *GroupBy
}

// Alpha is a helper for setting optional [Hybrid.Alpha].
func Alpha(a float32) *float32 { return &a }

const (
	// RANKED fusion algorithm.
	HybridFusionRanked = api.HybridFusionRanked
	// RELATIVE_SCORE fusion algorithm.
	HybridFusionRelativeScore = api.HybridFusionRelativeScore
)

// HybridFunc runs plain near text search.
type HybridFunc func(context.Context, Hybrid) (*Result, error)

// hybridFunc makes internal.Transport available to nearText via a closure.
func hybridFunc(t internal.Transport, rd api.RequestDefaults) HybridFunc {
	return func(ctx context.Context, h Hybrid) (*Result, error) {
		return query(ctx, t, request{
			RequestDefaults:        rd,
			Limit:                  h.Limit,
			AutoLimit:              h.AutoLimit,
			Offset:                 h.Offset,
			After:                  h.After,
			Filter:                 h.Filter,
			ReturnVectors:          h.ReturnVectors,
			ReturnMetadata:         h.ReturnMetadata,
			ReturnProperties:       h.ReturnProperties,
			ReturnNestedProperties: h.ReturnNestedProperties,
			ReturnReferences:       h.ReturnReferences,
			GroupBy:                h.groupBy,
		}, func(req *api.SearchRequest) {
			req.Hybrid = &api.Hybrid{
				Query:           h.Query,
				QueryProperties: h.QueryProperties,
				Alpha:           h.Alpha,
				Fusion:          h.Fusion,
				KeywordSimilarity: api.KeywordSimilarity{
					AllTokensMatch:     h.KeywordSimilarity.AllTokensMatch(),
					MinimumTokensMatch: h.KeywordSimilarity.MinimumTokensMatch(),
				},
				NearVector: nearVector(h.NearVector),
				NearText:   nearText(h.NearText),
			}
		})
	}
}

// GroupBy runs near text search with a GroupBy clause.
func (hf HybridFunc) GroupBy(ctx context.Context, h Hybrid, groupBy GroupBy) (*GroupByResult, error) {
	h.groupBy = &groupBy
	return queryGroupBy(ctx, hf, h)
}
