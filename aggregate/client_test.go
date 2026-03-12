package aggregate_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/aggregate"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestNewClient(t *testing.T) {
	require.Panics(t, func() {
		aggregate.NewClient(nil, api.RequestDefaults{CollectionName: "New"})
	}, "nil transport")
}

func TestClient_OverAll(t *testing.T) {
	rd := api.RequestDefaults{
		CollectionName:   "Songs",
		Tenant:           "john_doe",
		ConsistencyLevel: api.ConsistencyLevelQuorum,
	}

	// response is a helper to initialize all map fields in api.Aggregations.
	// internal/api should never return nil maps to the caller.
	// To reduce boilerplate in tests, it also populates TotalCount accordingly.
	response := func(aggs api.Aggregations) api.Aggregations {
		results := api.Aggregations{
			TotalCount: aggs.TotalCount,
			Text:       make(map[string]api.AggregateTextResult),
			Integer:    make(map[string]api.AggregateIntegerResult),
			Number:     make(map[string]api.AggregateNumberResult),
			Boolean:    make(map[string]api.AggregateBooleanResult),
			Date:       make(map[string]api.AggregateDateResult),
		}
		switch {
		case aggs.Text != nil:
			results.Text = aggs.Text
		case aggs.Integer != nil:
			results.Integer = aggs.Integer
		case aggs.Number != nil:
			results.Number = aggs.Number
		case aggs.Boolean != nil:
			results.Boolean = aggs.Boolean
		case aggs.Date != nil:
			results.Date = aggs.Date
		}
		return results
	}

	// result is a helper to initialize all map fields in aggregation.Results.
	// internal/api should never return nil maps to the caller.
	// To reduce boilerplate in tests, it also populates TotalCount accordingly.
	result := func(aggs aggregate.Aggregations) aggregate.Aggregations {
		results := aggregate.Aggregations{
			TotalCount: aggs.TotalCount,
			Text:       make(map[string]aggregate.TextResult),
			Integer:    make(map[string]aggregate.IntegerResult),
			Number:     make(map[string]aggregate.NumberResult),
			Boolean:    make(map[string]aggregate.BooleanResult),
			Date:       make(map[string]aggregate.DateResult),
		}
		switch {
		case aggs.Text != nil:
			results.Text = aggs.Text
		case aggs.Integer != nil:
			results.Integer = aggs.Integer
		case aggs.Number != nil:
			results.Number = aggs.Number
		case aggs.Boolean != nil:
			results.Boolean = aggs.Boolean
		case aggs.Date != nil:
			results.Date = aggs.Date
		}
		return results
	}

	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name    string
		request aggregate.OverAll // Aggregate query parameters.
		stubs   []testkit.Stub[api.AggregateRequest, api.AggregateResponse]
		want    *aggregate.Result // Expected return value.
		err     testkit.Error
	}{
		{
			name: "text properties",
			request: aggregate.OverAll{
				ObjectLimit: 2,
				Text: []aggregate.Text{
					{Property: "colour", Count: true, TopOccurrences: true},
					{Property: "tag", TopOccurrences: true, TopOccurencesCutoff: 10},
				},
			},
			stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
				{
					Request: &api.AggregateRequest{
						RequestDefaults: rd,
						ObjectLimit:     2,
						Text: []api.AggregateTextRequest{
							{Property: "colour", Count: true, TopOccurrences: true},
							{Property: "tag", TopOccurrences: true, TopOccurencesCutoff: 10},
						},
					},
					Response: api.AggregateResponse{
						TookSeconds: 92,
						Results: response(api.Aggregations{
							TotalCount: testkit.Ptr[int64](2),
							Text: map[string]api.AggregateTextResult{
								"colour": {
									Count: testkit.Ptr[int64](1),
									TopOccurrences: []api.TopOccurrence{
										{Value: "red", OccursTimes: 2},
										{Value: "blue", OccursTimes: 3},
									},
								},
								"tag": {
									TopOccurrences: []api.TopOccurrence{
										{Value: "casual", OccursTimes: 1},
										{Value: "comfy", OccursTimes: 2},
									},
								},
							},
						}),
					},
				},
			},
			want: &aggregate.Result{
				TookSeconds: 92,
				Aggregations: result(aggregate.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Text: map[string]aggregate.TextResult{
						"colour": {
							Count: testkit.Ptr[int64](1),
							TopOccurrences: []aggregate.TopOccurrence{
								{Value: "red", OccursTimes: 2},
								{Value: "blue", OccursTimes: 3},
							},
						},
						"tag": {
							TopOccurrences: []aggregate.TopOccurrence{
								{Value: "casual", OccursTimes: 1},
								{Value: "comfy", OccursTimes: 2},
							},
						},
					},
				}),
			},
			Only: true,
		},
		{
			name: "integer properties",
			request: aggregate.OverAll{
				Integer: []aggregate.Integer{
					{Property: "price", Sum: true, Min: true, Max: true},
					{Property: "size", Count: true, Mode: true, Median: true},
				},
			},
			stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
				{
					Request: &api.AggregateRequest{
						Integer: []api.AggregateIntegerRequest{
							{Property: "price", Sum: true, Min: true, Max: true},
							{Property: "size", Count: true, Mode: true, Median: true},
						},
					},
					Response: api.AggregateResponse{
						TookSeconds: 92,
						Results: response(api.Aggregations{
							TotalCount: testkit.Ptr[int64](2),
							Integer: map[string]api.AggregateIntegerResult{
								"price": {
									Sum: testkit.Ptr[int64](1),
									Min: testkit.Ptr[int64](2),
									Max: testkit.Ptr[int64](3),
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
			},
			want: &aggregate.Result{
				TookSeconds: 92,
				Aggregations: result(aggregate.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Integer: map[string]aggregate.IntegerResult{
						"price": {
							Sum: testkit.Ptr[int64](1),
							Min: testkit.Ptr[int64](2),
							Max: testkit.Ptr[int64](3),
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
			request: aggregate.OverAll{
				Number: []aggregate.Number{
					{Property: "price", Sum: true, Min: true, Max: true},
					{Property: "size", Count: true, Mode: true, Median: true},
				},
			},
			stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
				{
					Request: &api.AggregateRequest{
						Number: []api.AggregateNumberRequest{
							{Property: "price", Sum: true, Min: true, Max: true},
							{Property: "size", Count: true, Mode: true, Median: true},
						},
					},
					Response: api.AggregateResponse{
						TookSeconds: 92,
						Results: response(api.Aggregations{
							TotalCount: testkit.Ptr[int64](2),
							Number: map[string]api.AggregateNumberResult{
								"price": {
									Sum: testkit.Ptr[float64](1),
									Min: testkit.Ptr[float64](2),
									Max: testkit.Ptr[float64](3),
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
			},
			want: &aggregate.Result{
				TookSeconds: 92,
				Aggregations: result(aggregate.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Number: map[string]aggregate.NumberResult{
						"price": {
							Sum: testkit.Ptr[float64](1),
							Min: testkit.Ptr[float64](2),
							Max: testkit.Ptr[float64](3),
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
			request: aggregate.OverAll{
				Boolean: []aggregate.Boolean{
					{Property: "onSale", Type: true, PercentageTrue: true, PercentageFalse: true},
					{Property: "newArrival", Count: true, TotalTrue: true, TotalFalse: true},
				},
			},
			stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
				{
					Request: &api.AggregateRequest{
						Boolean: []api.AggregateBooleanRequest{
							{Property: "onSale", Type: true, PercentageTrue: true, PercentageFalse: true},
							{Property: "newArrival", Count: true, TotalTrue: true, TotalFalse: true},
						},
					},
					Response: api.AggregateResponse{
						TookSeconds: 92,
						Results: response(api.Aggregations{
							TotalCount: testkit.Ptr[int64](2),
							Boolean: map[string]api.AggregateBooleanResult{
								"onSale": {
									Type:            testkit.Ptr("black_friday"),
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
			},
			want: &aggregate.Result{
				TookSeconds: 92,
				Aggregations: result(aggregate.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Boolean: map[string]aggregate.BooleanResult{
						"onSale": {
							Type:            testkit.Ptr("black_friday"),
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
			request: aggregate.OverAll{
				Date: []aggregate.Date{
					{Property: "lastPurchase", Count: true, Min: true, Max: true},
					{Property: "lastReturn", Mode: true, Median: true},
				},
			},
			stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
				{
					Request: &api.AggregateRequest{
						Date: []api.AggregateDateRequest{
							{Property: "lastPurchase", Count: true, Min: true, Max: true},
							{Property: "lastReturn", Mode: true, Median: true},
						},
					},
					Response: api.AggregateResponse{
						TookSeconds: 92,
						Results: response(api.Aggregations{
							TotalCount: testkit.Ptr[int64](2),
							Date: map[string]api.AggregateDateResult{
								"lastPurchase": {
									Count: testkit.Ptr[int64](1),
									Min:   &testkit.Now,
									Max:   &testkit.Now,
								},
								"lastReturn": {
									Mode:   &testkit.Now,
									Median: &testkit.Now,
								},
							},
						}),
					},
				},
			},
			want: &aggregate.Result{
				TookSeconds: 92,
				Aggregations: result(aggregate.Aggregations{
					TotalCount: testkit.Ptr[int64](2),
					Date: map[string]aggregate.DateResult{
						"lastPurchase": {
							Count: testkit.Ptr[int64](1),
							Min:   &testkit.Now,
							Max:   &testkit.Now,
						},
						"lastReturn": {
							Mode:   &testkit.Now,
							Median: &testkit.Now,
						},
					},
				}),
			},
		},
		{
			name: "request error",
			stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
				{Err: testkit.ErrWhaam},
			},
			err: testkit.ExpectError,
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			transport := testkit.NewTransport(t, tt.stubs)
			c := aggregate.NewClient(transport, rd)
			require.NotNil(t, c, "nil client")

			got, err := c.OverAll(t.Context(), tt.request)

			tt.err.Require(t, err, "request error")
			require.EqualExportedValues(t, tt.want, got, "bad result")
		})
	}

	t.Run("group by", func(t *testing.T) {
		for _, tt := range testkit.WithOnly(t, []struct {
			testkit.Only

			name    string
			request aggregate.OverAll // Aggregate query parameters.
			groupBy aggregate.GroupBy // GroupBy clause.
			stubs   []testkit.Stub[api.AggregateRequest, api.AggregateResponse]
			want    *aggregate.GroupByResult // Expected return value.
			err     testkit.Error
		}{
			{
				name: "request ok",
				request: aggregate.OverAll{
					TotalCount: true,
					Boolean: []aggregate.Boolean{
						{Property: "onSale", PercentageTrue: true, PercentageFalse: true},
					},
					Number: []aggregate.Number{
						{Property: "onSale", Sum: true, Min: true, Max: true},
					},
				},
				groupBy: aggregate.GroupBy{Property: "album", Limit: 1},
				stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
					{
						Request: &api.AggregateRequest{
							RequestDefaults: rd,
							TotalCount:      true,
							Boolean: []api.AggregateBooleanRequest{
								{Property: "onSale", PercentageTrue: true, PercentageFalse: true},
							},
							Number: []api.AggregateNumberRequest{
								{Property: "onSale", Sum: true, Min: true, Max: true},
							},
							GroupBy: &api.GroupBy{
								Property: "album",
								Limit:    1,
							},
						},
						Response: api.AggregateResponse{
							TookSeconds: 92,
							GroupByResults: []api.AggregateGroup{
								{
									Property: "onSale",
									Value:    true,
									Results: api.Aggregations{
										TotalCount: testkit.Ptr(int64(1)),
										Boolean: map[string]api.AggregateBooleanResult{
											"onSale": {
												Type:            testkit.Ptr("black_friday"),
												PercentageTrue:  testkit.Ptr[float64](1),
												PercentageFalse: testkit.Ptr[float64](2),
											},
										},
									},
								},
								{
									Property: "price",
									Value:    float64(4),
									Results: response(api.Aggregations{
										TotalCount: testkit.Ptr[int64](1),
										Number: map[string]api.AggregateNumberResult{
											"price": {
												Sum: testkit.Ptr[float64](1),
												Min: testkit.Ptr[float64](2),
												Max: testkit.Ptr[float64](3),
											},
										},
									}),
								},
							},
						},
					},
				},
				want: &aggregate.GroupByResult{
					Groups: []aggregate.Group{
						{
							Property: "onSale",
							Value:    true,
							Aggregations: result(aggregate.Aggregations{
								TotalCount: testkit.Ptr(int64(1)),
								Boolean: map[string]aggregate.BooleanResult{
									"onSale": {
										Type:            testkit.Ptr("black_friday"),
										PercentageTrue:  testkit.Ptr[float64](1),
										PercentageFalse: testkit.Ptr[float64](2),
									},
								},
							}),
						},
						{
							Property: "price",
							Value:    float64(4),
							Aggregations: result(aggregate.Aggregations{
								TotalCount: testkit.Ptr[int64](1),
								Number: map[string]aggregate.NumberResult{
									"price": {
										Sum: testkit.Ptr[float64](1),
										Min: testkit.Ptr[float64](2),
										Max: testkit.Ptr[float64](3),
									},
								},
							}),
						},
					},
				},
			},
			{
				name: "request error",
				stubs: []testkit.Stub[api.AggregateRequest, api.AggregateResponse]{
					{Err: testkit.ErrWhaam},
				},
				err: testkit.ExpectError,
			},
		}) {
			t.Run(tt.name, func(t *testing.T) {
				transport := testkit.NewTransport(t, tt.stubs)
				c := aggregate.NewClient(transport, rd)
				require.NotNil(t, c, "nil client")

				got, err := c.OverAll.GroupBy(t.Context(), tt.request, tt.groupBy)

				tt.err.Require(t, err, "request error")
				require.EqualExportedValues(t, tt.want, got, "bad result")
			})
		}
	})
}
