package boost_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
	"github.com/weaviate/weaviate-go-client/v6/query/boost"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

//	Conds: []api.BoostCond{
//		{
//			Weight: 3,
//			Func: api.BoostFunc{
//				TimeDecay: &api.TimeDecay{
//					Property: "time_decay_optional",
//					Origin:   testkit.Now,
//					Scale:    120 * time.Hour,
//					Offset:   testkit.Ptr(10 * time.Minute),
//					Curve:    api.BoostCurveGauss,
//					Decay:    testkit.Ptr[float32](92),
//				},
//			},
//		},
//		{
//			Func: api.BoostFunc{
//				TimeDecay: &api.TimeDecay{
//					Property: "time_decay_required",
//					Origin:   testkit.Now,
//					Scale:    120 * time.Hour,
//				},
//			},
//		},
//		{
//			Weight: 3,
//			Func: api.BoostFunc{
//				NumericDecay: &api.NumericDecay{
//					Property: "numeric_decay_optional",
//					Origin:   1,
//					Scale:    2,
//					Offset:   testkit.Ptr[float64](3),
//					Curve:    api.BoostCurveGauss,
//					Decay:    testkit.Ptr[float32](92),
//				},
//			},
//		},
//		{
//			Func: api.BoostFunc{
//				NumericDecay: &api.NumericDecay{
//					Property: "numeric_decay_required",
//					Origin:   1,
//					Scale:    2,
//				},
//			},
//		},
//		{
//			Weight: 3,
//			Func: api.BoostFunc{
//				PropertyValue: &api.PropertyValue{
//					Property: "title_optional",
//					Modifier: api.BoostModifierLOG1P,
//				},
//			},
//		},
//		{
//			Func: api.BoostFunc{
//				PropertyValue: &api.PropertyValue{
//					Property: "title_required",
//					Modifier: api.BoostModifierLOG1P,
//				},
//			},
//		},
//		{
//			Weight: 3,
//			Func: api.BoostFunc{
//				Filter: &api.FilterExpr{
//					Target:   []string{"duration_sec"},
//					Operator: api.FilterOperatorGreaterThan,
//					Value:    92,
//				},
//			},
//		},
//	},

func TestCond(t *testing.T) {
	for _, tt := range testkit.WithOnly(t, []struct {
		testkit.Only

		name string
		cond boost.Cond
		expr api.BoostExpr
	}{
		{
			name: "time decay",
			cond: boost.Cond{
				Weight: 3,
				Func: boost.TimeDecay{
					Property: "created_at",
					Origin:   testkit.Now,
					Scale:    120 * time.Hour,
					Offset:   boost.Offset(10 * time.Minute),
					Curve:    boost.Gauss,
					Decay:    boost.Decay(92),
				},
			},
			expr: api.BoostExpr{
				Weight: 3,
				Conds: []api.BoostCond{
					{
						Weight: 3,
						Func: api.BoostFunc{
							TimeDecay: &api.TimeDecay{
								Property: "created_at",
								Origin:   testkit.Now,
								Scale:    120 * time.Hour,
								Offset:   testkit.Ptr(10 * time.Minute),
								Curve:    api.BoostCurveGauss,
								Decay:    testkit.Ptr[float32](92),
							},
						},
					},
				},
			},
		},
		{
			name: "numeric decay",
			cond: boost.Cond{
				Weight: 3,
				Func: boost.NumericDecay{
					Property: "price",
					Origin:   1,
					Scale:    2,
					Offset:   testkit.Ptr[float64](3),
					Curve:    boost.Gauss,
					Decay:    testkit.Ptr[float32](92),
				},
			},
			expr: api.BoostExpr{
				Weight: 3,
				Conds: []api.BoostCond{
					{
						Weight: 3,
						Func: api.BoostFunc{
							NumericDecay: &api.NumericDecay{
								Property: "price",
								Origin:   1,
								Scale:    2,
								Offset:   testkit.Ptr[float64](3),
								Curve:    api.BoostCurveGauss,
								Decay:    testkit.Ptr[float32](92),
							},
						},
					},
				},
			},
		},
		{
			name: "property value function",
			cond: boost.Cond{
				Weight: 3,
				Func: boost.PropertyValue{
					Property: "title",
					Modifier: boost.Log1P,
				},
			},
			expr: api.BoostExpr{
				Weight: 3,
				Conds: []api.BoostCond{
					{
						Weight: 3,
						Func: api.BoostFunc{
							PropertyValue: &api.PropertyValue{
								Property: "title",
								Modifier: api.BoostModifierLog1P,
							},
						},
					},
				},
			},
		},
		{
			name: "filter",
			cond: boost.Cond{
				Weight: 3,
				Func: boost.Filter{
					Expr: &filter.Cond{
						Target:   "duration_sec",
						Operator: filter.GreaterThan,
						Value:    92,
					},
				},
			},
			expr: api.BoostExpr{
				Weight: 3,
				Conds: []api.BoostCond{
					{
						Weight: 3,
						Func: api.BoostFunc{
							Filter: &api.FilterExpr{
								Target:   []string{"duration_sec"},
								Operator: api.FilterOperatorGreaterThan,
								Value:    92,
							},
						},
					},
				},
			},
		},
	}) {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.cond.Expr()
			require.Equal(t, tt.expr, expr)
		})
	}
}

func TestBlend(t *testing.T) {
	b := boost.Blend{
		Depth:  1,
		Weight: 5,
		Conds: []boost.Cond{
			{Func: boost.PropertyValue{Property: "has_weight"}, Weight: 20},
			{Func: boost.PropertyValue{Property: "no_weight"}},
		},
	}

	want := api.BoostExpr{
		Depth:  1,
		Weight: 5,
		Conds: []api.BoostCond{
			{
				Func:   api.BoostFunc{PropertyValue: &api.PropertyValue{Property: "has_weight"}},
				Weight: 20,
			},
			{
				Func:   api.BoostFunc{PropertyValue: &api.PropertyValue{Property: "no_weight"}},
				Weight: 5,
			},
		},
	}

	require.Implements(t, (*boost.Expr)(nil), b, "blend is a valid boost.Expr")
	expr := b.Expr()

	require.Equal(t, want, expr, "bad blend expr")
}
