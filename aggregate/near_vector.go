package aggregate

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type NearVector struct {
	*query.NearVector

	Text    []Text
	Integer []Integer
	Number  []Number
	Boolean []Boolean
	Date    []Date

	TotalCount  bool
	Limit       int32
	ObjectLimit int32
}

type NearVectorFunc func(context.Context, NearVector) (any, error)

func nearVector(ctx context.Context, t internal.Transport, rd api.RequestDefaults, nv NearVector) (any, error) {
	req := &api.AggregateRequest{
		RequestDefaults: rd,
		TotalCount:      nv.TotalCount,
		Limit:           nv.Limit,
		ObjectLimit:     nv.ObjectLimit,

		Text:    make(map[string]*api.AggregateTextRequest, len(nv.Text)),
		Integer: make(map[string]*api.AggregateIntegerRequest, len(nv.Integer)),
		Number:  make(map[string]*api.AggregateNumberRequest, len(nv.Number)),
		Boolean: make(map[string]*api.AggregateBooleanRequest, len(nv.Boolean)),
		Date:    make(map[string]*api.AggregateDateRequest, len(nv.Date)),
	}

	for _, txt := range nv.Text {
		req.Text[txt.Property] = &api.AggregateTextRequest{
			Count:               txt.Count,
			TopOccurrences:      txt.TopOccurrences,
			TopOccurencesCutoff: txt.TopOccurencesCutoff,
		}
	}
	for _, int := range nv.Integer {
		req.Integer[int.Property] = &api.AggregateIntegerRequest{
			Count:   int.Count,
			Minimum: int.Min,
			Maximum: int.Max,
			Mode:    int.Mode,
			Mean:    int.Mean,
			Median:  int.Median,
		}
	}
	for _, num := range nv.Number {
		req.Number[num.Property] = &api.AggregateNumberRequest{
			Count:   num.Count,
			Minimum: num.Min,
			Maximum: num.Max,
			Mode:    num.Mode,
			Mean:    num.Mean,
			Median:  num.Median,
		}
	}
	for _, bool := range nv.Boolean {
		req.Boolean[bool.Property] = &api.AggregateBooleanRequest{
			Count:           bool.Count,
			PercentageTrue:  bool.PercentageTrue,
			PercentageFalse: bool.PercentageFalse,
			TotalTrue:       bool.TotalTrue,
			TotalFalse:      bool.TotalFalse,
		}
	}
	for _, date := range nv.Date {
		req.Date[date.Property] = &api.AggregateDateRequest{
			Count:   date.Count,
			Minimum: date.Min,
			Maximum: date.Max,
			Mode:    date.Mode,
			Median:  date.Median,
		}
	}

	var resp api.AggregateResponse
	if err := t.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("near vector: %w", err)
	}

	result := &Result{
		Text:    make(map[string]TextResult, len(resp.Results.Text)),
		Integer: make(map[string]IntegerResult, len(resp.Results.Integer)),
		Number:  make(map[string]NumberResult, len(resp.Results.Number)),
		Boolean: make(map[string]BooleanResult, len(resp.Results.Boolean)),
		Date:    make(map[string]DateResult, len(resp.Results.Date)),
	}

	for property, txt := range resp.Results.Text {
		top := make([]TopOccurrence, len(txt.TopOccurences))
		for i, item := range txt.TopOccurences {
			top[i] = TopOccurrence{
				Value:       item.Value,
				OccursTimes: item.Occurs,
			}
		}
		result.Text[property] = TextResult{
			Count:          txt.Count,
			TopOccurrences: top,
		}
	}
	for property, int := range resp.Results.Integer {
		result.Integer[property] = IntegerResult{
			Count:  int.Count,
			Min:    int.Minimum,
			Max:    int.Minimum,
			Mode:   int.Mode,
			Mean:   int.Mean,
			Median: int.Median,
		}
	}
	for property, num := range resp.Results.Number {
		result.Number[property] = NumberResult{
			Count:  num.Count,
			Min:    num.Minimum,
			Max:    num.Minimum,
			Mode:   num.Mode,
			Mean:   num.Mean,
			Median: num.Median,
		}
	}
	for property, bool := range resp.Results.Boolean {
		result.Boolean[property] = BooleanResult{
			Count:           bool.Count,
			PercentageTrue:  bool.PercentageTrue,
			PercentageFalse: bool.PercentageFalse,
			TotalTrue:       bool.TotalTrue,
			TotalFalse:      bool.TotalFalse,
		}
	}
	for property, date := range resp.Results.Date {
		result.Date[property] = DateResult{
			Count:  date.Count,
			Min:    date.Minimum,
			Max:    date.Minimum,
			Mode:   date.Mode,
			Median: date.Median,
		}
	}

	return result, nil
}

func nearVectorFunc(t internal.Transport, rd api.RequestDefaults) NearVectorFunc {
	return func(ctx context.Context, nv NearVector) (any, error) {
		return nearVector(ctx, t, rd, nv)
	}
}
