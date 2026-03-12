package api

import (
	"fmt"
	"time"

	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type AggregateRequest struct {
	RequestDefaults

	Text    []AggregateTextRequest
	Integer []AggregateIntegerRequest
	Number  []AggregateNumberRequest
	Boolean []AggregateBooleanRequest
	Date    []AggregateDateRequest

	TotalCount  bool
	ObjectLimit int32
	GroupBy     *GroupBy

	NearVector *NearVector
}
type (
	AggregateTextRequest struct {
		Property string

		Count               bool
		TopOccurrences      bool
		TopOccurencesCutoff int32
	}
	AggregateIntegerRequest struct {
		Property string

		Count  bool
		Sum    bool
		Min    bool
		Max    bool
		Mode   bool
		Mean   bool
		Median bool
	}
	AggregateNumberRequest struct {
		Property string

		Count  bool
		Sum    bool
		Min    bool
		Max    bool
		Mode   bool
		Mean   bool
		Median bool
	}
	AggregateBooleanRequest struct {
		Property string

		Count           bool
		Type            bool
		PercentageTrue  bool
		PercentageFalse bool
		TotalTrue       bool
		TotalFalse      bool
	}
	AggregateDateRequest struct {
		Property string

		Count  bool
		Min    bool
		Max    bool
		Mode   bool
		Median bool
	}
)

func (r *AggregateRequest) Method() transport.MethodFunc[proto.AggregateRequest, proto.AggregateReply] {
	return proto.WeaviateClient.Aggregate
}
func (r *AggregateRequest) Body() transport.MessageMarshaler[proto.AggregateRequest] { return r }

// MarshalMessage implements [Message].
func (r *AggregateRequest) MarshalMessage() (*proto.AggregateRequest, error) {
	var aggregations []*proto.AggregateRequest_Aggregation
	for _, txt := range r.Text {
		aggregations = append(aggregations, &proto.AggregateRequest_Aggregation{
			Property: txt.Property,
			Aggregation: &proto.AggregateRequest_Aggregation_Text_{
				Text: &proto.AggregateRequest_Aggregation_Text{
					Count:              txt.Count,
					TopOccurences:      txt.TopOccurrences,
					TopOccurencesLimit: nilZero(uint32(txt.TopOccurencesCutoff)),
				},
			},
		})
	}
	for _, int := range r.Integer {
		aggregations = append(aggregations, &proto.AggregateRequest_Aggregation{
			Property: int.Property,
			Aggregation: &proto.AggregateRequest_Aggregation_Int{
				Int: &proto.AggregateRequest_Aggregation_Integer{
					Count:   int.Count,
					Sum:     int.Sum,
					Minimum: int.Min,
					Maximum: int.Max,
					Mode:    int.Mode,
					Mean:    int.Mean,
					Median:  int.Median,
				},
			},
		})
	}
	for _, num := range r.Number {
		aggregations = append(aggregations, &proto.AggregateRequest_Aggregation{
			Property: num.Property,
			Aggregation: &proto.AggregateRequest_Aggregation_Number_{
				Number: &proto.AggregateRequest_Aggregation_Number{
					Count:   num.Count,
					Sum:     num.Sum,
					Minimum: num.Min,
					Maximum: num.Max,
					Mode:    num.Mode,
					Mean:    num.Mean,
					Median:  num.Median,
				},
			},
		})
	}
	for _, bool := range r.Boolean {
		aggregations = append(aggregations, &proto.AggregateRequest_Aggregation{
			Property: bool.Property,
			Aggregation: &proto.AggregateRequest_Aggregation_Boolean_{
				Boolean: &proto.AggregateRequest_Aggregation_Boolean{
					Count:           bool.Count,
					Type:            bool.Type,
					PercentageTrue:  bool.PercentageTrue,
					PercentageFalse: bool.PercentageFalse,
					TotalTrue:       bool.TotalTrue,
					TotalFalse:      bool.TotalFalse,
				},
			},
		})
	}
	for _, date := range r.Date {
		aggregations = append(aggregations, &proto.AggregateRequest_Aggregation{
			Property: date.Property,
			Aggregation: &proto.AggregateRequest_Aggregation_Date_{
				Date: &proto.AggregateRequest_Aggregation_Date{
					Count:   date.Count,
					Minimum: date.Min,
					Maximum: date.Max,
					Mode:    date.Mode,
					Median:  date.Median,
				},
			},
		})
	}

	req := &proto.AggregateRequest{
		Collection: r.CollectionName,
		Tenant:     r.Tenant,

		ObjectsCount: r.TotalCount,
		ObjectLimit:  nilZero(uint32(r.ObjectLimit)),
		Aggregations: aggregations,
	}

	if r.GroupBy != nil {
		req.Limit = ptr(uint32(r.GroupBy.Limit))
		req.GroupBy = &proto.AggregateRequest_GroupBy{
			Collection: r.CollectionName,
			Property:   r.GroupBy.Property,
		}
	}

	switch {
	case r.NearVector != nil:
		nv, err := marshalNearVector(r.NearVector)
		if err != nil {
			return nil, err
		}
		req.Search = &proto.AggregateRequest_NearVector{NearVector: nv}
	default:
		// It is not a mistake to leave search method unset.
		// This would be the case when fetch objects with a conventional filter.
	}

	return req, nil
}

type AggregateResponse struct {
	TookSeconds    float32
	Results        Aggregations
	GroupByResults []AggregateGroup
}

type (
	Aggregations struct {
		TotalCount *int64
		Text       map[string]AggregateTextResult
		Integer    map[string]AggregateIntegerResult
		Number     map[string]AggregateNumberResult
		Boolean    map[string]AggregateBooleanResult
		Date       map[string]AggregateDateResult
	}
	AggregateTextResult struct {
		Count          *int64
		TopOccurrences []TopOccurrence
	}
	AggregateTopOccurrence struct {
		Value       string
		OccursTimes int64
	}
	AggregateIntegerResult struct {
		Count  *int64
		Sum    *int64
		Min    *int64
		Max    *int64
		Mode   *int64
		Mean   *float64
		Median *float64
	}
	AggregateNumberResult struct {
		Count  *int64
		Sum    *float64
		Min    *float64
		Max    *float64
		Mode   *float64
		Mean   *float64
		Median *float64
	}
	AggregateBooleanResult struct {
		Count           *int64
		Type            *string
		PercentageTrue  *float64
		PercentageFalse *float64
		TotalTrue       *int64
		TotalFalse      *int64
	}
	AggregateDateResult struct {
		Count  *int64
		Min    *time.Time
		Max    *time.Time
		Mode   *time.Time
		Median *time.Time
	}
	TopOccurrence struct {
		Value       string
		OccursTimes int64
	}
	AggregateGroup struct {
		Property string
		Value    any
		Results  Aggregations
	}
)

func (r *AggregateResponse) UnmarshalMessage(reply *proto.AggregateReply) error {
	var results Aggregations
	var groups []AggregateGroup

	single := reply.GetSingleResult()
	if single != nil {
		aggregations, err := unmarshalAggregations(single.GetAggregations().Aggregations)
		if err != nil {
			return err
		}
		if aggregations != nil {
			results = *aggregations
			results.TotalCount = single.ObjectsCount
		}
	}

	grouped := reply.GetGroupedResults()
	if grouped != nil {
		for _, group := range grouped.Groups {
			by := group.GetGroupedBy()
			property := by.GetPath()[0]

			var results Aggregations
			aggregations, err := unmarshalAggregations(group.GetAggregations().Aggregations)
			if err != nil {
				return err
			}
			if aggregations != nil {
				results = *aggregations
				results.TotalCount = group.ObjectsCount
			}

			var value any
			switch by.GetValue().(type) {
			case *proto.AggregateReply_Group_GroupedBy_Text:
				value = by.GetText()
			case *proto.AggregateReply_Group_GroupedBy_Texts:
				value = by.GetTexts().GetValues()
			case *proto.AggregateReply_Group_GroupedBy_Int:
				value = by.GetInt()
			case *proto.AggregateReply_Group_GroupedBy_Ints:
				value = by.GetInts().GetValues()
			case *proto.AggregateReply_Group_GroupedBy_Number:
				value = by.GetNumber()
			case *proto.AggregateReply_Group_GroupedBy_Numbers:
				value = by.GetNumbers().GetValues()
			case *proto.AggregateReply_Group_GroupedBy_Boolean:
				value = by.GetBoolean()
			case *proto.AggregateReply_Group_GroupedBy_Booleans:
				value = by.GetBooleans().GetValues()
			default:
				// TODO(dyma): support geo
			}

			dev.AssertNotNil(aggregations, "result")
			groups = append(groups, AggregateGroup{
				Property: property,
				Value:    value,
				Results:  results,
			})
		}
	}

	dev.AssertNotNil(results, "results")

	// TODO(dyma): replace took seconds with took time.Duration
	*r = AggregateResponse{
		TookSeconds:    reply.GetTook(),
		Results:        results,
		GroupByResults: groups,
	}
	return nil
}

func unmarshalAggregations(aggregations []*proto.AggregateReply_Aggregations_Aggregation) (*Aggregations, error) {
	out := Aggregations{
		Text:    make(map[string]AggregateTextResult),
		Integer: make(map[string]AggregateIntegerResult),
		Number:  make(map[string]AggregateNumberResult),
		Boolean: make(map[string]AggregateBooleanResult),
		Date:    make(map[string]AggregateDateResult),
	}
	for _, agg := range aggregations {
		property := agg.GetProperty()
		switch {
		case agg.GetText() != nil:
			txt := agg.GetText()
			top := make([]TopOccurrence, len(txt.GetTopOccurences().GetItems()))
			for i, item := range txt.GetTopOccurences().GetItems() {
				top[i] = TopOccurrence{
					Value:       item.Value,
					OccursTimes: item.Occurs,
				}
			}
			out.Text[property] = AggregateTextResult{
				Count:          txt.Count,
				TopOccurrences: top,
			}
		case agg.GetDate() != nil:
			date := agg.GetDate()
			minimum, err := timeFromString(date.GetMinimum())
			if err != nil {
				return nil, fmt.Errorf("%q minimum: %w", property, err)
			}
			maximum, err := timeFromString(date.GetMaximum())
			if err != nil {
				return nil, fmt.Errorf("%q maximum: %w", property, err)
			}
			mode, err := timeFromString(date.GetMode())
			if err != nil {
				return nil, fmt.Errorf("%q mode: %w", property, err)
			}
			median, err := timeFromString(date.GetMedian())
			if err != nil {
				return nil, fmt.Errorf("%q median: %w", property, err)
			}
			out.Date[property] = AggregateDateResult{
				Count:  date.Count,
				Min:    minimum,
				Max:    maximum,
				Mode:   mode,
				Median: median,
			}
		case agg.GetInt() != nil:
			int := agg.GetInt()
			out.Integer[property] = AggregateIntegerResult{
				Count:  int.Count,
				Sum:    int.Sum,
				Min:    int.Minimum,
				Max:    int.Maximum,
				Mode:   int.Mode,
				Median: int.Median,
				Mean:   int.Mean,
			}
		case agg.GetNumber() != nil:
			num := agg.GetNumber()
			out.Number[property] = AggregateNumberResult{
				Count:  num.Count,
				Sum:    num.Sum,
				Min:    num.Minimum,
				Max:    num.Maximum,
				Mode:   num.Mode,
				Median: num.Median,
				Mean:   num.Mean,
			}
		case agg.GetBoolean() != nil:
			bool := agg.GetBoolean()
			out.Boolean[property] = AggregateBooleanResult{
				Count:           bool.Count,
				Type:            bool.Type,
				PercentageTrue:  bool.PercentageTrue,
				PercentageFalse: bool.PercentageFalse,
				TotalTrue:       bool.TotalTrue,
				TotalFalse:      bool.TotalFalse,
			}
		}
	}
	return &out, nil
}

// nilZero returns a pointer to v if it is not the zero value for T and nil otherwise.
func nilZero[T comparable](v T) *T {
	if v == *new(T) {
		return nil
	}
	return &v
}
