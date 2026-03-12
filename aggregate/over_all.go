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
		Limit       int32
		ObjectLimit int32
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
			Limit:       oa.Limit,
			ObjectLimit: oa.ObjectLimit,
		}, nil, "over all")
	}
}
