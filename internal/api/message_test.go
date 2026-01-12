package api_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

type MessageMarshalerTest[In api.RequestMessage, Out api.ReplyMessage] struct {
	name string
	req  api.Message[In, Out]
	want *In
}

// testMessageMarshaler runs test cases for [transports.MessageMarshaler] implementations.
func testMessageMarshaler[In api.RequestMessage, Out api.ReplyMessage](t *testing.T, tests []MessageMarshalerTest[In, Out]) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.req.MarshalMessage()
			require.NoError(t, err)
			require.EqualExportedValues(t, tt.want, got)
		})
	}
}

var (
	singleVector      = []float32{1, 2, 3}
	singleVectorBytes = []byte{
		0x0, 0x0, 0x80, 0x3f,
		0x0, 0x0, 0x0, 0x40,
		0x0, 0x0, 0x40, 0x40,
	}
	multiVector      = [][]float32{{1, 2, 3}, {1, 2, 3}}
	multiVectorBytes = []byte{
		0x3, 0x0, // inner array size, uint16(3)
		0x0, 0x0, 0x80, 0x3f, // first vector
		0x0, 0x0, 0x0, 0x40,
		0x0, 0x0, 0x40, 0x40,
		0x0, 0x0, 0x80, 0x3f, // second vector
		0x0, 0x0, 0x0, 0x40,
		0x0, 0x0, 0x40, 0x40,
	}
)

// TestAggregateRequest_MarshalMessage tests that api.AggregateRequest creates
// the expected proto.AggregateRequest when its MarshalMessage is called.
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
	testMessageMarshaler(t, []MessageMarshalerTest[proto.AggregateRequest, proto.AggregateReply]{
		{
			name: "with limit",
			req:  &api.AggregateRequest{Limit: 10},
			want: &proto.AggregateRequest{Limit: testkit.Ptr[uint32](10)},
		},
		{
			name: "with object limit",
			req:  &api.AggregateRequest{ObjectLimit: 10},
			want: &proto.AggregateRequest{ObjectLimit: testkit.Ptr[uint32](10)},
		},
		{
			name: "text properties",
			req: &api.AggregateRequest{
				Text: map[string]*api.AggregateTextRequest{
					"colour": {Count: true, TopOccurrences: true},
					"tag":    {TopOccurrences: true, TopOccurencesCutoff: 10},
				},
			},
			want: &proto.AggregateRequest{
				Aggregations: sort([]*proto.AggregateRequest_Aggregation{
					{Property: "colour", Aggregation: &proto.AggregateRequest_Aggregation_Text_{
						Text: &proto.AggregateRequest_Aggregation_Text{
							Count: true, TopOccurences: true,
						},
					}},
					{Property: "tag", Aggregation: &proto.AggregateRequest_Aggregation_Text_{
						Text: &proto.AggregateRequest_Aggregation_Text{
							TopOccurences: true, TopOccurencesLimit: testkit.Ptr[uint32](10),
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
		{
			name: "near vector (single vector)",
			req: &api.AggregateRequest{
				NearVector: &api.NearVector{
					Distance: testkit.Ptr(0.5),
					Target: api.SearchTarget{Vectors: []api.TargetVector{
						{Vector: api.Vector{Name: "1d", Single: singleVector}},
					}},
				},
			},
			want: &proto.AggregateRequest{
				Search: &proto.AggregateRequest_NearVector{
					NearVector: &proto.NearVector{
						Distance: testkit.Ptr(0.5),
						Vectors: []*proto.Vectors{{
							Name:        "1d",
							VectorBytes: singleVectorBytes,
							Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
						}},
					},
				},
			},
		},
		{
			name: "near vector (multi vector)",
			req: &api.AggregateRequest{
				NearVector: &api.NearVector{
					Distance: testkit.Ptr(0.5),
					Target: api.SearchTarget{Vectors: []api.TargetVector{
						{Vector: api.Vector{Name: "2d", Multi: multiVector}},
					}},
				},
			},
			want: &proto.AggregateRequest{
				Search: &proto.AggregateRequest_NearVector{
					NearVector: &proto.NearVector{
						Distance: testkit.Ptr(0.5),
						Vectors: []*proto.Vectors{{
							Name:        "2d",
							VectorBytes: multiVectorBytes,
							Type:        proto.Vectors_VECTOR_TYPE_MULTI_FP32,
						}},
					},
				},
			},
		},
		{
			name: "near vector (multi-target average)",
			req: &api.AggregateRequest{
				NearVector: &api.NearVector{
					Distance: testkit.Ptr(0.5),
					Target: api.SearchTarget{
						CombinationMethod: api.CombinationMethodAverage,
						Vectors: []api.TargetVector{
							{Vector: api.Vector{Name: "1d", Single: singleVector}},
							{Vector: api.Vector{Name: "2d", Multi: multiVector}},
						},
					},
				},
			},
			want: &proto.AggregateRequest{
				Search: &proto.AggregateRequest_NearVector{
					NearVector: &proto.NearVector{
						Distance: testkit.Ptr(0.5),
						Targets: &proto.Targets{
							Combination:   proto.CombinationMethod_COMBINATION_METHOD_TYPE_AVERAGE,
							TargetVectors: []string{"1d", "2d"},
						},
						VectorForTargets: []*proto.VectorForTarget{
							{Name: "1d", Vectors: []*proto.Vectors{{
								Name:        "1d",
								VectorBytes: singleVectorBytes,
								Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
							}}},
							{Name: "2d", Vectors: []*proto.Vectors{{
								Name:        "2d",
								VectorBytes: multiVectorBytes,
								Type:        proto.Vectors_VECTOR_TYPE_MULTI_FP32,
							}}},
						},
					},
				},
			},
		},
		{
			name: "near vector (multi-target manual weights)",
			req: &api.AggregateRequest{
				NearVector: &api.NearVector{
					Distance: testkit.Ptr(0.5),
					Target: api.SearchTarget{
						CombinationMethod: api.CombinationMethodManualWeights,
						Vectors: []api.TargetVector{
							{
								Vector: api.Vector{Name: "1d_3", Single: singleVector},
								Weight: testkit.Ptr[float32](.3),
							},
							{
								Vector: api.Vector{Name: "1d_5", Single: singleVector},
								Weight: testkit.Ptr[float32](.5),
							},
							{
								Vector: api.Vector{Name: "2d", Multi: multiVector},
								Weight: testkit.Ptr[float32](.2),
							},
						},
					},
				},
			},
			want: &proto.AggregateRequest{
				Search: &proto.AggregateRequest_NearVector{
					NearVector: &proto.NearVector{
						Distance: testkit.Ptr(0.5),
						Targets: &proto.Targets{
							Combination:   proto.CombinationMethod_COMBINATION_METHOD_TYPE_MANUAL,
							TargetVectors: []string{"1d_3", "1d_5", "2d"},
							WeightsForTargets: []*proto.WeightsForTarget{
								{Target: "1d_3", Weight: .3},
								{Target: "1d_5", Weight: .5},
								{Target: "2d", Weight: .2},
							},
						},
						VectorForTargets: []*proto.VectorForTarget{
							{Name: "1d_3", Vectors: []*proto.Vectors{{
								Name:        "1d_3",
								VectorBytes: singleVectorBytes,
								Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
							}}},
							{Name: "1d_5", Vectors: []*proto.Vectors{{
								Name:        "1d_5",
								VectorBytes: singleVectorBytes,
								Type:        proto.Vectors_VECTOR_TYPE_SINGLE_FP32,
							}}},
							{Name: "2d", Vectors: []*proto.Vectors{{
								Name:        "2d",
								VectorBytes: multiVectorBytes,
								Type:        proto.Vectors_VECTOR_TYPE_MULTI_FP32,
							}}},
						},
					},
				},
			},
		},
	})
}

// ----------------------------------------------------------------------------

type MessageUnmarshalerTest[R api.ReplyMessage, Dest any] struct {
	name      string
	reply     *R
	want      *Dest
	expectErr func(*testing.T, error) // Use require.NoError if nil.
}

// testMessageMarshaler runs test cases for [api.MessageUnmarshaler] implementations.
func testMessageUnmarshaler[R api.ReplyMessage, Dest any](t *testing.T, tests []MessageUnmarshalerTest[R, Dest]) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			var dest any = new(Dest)
			require.Implements(t, (*api.MessageUnmarshaler[R])(nil), dest)

			// Act
			err := dest.(api.MessageUnmarshaler[R]).UnmarshalMessage(tt.reply)

			// Assert
			if tt.expectErr == nil {
				require.NoError(t, err)
			} else {
				tt.expectErr(t, err)
			}
			require.EqualExportedValues(t, tt.want, dest)
		})
	}
}

// TestAggregateRequest_UnmarshalMessage tests that api.AggregateResponse reads
// proto.AggregateRequest correctly when its UnmarshalMessage is called.
func TestAggregateRequest_UnmarshalMessage(t *testing.T) {
	type Aggregations []*proto.AggregateReply_Aggregations_Aggregation

	// result is a helper function to wrap returned aggregations in the layers or protobuf bureaucracy.
	result := func(aggs Aggregations) *proto.AggregateReply_SingleResult {
		return &proto.AggregateReply_SingleResult{
			SingleResult: &proto.AggregateReply_Single{
				ObjectsCount: testkit.Ptr(int64(len(aggs))),
				Aggregations: &proto.AggregateReply_Aggregations{
					Aggregations: aggs,
				},
			},
		}
	}

	// response is a helper to initialize all map fields in api.Aggregations.
	// internal/api should never return nil maps to the caller.
	// To reduce boilerplate in tests, it also populates TotalCount accordingly.
	response := func(aggs api.Aggregations) api.Aggregations {
		out := api.Aggregations{
			Text:    make(map[string]*api.AggregateTextResult),
			Integer: make(map[string]*api.AggregateIntegerResult),
			Number:  make(map[string]*api.AggregateNumberResult),
			Boolean: make(map[string]*api.AggregateBooleanResult),
			Date:    make(map[string]*api.AggregateDateResult),
		}
		switch {
		case aggs.Text != nil:
			out.TotalCount = testkit.Ptr(int64(len(aggs.Text)))
			out.Text = aggs.Text
		case aggs.Integer != nil:
			out.TotalCount = testkit.Ptr(int64(len(aggs.Integer)))
			out.Integer = aggs.Integer
		case aggs.Number != nil:
			out.TotalCount = testkit.Ptr(int64(len(aggs.Number)))
			out.Number = aggs.Number
		case aggs.Boolean != nil:
			out.TotalCount = testkit.Ptr(int64(len(aggs.Boolean)))
			out.Boolean = aggs.Boolean
		case aggs.Date != nil:
			out.TotalCount = testkit.Ptr(int64(len(aggs.Date)))
			out.Date = aggs.Date
		}
		return out
	}

	testMessageUnmarshaler(t, []MessageUnmarshalerTest[proto.AggregateReply, api.AggregateResponse]{
		{
			name: "text properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: result(Aggregations{
					{Property: "colour", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Text_{
						Text: &proto.AggregateReply_Aggregations_Aggregation_Text{
							Count: testkit.Ptr[int64](1),
							TopOccurences: &proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences{
								Items: []*proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences_TopOccurrence{
									{Value: "red", Occurs: 2},
									{Value: "blue", Occurs: 3},
								},
							},
						},
					}},
					{Property: "tag", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Text_{
						Text: &proto.AggregateReply_Aggregations_Aggregation_Text{
							TopOccurences: &proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences{
								Items: []*proto.AggregateReply_Aggregations_Aggregation_Text_TopOccurrences_TopOccurrence{
									{Value: "casual", Occurs: 1},
									{Value: "comfy", Occurs: 2},
								},
							},
						},
					}},
				}),
			},
			want: &api.AggregateResponse{
				TookSeconds: 92,
				Results: response(api.Aggregations{
					Text: map[string]*api.AggregateTextResult{
						"colour": {
							Count: testkit.Ptr[int64](1),
							TopOccurences: []*api.TopOccurence{
								{Value: "red", Occurs: 2},
								{Value: "blue", Occurs: 3},
							},
						},
						"tag": {
							TopOccurences: []*api.TopOccurence{
								{Value: "casual", Occurs: 1},
								{Value: "comfy", Occurs: 2},
							},
						},
					},
				}),
			},
		},
		{
			name: "integer properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: result(Aggregations{
					{Property: "price", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Int{
						Int: &proto.AggregateReply_Aggregations_Aggregation_Integer{
							Sum:     testkit.Ptr[int64](1),
							Minimum: testkit.Ptr[int64](2),
							Maximum: testkit.Ptr[int64](3),
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Int{
						Int: &proto.AggregateReply_Aggregations_Aggregation_Integer{
							Count:  testkit.Ptr[int64](1),
							Mode:   testkit.Ptr[int64](2),
							Median: testkit.Ptr[float64](3),
						},
					}},
				}),
			},
			want: &api.AggregateResponse{
				TookSeconds: 92,
				Results: response(api.Aggregations{
					Integer: map[string]*api.AggregateIntegerResult{
						"price": {
							Sum:     testkit.Ptr[int64](1),
							Minimum: testkit.Ptr[int64](2),
							Maximum: testkit.Ptr[int64](3),
						},
						"size": {
							Count:  testkit.Ptr[int64](1),
							Mode:   testkit.Ptr[int64](2),
							Median: testkit.Ptr[float64](3),
						},
					},
				}),
			},
		},
		{
			name: "number properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: result(Aggregations{
					{Property: "price", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Number_{
						Number: &proto.AggregateReply_Aggregations_Aggregation_Number{
							Sum:     testkit.Ptr[float64](1),
							Minimum: testkit.Ptr[float64](2),
							Maximum: testkit.Ptr[float64](3),
						},
					}},
					{Property: "size", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Number_{
						Number: &proto.AggregateReply_Aggregations_Aggregation_Number{
							Count:  testkit.Ptr[int64](1),
							Mode:   testkit.Ptr[float64](2),
							Median: testkit.Ptr[float64](3),
						},
					}},
				}),
			},
			want: &api.AggregateResponse{
				TookSeconds: 92,
				Results: response(api.Aggregations{
					Number: map[string]*api.AggregateNumberResult{
						"price": {
							Sum:     testkit.Ptr[float64](1),
							Minimum: testkit.Ptr[float64](2),
							Maximum: testkit.Ptr[float64](3),
						},
						"size": {
							Count:  testkit.Ptr[int64](1),
							Mode:   testkit.Ptr[float64](2),
							Median: testkit.Ptr[float64](3),
						},
					},
				}),
			},
		},
		{
			name: "boolean properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: result(Aggregations{
					{Property: "onSale", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Boolean_{
						Boolean: &proto.AggregateReply_Aggregations_Aggregation_Boolean{
							PercentageTrue:  testkit.Ptr[float64](1),
							PercentageFalse: testkit.Ptr[float64](2),
						},
					}},
					{Property: "newArrival", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Boolean_{
						Boolean: &proto.AggregateReply_Aggregations_Aggregation_Boolean{
							Count:      testkit.Ptr[int64](1),
							TotalTrue:  testkit.Ptr[int64](2),
							TotalFalse: testkit.Ptr[int64](3),
						},
					}},
				}),
			},
			want: &api.AggregateResponse{
				TookSeconds: 92,
				Results: response(api.Aggregations{
					Boolean: map[string]*api.AggregateBooleanResult{
						"onSale": {
							PercentageTrue:  testkit.Ptr[float64](1),
							PercentageFalse: testkit.Ptr[float64](2),
						},
						"newArrival": {
							Count:      testkit.Ptr[int64](1),
							TotalTrue:  testkit.Ptr[int64](2),
							TotalFalse: testkit.Ptr[int64](3),
						},
					},
				}),
			},
		},
		{
			name: "date properties",
			reply: &proto.AggregateReply{
				Took: 92,
				Result: result(Aggregations{
					{Property: "lastPurchase", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Date_{
						Date: &proto.AggregateReply_Aggregations_Aggregation_Date{
							Count:   testkit.Ptr[int64](1),
							Minimum: testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
							Maximum: testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
						},
					}},
					{Property: "lastReturn", Aggregation: &proto.AggregateReply_Aggregations_Aggregation_Date_{
						Date: &proto.AggregateReply_Aggregations_Aggregation_Date{
							Mode:   testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
							Median: testkit.Ptr(testkit.Now.Format(time.RFC3339Nano)),
						},
					}},
				}),
			},
			want: &api.AggregateResponse{
				TookSeconds: 92,
				Results: response(api.Aggregations{
					Date: map[string]*api.AggregateDateResult{
						"lastPurchase": {
							Count:   testkit.Ptr[int64](1),
							Minimum: &testkit.Now,
							Maximum: &testkit.Now,
						},
						"lastReturn": {
							Mode:   &testkit.Now,
							Median: &testkit.Now,
						},
					},
				}),
			},
		},
	})
}
