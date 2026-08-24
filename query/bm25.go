package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

type BM25 struct {
	Limit                  int              // Limit the number of results returned for the query.
	Offset                 int              // Skip the first N objects in the collection.
	AutoLimit              int              // Return objects in the first N similarity clusters.
	After                  uuid.UUID        // Skip all objects before the one with this ID.
	Filter                 filter.Expr      // Filter results based on their properties.
	ReturnMetadata         ReturnMetadata   // Select query and object metadata to return for each object.
	ReturnVectors          []string         // List vectors to return for each object in the result set.
	ReturnReferences       []Reference      // Select reference properties to return.
	ReturnNestedProperties []NestedProperty // Return object properties and a subset of their nested properties.

	// Select a subset of properties to return. By default, all properties are returned.
	// To not return any properties, initialize this value to an empty slice explicitly.
	ReturnProperties []string

	// Query is tokenized and the resulting keywords are searched across [QueryProeprties].
	// Required parameter.
	Query string

	// Select a subset of properties to search for keywords. By default, all searchable
	// text/text[] properties are considered.
	QueryProperties []string

	// Similarity threshold.
	KeywordSimilarity KeywordSimilarity

	// groupBy can only be set by [BM25Func.GroupBy], as it changes the shape of the response.
	groupBy *GroupBy
}

type (
	// KeywordSimilarity conrols the similarity threshold for BM25 (keyword) search.
	KeywordSimilarity struct {
		// allTokensMatch requires each token in the search query to be present
		// in a candidate object for it to be considered a match.
		allTokensMatch bool

		// crossProperty counts token matches from all query properties towards
		// the similarity target.
		crossProperty bool

		// minimumTokensMatch is the lower threshold for the number of times
		// _each_ token needs be present in a candidate object for it to be considered a match.
		mininumTokensMatch *int32
	}
)

func (kws *KeywordSimilarity) AllTokensMatch() bool       { return kws.allTokensMatch }
func (kws *KeywordSimilarity) CrossProperty() bool        { return kws.crossProperty }
func (kws *KeywordSimilarity) MinimumTokensMatch() *int32 { return kws.mininumTokensMatch }

// AllTokensMatch is a [KeywordSimilarity] parameter with allTokensMatch=true.
var AllTokensMatch = KeywordSimilarity{allTokensMatch: true}

// AllTokensMatchCross is a [KeywordSimilarity] parameter like [AllTokensMatch],
// but counting matches from all query properties towards the similarity target.
var AllTokensMatchCross = KeywordSimilarity{allTokensMatch: true, crossProperty: true}

// MinimumTokensMatch returns [KeywordSimilarity] with MinimumTokensMatch=n.
func MinimumTokensMatch(n int32) KeywordSimilarity {
	return KeywordSimilarity{mininumTokensMatch: &n}
}

// BM25Func runs plain near text search.
type BM25Func func(context.Context, BM25) (*Result, error)

// bm25Func makes internal.Transport available to [query] via a closure.
func bm25Func(t internal.Transport, rd api.RequestDefaults) BM25Func {
	return func(ctx context.Context, bm25 BM25) (*Result, error) {
		return query(ctx, t, request{
			RequestDefaults:        rd,
			Limit:                  int32(bm25.Limit),
			AutoLimit:              int32(bm25.AutoLimit),
			Offset:                 int32(bm25.Offset),
			After:                  bm25.After,
			Filter:                 bm25.Filter,
			ReturnVectors:          bm25.ReturnVectors,
			ReturnMetadata:         bm25.ReturnMetadata,
			ReturnProperties:       bm25.ReturnProperties,
			ReturnNestedProperties: bm25.ReturnNestedProperties,
			ReturnReferences:       bm25.ReturnReferences,
			GroupBy:                bm25.groupBy,
		}, func(req *api.SearchRequest) {
			req.BM25 = &api.BM25{
				Query:           bm25.Query,
				QueryProperties: bm25.QueryProperties,
				KeywordSimilarity: api.KeywordSimilarity{
					AllTokensMatch:     bm25.KeywordSimilarity.AllTokensMatch(),
					CrossProperty:      bm25.KeywordSimilarity.CrossProperty(),
					MinimumTokensMatch: bm25.KeywordSimilarity.MinimumTokensMatch(),
				},
			}
		})
	}
}

// GroupBy runs near text search with a GroupBy clause.
func (bf BM25Func) GroupBy(ctx context.Context, bm25 BM25, groupBy GroupBy) (*GroupByResult, error) {
	bm25.groupBy = &groupBy
	return queryGroupBy(ctx, bf, bm25)
}
