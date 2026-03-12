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
	Limit       int32
	ObjectLimit int32
}

type (
	Text    api.AggregateTextRequest
	Integer api.AggregateIntegerRequest
	Number  api.AggregateNumberRequest
	Boolean api.AggregateBooleanRequest
	Date    api.AggregateDateRequest
)

type GroupBy struct {
	Collection string
	Property   string
}

type Result struct {
	Text    map[string]TextResult
	Integer map[string]IntegerResult
	Number  map[string]NumberResult
	Boolean map[string]BooleanResult
	Date    map[string]DateResult

	TotalCount  *int64
	TookSeconds float32
}

type (
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
		Limit:           r.Limit,
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
		// Conversion to any, while a bit awkward, enables us to do
		// value-type dispatch here; it smartly bridges the gap between
		// two disparate type sets: what `req` may be (query, aggregate, generate),
		// and what Query.Request() may return (near vector, hybrid, bm25).
		switch q := search.(type) {
		case *api.NearVector:
			req.NearVector = q
		}
	}

	var resp api.AggregateResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("%s: %w", label, err)
	}

	result := &Result{
		TotalCount:  resp.TotalCount,
		TookSeconds: resp.TookSeconds,

		Text:    make(map[string]TextResult, len(resp.Text)),
		Integer: make(map[string]IntegerResult, len(resp.Integer)),
		Number:  make(map[string]NumberResult, len(resp.Number)),
		Boolean: make(map[string]BooleanResult, len(resp.Boolean)),
		Date:    make(map[string]DateResult, len(resp.Date)),
	}
	for property, txt := range resp.Text {
		top := make([]TopOccurrence, len(txt.TopOccurrences))
		for i, item := range txt.TopOccurrences {
			top[i] = TopOccurrence(item)
		}
		result.Text[property] = TextResult{
			Count:          txt.Count,
			TopOccurrences: top,
		}
	}
	for property, int := range resp.Integer {
		result.Integer[property] = IntegerResult(int)
	}
	for property, num := range resp.Number {
		result.Number[property] = NumberResult(num)
	}
	for property, bool := range resp.Boolean {
		result.Boolean[property] = BooleanResult(bool)
	}
	for property, date := range resp.Date {
		result.Date[property] = DateResult(date)
	}
	return result, nil
}
