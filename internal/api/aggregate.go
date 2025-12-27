package api

import (
	"cmp"
	"fmt"
	"iter"
	"log"
	"maps"
	"slices"
	"time"

	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

type AggregateRequest struct {
	RequestDefaults

	Text    map[string]*AggregateTextRequest
	Integer map[string]*AggregateIntegerRequest
	Number  map[string]*AggregateNumberRequest
	Boolean map[string]*AggregateBooleanRequest
	Date    map[string]*AggregateDateRequest

	TotalCount  bool
	Limit       int32
	ObjectLimit int32
}

var _ transport.MessageMarshaler[proto.AggregateRequest] = (*AggregateRequest)(nil)

func (r *AggregateRequest) MarshalMessage() *proto.AggregateRequest {
	var aggs []*proto.AggregateRequest_Aggregation
	for property, txt := range sortedMap(r.Text) {
		aggs = append(aggs, &proto.AggregateRequest_Aggregation{
			Property: property,
			Aggregation: &proto.AggregateRequest_Aggregation_Text_{
				Text: &proto.AggregateRequest_Aggregation_Text{
					Count:              txt.Count,
					TopOccurences:      txt.TopOccurrences,
					TopOccurencesLimit: nilZero(uint32(txt.TopOccurencesCutoff)),
				},
			},
		})
	}
	for property, int := range sortedMap(r.Integer) {
		aggs = append(aggs, &proto.AggregateRequest_Aggregation{
			Property: property,
			Aggregation: &proto.AggregateRequest_Aggregation_Int{
				Int: (*proto.AggregateRequest_Aggregation_Integer)(int),
			},
		})
	}
	for property, num := range sortedMap(r.Number) {
		aggs = append(aggs, &proto.AggregateRequest_Aggregation{
			Property: property,
			Aggregation: &proto.AggregateRequest_Aggregation_Number_{
				Number: (*proto.AggregateRequest_Aggregation_Number)(num),
			},
		})
	}
	for property, bool := range sortedMap(r.Boolean) {
		aggs = append(aggs, &proto.AggregateRequest_Aggregation{
			Property: property,
			Aggregation: &proto.AggregateRequest_Aggregation_Boolean_{
				Boolean: (*proto.AggregateRequest_Aggregation_Boolean)(bool),
			},
		})
	}
	for property, date := range sortedMap(r.Date) {
		aggs = append(aggs, &proto.AggregateRequest_Aggregation{
			Property: property,
			Aggregation: &proto.AggregateRequest_Aggregation_Date_{
				Date: (*proto.AggregateRequest_Aggregation_Date)(date),
			},
		})
	}
	return &proto.AggregateRequest{
		Collection: r.CollectionName,
		Tenant:     r.Tenant,

		ObjectsCount: r.TotalCount,
		Limit:        nilZero(uint32(r.Limit)),
		ObjectLimit:  nilZero(uint32(r.ObjectLimit)),
		Aggregations: aggs,
	}
}

type (
	AggregateTextRequest struct {
		Count               bool
		TopOccurrences      bool
		TopOccurencesCutoff int32
	}
	AggregateIntegerRequest proto.AggregateRequest_Aggregation_Integer
	AggregateNumberRequest  proto.AggregateRequest_Aggregation_Number
	AggregateBooleanRequest proto.AggregateRequest_Aggregation_Boolean
	AggregateDateRequest    proto.AggregateRequest_Aggregation_Date
)

type Aggregations struct {
	Text    map[string]*AggregateTextResult
	Integer map[string]*AggregateIntegerResult
	Number  map[string]*AggregateNumberResult
	Boolean map[string]*AggregateBooleanResult
	Date    map[string]*AggregateDateResult

	TotalCount *int64
}

type AggregateResponse struct {
	Results Aggregations

	TookSeconds float32
}

var _ transport.MessageUnmarshaler[proto.AggregateReply] = (*AggregateResponse)(nil)

func (r *AggregateResponse) UnmarshalMessage(reply *proto.AggregateReply) error {
	result := Aggregations{
		Text:    make(map[string]*AggregateTextResult),
		Integer: make(map[string]*AggregateIntegerResult),
		Number:  make(map[string]*AggregateNumberResult),
		Boolean: make(map[string]*AggregateBooleanResult),
		Date:    make(map[string]*AggregateDateResult),
	}
	single := reply.GetSingleResult()
	if single != nil {
		result.TotalCount = single.ObjectsCount
		for _, agg := range single.GetAggregations().GetAggregations() {
			property := agg.GetProperty()
			switch {
			case agg.GetText() != nil:
				txt := agg.GetText()
				top := make([]*TopOccurence, len(txt.GetTopOccurences().GetItems()))
				for i, item := range txt.GetTopOccurences().GetItems() {
					top[i] = (*TopOccurence)(item)
				}
				result.Text[property] = &AggregateTextResult{
					Count:         txt.Count,
					TopOccurences: top,
				}
			case agg.GetDate() != nil:
				date := agg.GetDate()
				minimum, err := parseDate(date.GetMinimum())
				if err != nil {
					return fmt.Errorf("%q minimum: %w", property, err)
				}
				log.Print(testkit.Now.Format(time.RFC3339))
				log.Print(date.GetMinimum())
				log.Print(minimum)
				maximum, err := parseDate(date.GetMaximum())
				if err != nil {
					return fmt.Errorf("%q maximum: %w", property, err)
				}
				mode, err := parseDate(date.GetMode())
				if err != nil {
					return fmt.Errorf("%q mode: %w", property, err)
				}
				median, err := parseDate(date.GetMedian())
				if err != nil {
					return fmt.Errorf("%q median: %w", property, err)
				}
				result.Date[property] = &AggregateDateResult{
					Count:   date.Count,
					Minimum: minimum,
					Maximum: maximum,
					Mode:    mode,
					Median:  median,
				}
				log.Print(testkit.Now)
				log.Print(result.Date[property].Minimum)
			case agg.GetInt() != nil:
				result.Integer[property] = (*AggregateIntegerResult)(agg.GetInt())
			case agg.GetNumber() != nil:
				result.Number[property] = (*AggregateNumberResult)(agg.GetNumber())
			case agg.GetBoolean() != nil:
				result.Boolean[property] = (*AggregateBooleanResult)(agg.GetBoolean())
			}
		}
	}

	*r = AggregateResponse{
		TookSeconds: reply.GetTook(),
		Results:     result,
	}
	return nil
}

type (
	TopOccurence        proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences_TopOccurrence
	AggregateTextResult struct {
		Count         *int64
		TopOccurences []*TopOccurence
	}
	AggregateIntegerResult proto.AggregateReply_Aggregations_Aggregation_Integer
	AggregateNumberResult  proto.AggregateReply_Aggregations_Aggregation_Number
	AggregateBooleanResult proto.AggregateReply_Aggregations_Aggregation_Boolean
	AggregateDateResult    struct {
		Count   *int64
		Minimum *time.Time
		Maximum *time.Time
		Mode    *time.Time
		Median  *time.Time
	}
)

// parseDate parse date string assuming RFC3339 format.
// It returns nil if the input string is empty.
func parseDate(date string) (*time.Time, error) {
	if date == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// sortedMap returns an iterator over key-value pairs from m;
// similar to [maps.All], but with pairs sorted by key.
func sortedMap[Map ~map[K]V, K cmp.Ordered, V any](m Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range slices.Sorted(maps.Keys(m)) {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
