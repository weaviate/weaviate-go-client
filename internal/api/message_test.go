package api_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

type MessageMarshalerTest[R transport.RequestMessage] struct {
	name string
	req  transport.MessageMarshaler[R]
	want *R
}

// TestAggregateRequest_MarshalMessage tests that api.AggregateRequest creates
// the expected proto.AggregateRequest when it's MarshalMessage is called.
// Most of the test cases are split by property type to simplify error logs;
// there's nothing stopping a test case from having mixed property types though.
func TestAggregateRequest_MarshalMessage(t *testing.T) {
	// Sort orders expected aggregations by their property name to match
	// the ordering we expect MarshalMessage to apply.
	sort := func(in []*proto.AggregateRequest_Aggregation) []*proto.AggregateRequest_Aggregation {
		slices.SortFunc(in, func(a, b *proto.AggregateRequest_Aggregation) int {
			return strings.Compare(a.Property, b.Property)
		})
		return in
	}
	testMessageMarshaler(t, []MessageMarshalerTest[proto.AggregateRequest]{
		{
			name: "with limit",
			req:  &api.AggregateRequest{Limit: 10},
			want: &proto.AggregateRequest{Limit: testkit.Ptr(uint32(10))},
		},
		{
			name: "with object limit",
			req:  &api.AggregateRequest{ObjectLimit: 10},
			want: &proto.AggregateRequest{ObjectLimit: testkit.Ptr(uint32(10))},
		},
		{
			name: "text properties",
			req: &api.AggregateRequest{
				Text: map[string]*api.AggregateTextRequest{
					"colour": {Count: true, TopOccurrences: true},
					"tags":   {TopOccurrences: true, TopOccurencesCutoff: 10},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: sort([]*proto.AggregateRequest_Aggregation{
					{Property: "colour", Aggregation: &proto.AggregateRequest_Aggregation_Text_{
						Text: &proto.AggregateRequest_Aggregation_Text{
							Count: true, TopOccurences: true,
						},
					}},
					{Property: "tags", Aggregation: &proto.AggregateRequest_Aggregation_Text_{
						Text: &proto.AggregateRequest_Aggregation_Text{
							TopOccurences: true, TopOccurencesLimit: testkit.Ptr(uint32(10)),
						},
					}},
				}),
			},
		},
		{
			name: "integer properties",
			req: &api.AggregateRequest{
				Integer: map[string]*api.AggregateIntegerRequest{
					"price": {Sum: true, Minimum: true, Maximum: true},
					"size":  {Count: true, Mode: true, Median: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: sort([]*proto.AggregateRequest_Aggregation{
					{Property: "price", Aggregation: &proto.AggregateRequest_Aggregation_Int{
						Int: &proto.AggregateRequest_Aggregation_Integer{
							Sum: true, Minimum: true, Maximum: true,
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateRequest_Aggregation_Int{
						Int: &proto.AggregateRequest_Aggregation_Integer{
							Count: true, Mode: true, Median: true,
						},
					}},
				}),
			},
		},
		{
			name: "number properties",
			req: &api.AggregateRequest{
				Number: map[string]*api.AggregateNumberRequest{
					"price": {Sum: true, Minimum: true, Maximum: true},
					"size":  {Count: true, Mode: true, Median: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: sort([]*proto.AggregateRequest_Aggregation{
					{Property: "price", Aggregation: &proto.AggregateRequest_Aggregation_Number_{
						Number: &proto.AggregateRequest_Aggregation_Number{
							Sum: true, Minimum: true, Maximum: true,
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateRequest_Aggregation_Number_{
						Number: &proto.AggregateRequest_Aggregation_Number{
							Count: true, Mode: true, Median: true,
						},
					}},
				}),
			},
		},
		{
			name: "boolean properties",
			req: &api.AggregateRequest{
				Boolean: map[string]*api.AggregateBooleanRequest{
					"onSale":     {PercentageTrue: true, PercentageFalse: true},
					"newArrival": {Count: true, TotalTrue: true, TotalFalse: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: sort([]*proto.AggregateRequest_Aggregation{
					{Property: "onSale", Aggregation: &proto.AggregateRequest_Aggregation_Boolean_{
						Boolean: &proto.AggregateRequest_Aggregation_Boolean{
							PercentageTrue: true, PercentageFalse: true,
						},
					}},
					{Property: "newArrival", Aggregation: &proto.AggregateRequest_Aggregation_Boolean_{
						Boolean: &proto.AggregateRequest_Aggregation_Boolean{
							Count: true, TotalTrue: true, TotalFalse: true,
						},
					}},
				}),
			},
		},
		{
			name: "date properties",
			req: &api.AggregateRequest{
				Date: map[string]*api.AggregateDateRequest{
					"lastPurchase": {Count: true, Minimum: true, Maximum: true},
					"lastReturn":   {Mode: true, Median: true},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: sort([]*proto.AggregateRequest_Aggregation{
					{Property: "lastPurchase", Aggregation: &proto.AggregateRequest_Aggregation_Date_{
						Date: &proto.AggregateRequest_Aggregation_Date{
							Count: true, Minimum: true, Maximum: true,
						},
					}},
					{Property: "lastReturn", Aggregation: &proto.AggregateRequest_Aggregation_Date_{
						Date: &proto.AggregateRequest_Aggregation_Date{
							Mode: true, Median: true,
						},
					}},
				}),
			},
		},
	})
}

func testMessageMarshaler[R transport.RequestMessage](t *testing.T, tests []MessageMarshalerTest[R]) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.MarshalMessage()
			require.Equal(t, tt.want, got)
		})
	}
}
