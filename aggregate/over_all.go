package aggregate

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type (
	// Request parameters for an aggregation without a query-filter.
	OverAll struct {
		Text    []Text    // Aggregations for text properties.
		Integer []Integer // Aggregations for integer properties.
		Number  []Number  // Aggregations for number properties.
		Boolean []Boolean // Aggregations for boolean properties.
		Date    []Date    // Aggregations for date properties.

		TotalCount  bool // Return total object count.
		ObjectLimit int32
		groupBy     *GroupBy
	}
	OverAllFunc func(context.Context, OverAll) (*Result, error)
)

// nearVectorFunc encloses transport and request defaults into OverAllFunc scope.
func overAllFunc(t internal.Transport, rd api.RequestDefaults) OverAllFunc {
	return func(ctx context.Context, oa OverAll) (*Result, error) {
		return aggregate(ctx, t, rd, &Request[any]{
			Text:        oa.Text,
			Integer:     oa.Integer,
			Number:      oa.Number,
			Boolean:     oa.Boolean,
			Date:        oa.Date,
			TotalCount:  oa.TotalCount,
			ObjectLimit: oa.ObjectLimit,
			groupBy:     oa.groupBy,
		}, nil, "over all")
	}
}

// GroupBy runs over all aggregation with a GroupBy clause.
func (oaf OverAllFunc) GroupBy(ctx context.Context, oa OverAll, groupBy GroupBy) (*GroupByResult, error) {
	oa.groupBy = &groupBy
	ctx = contextWithGroupByResult(ctx) // safe to reassign since we hold the copy of the original context.
	if _, err := oaf(ctx, oa); err != nil {
		return nil, err
	}
	return getGroupByResult(ctx), nil
}
