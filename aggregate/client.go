package aggregate

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewClient(t internal.Transport, rd api.RequestDefaults) *Client {
	dev.AssertNotNil(t, "transport")
	return &Client{
		transport:  t,
		defaults:   rd,
		OverAll:    overAllFunc(t, rd),
		NearVector: nearVectorFunc(t, rd),
	}
}

type Client struct {
	transport internal.Transport
	defaults  api.RequestDefaults

	OverAll    OverAllFunc
	NearVector NearVectorFunc
}

// Request contains common aggregation parameters.
type Request[Query any] struct {
	Query Query // Query-filter portion of the request.

	Text    []Text    // Aggregations for text properties.
	Integer []Integer // Aggregations for integer properties.
	Number  []Number  // Aggregations for number properties.
	Boolean []Boolean // Aggregations for boolean properties.
	Date    []Date    // Aggregations for date properties.

	TotalCount  bool // Return total object count.
	ObjectLimit int32

	// groupBy can only be set by a GroupBy method, as it changes the shape of the response.
	groupBy *GroupBy
}

type GroupBy struct {
	Property string
	Limit    int
}

type (
	Text    api.AggregateTextRequest
	Integer api.AggregateIntegerRequest
	Number  api.AggregateNumberRequest
	Boolean api.AggregateBooleanRequest
	Date    api.AggregateDateRequest
)

type Result struct {
	TookSeconds float32
	Aggregations
}

type GroupByResult struct {
	Groups []Group
}

type Group struct {
	Property string
	Value    any
	Aggregations
}

type (
	Aggregations struct {
		TotalCount *int64
		Text       map[string]TextResult
		Integer    map[string]IntegerResult
		Number     map[string]NumberResult
		Boolean    map[string]BooleanResult
		Date       map[string]DateResult
	}
	TextResult struct {
		Count          *int64
		TopOccurrences []TopOccurrence
	}
	TopOccurrence api.TopOccurrence
	IntegerResult api.AggregateIntegerResult
	NumberResult  api.AggregateNumberResult
	BooleanResult api.AggregateBooleanResult
	DateResult    api.AggregateDateResult
)

func aggregate[Query any](ctx context.Context, t internal.Transport, rd api.RequestDefaults, r *Request[Query], search any, label string) (*Result, error) {
	req := &api.AggregateRequest{
		RequestDefaults: rd,
		TotalCount:      r.TotalCount,
		ObjectLimit:     r.ObjectLimit,
	}
	for _, txt := range r.Text {
		req.Text = append(req.Text, api.AggregateTextRequest(txt))
	}
	for _, int := range r.Integer {
		req.Integer = append(req.Integer, api.AggregateIntegerRequest(int))
	}
	for _, num := range r.Number {
		req.Number = append(req.Number, api.AggregateNumberRequest(num))
	}
	for _, bool := range r.Boolean {
		req.Boolean = append(req.Boolean, api.AggregateBooleanRequest(bool))
	}
	for _, date := range r.Date {
		req.Date = append(req.Date, api.AggregateDateRequest(date))
	}

	if search != nil {
		switch q := search.(type) {
		case *api.NearVector:
			req.NearVector = q
		}
	}

	if r.groupBy != nil {
		req.GroupBy = &api.GroupBy{
			Property: r.groupBy.Property,
			Limit:    int32(r.groupBy.Limit),
		}
	}

	var resp api.AggregateResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}

	// aggregate was called from a GroupBy() method.
	// This means we should put GroupByResult in the context,
	// as the first return value will be discarded.
	if r.groupBy != nil {
		groups := make([]Group, len(resp.GroupByResults))
		for gi, group := range resp.GroupByResults {
			groups[gi] = Group{
				Property:     group.Property,
				Value:        group.Value,
				Aggregations: aggregationsFromAPI(group.Results),
			}
		}
		setGroupByResult(ctx, &GroupByResult{Groups: groups})
		return nil, nil
	}

	result := &Result{
		TookSeconds:  resp.TookSeconds,
		Aggregations: aggregationsFromAPI(resp.Results),
	}
	return result, nil
}

func aggregationsFromAPI(aggregations api.Aggregations) Aggregations {
	out := Aggregations{
		TotalCount: aggregations.TotalCount,
		Text:       make(map[string]TextResult, len(aggregations.Text)),
		Integer:    make(map[string]IntegerResult, len(aggregations.Integer)),
		Number:     make(map[string]NumberResult, len(aggregations.Number)),
		Boolean:    make(map[string]BooleanResult, len(aggregations.Boolean)),
		Date:       make(map[string]DateResult, len(aggregations.Date)),
	}
	for property, txt := range aggregations.Text {
		top := make([]TopOccurrence, len(txt.TopOccurrences))
		for i, item := range txt.TopOccurrences {
			top[i] = TopOccurrence(item)
		}
		out.Text[property] = TextResult{
			Count:          txt.Count,
			TopOccurrences: top,
		}
	}
	for property, int := range aggregations.Integer {
		out.Integer[property] = IntegerResult(int)
	}
	for property, num := range aggregations.Number {
		out.Number[property] = NumberResult(num)
	}
	for property, bool := range aggregations.Boolean {
		out.Boolean[property] = BooleanResult(bool)
	}
	for property, date := range aggregations.Date {
		out.Date[property] = DateResult(date)
	}
	return out
}

// groupByResultKey is used to pass grouped query results to the GroupBy caller.
var groupByResultKey = internal.ContextKey{}

// contextWithGorupByResult creates a placeholder for *GroupByResult in the ctx.Values store.
func contextWithGroupByResult(ctx context.Context) context.Context {
	return internal.ContextWithPlaceholder[GroupByResult](ctx, groupByResultKey)
}

// getGroupByResult extracts *GroupByResult from the context.
func getGroupByResult(ctx context.Context) *GroupByResult {
	return internal.ValueFromContext[GroupByResult](ctx, groupByResultKey)
}

// setGroupByResult replaces *GroupByResult placeholder
// in the context with the value at r.
func setGroupByResult(ctx context.Context, r *GroupByResult) {
	internal.SetContextValue(ctx, groupByResultKey, r)
}
