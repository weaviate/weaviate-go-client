package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

type OverAll struct {
	Limit                  int              // Limit the number of results returned for the query.
	Offset                 int              // Skip the first N objects in the collection.
	AutoLimit              int              // Return objects in the first N similarity clusters.
	After                  uuid.UUID        // Skip all objects before the one with this ID.
	Filter                 filter.Expr      // Filter results based on their properties.
	Sort                   []SortExpr       // Sort search results by property value.
	ReturnMetadata         ReturnMetadata   // Select query and object metadata to return for each object.
	ReturnVectors          []string         // List vectors to return for each object in the result set.
	ReturnReferences       []Reference      // Select reference properties to return.
	ReturnNestedProperties []NestedProperty // Return object properties and a subset of their nested properties.

	// Select a subset of properties to return. By default, all properties are returned.
	// To not return any properties, initialize this value to an empty slice explicitly.
	ReturnProperties []string

	// groupBy can only be set by [NearVectorFunc.GroupBy], as it changes the shape of the response.
	groupBy *GroupBy
}

// OverAllFunc runs plain OverAll search.
type OverAllFunc func(context.Context, OverAll) (*Result, error)

// overAllFunc makes internal.Transport available to [query] via a closure.
func overAllFunc(t internal.Transport, rd api.RequestDefaults) OverAllFunc {
	return func(ctx context.Context, oa OverAll) (*Result, error) {
		return query(ctx, t, oa.QueryRequest(rd, oa.groupBy))
	}
}

func (oa *OverAll) QueryRequest(rd api.RequestDefaults, gb *GroupBy) Request {
	return Request{
		RequestDefaults:        rd,
		Limit:                  int32(oa.Limit),
		AutoLimit:              int32(oa.AutoLimit),
		Offset:                 int32(oa.Offset),
		After:                  oa.After,
		Filter:                 oa.Filter,
		Sort:                   oa.Sort,
		ReturnVectors:          oa.ReturnVectors,
		ReturnMetadata:         oa.ReturnMetadata,
		ReturnProperties:       oa.ReturnProperties,
		ReturnNestedProperties: oa.ReturnNestedProperties,
		ReturnReferences:       oa.ReturnReferences,
		GroupBy:                gb,
		Label:                  "over all",
	}
}

// GroupBy runs OverAll search with a GroupBy clause.
func (oaf OverAllFunc) GroupBy(ctx context.Context, oa OverAll, groupBy GroupBy) (*GroupByResult, error) {
	oa.groupBy = &groupBy
	return queryGroupBy(ctx, oaf, oa)
}
