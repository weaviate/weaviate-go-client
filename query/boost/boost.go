package boost

import (
	"time"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

type (
	Expr interface{ Expr() api.BoostExpr }
	Func interface{ Func() api.BoostFunc }
)

type Cond struct {
	Weight float32
	Func   Func
}

var _ Expr = (*Cond)(nil)

func (c Cond) Expr() api.BoostExpr {
	return api.BoostExpr{
		Weight: c.Weight,
		Conds: []api.BoostCond{
			{
				Func:   c.Func.Func(),
				Weight: c.Weight,
			},
		},
	}
}

type Blend struct {
	Weight float32
	Depth  int
	Conds  []Cond
}

var _ Expr = (*Blend)(nil)

func (b Blend) Expr() api.BoostExpr {
	conds := make([]api.BoostCond, len(b.Conds))
	for i, c := range b.Conds {
		cond := api.BoostCond{
			Func:   c.Func.Func(),
			Weight: c.Weight,
		}
		if c.Weight == 0 {
			cond.Weight = b.Weight
		}
		conds[i] = cond
	}
	return api.BoostExpr{
		Weight: b.Weight,
		Depth:  b.Depth,
		Conds:  conds,
	}
}

func Decay(d float32) *float32                 { return &d }
func Offset[T time.Duration | float64](o T) *T { return &o }

// Day is a [time.Duration] of 24h.
const Day = 24 * time.Hour

type TimeDecay struct {
	Property string        // Required parameter.
	Origin   time.Time     // Required parameter.
	Scale    time.Duration // Required parameter.
	Offset   *time.Duration
	Decay    *float32
	Curve    Curve
}

var _ Func = (*TimeDecay)(nil)

func (td TimeDecay) Func() api.BoostFunc {
	return api.BoostFunc{
		TimeDecay: &api.TimeDecay{
			Property: td.Property,
			Scale:    td.Scale,
			Origin:   td.Origin,
			Offset:   td.Offset,
			Curve:    api.BoostCurve(td.Curve),
			Decay:    td.Decay,
		},
	}
}

type NumericDecay struct {
	Property string  // Required parameter.
	Origin   float64 // Required parameter.
	Scale    float64 // Required parameter.
	Offset   *float64
	Decay    *float32
	Curve    Curve
}

var _ Func = (*NumericDecay)(nil)

func (nd NumericDecay) Func() api.BoostFunc {
	return api.BoostFunc{
		NumericDecay: &api.NumericDecay{
			Property: nd.Property,
			Scale:    nd.Scale,
			Origin:   nd.Origin,
			Offset:   nd.Offset,
			Curve:    api.BoostCurve(nd.Curve),
			Decay:    nd.Decay,
		},
	}
}

type Curve api.BoostCurve

const (
	Gauss       = Curve(api.BoostCurveGauss)
	Linear      = Curve(api.BoostCurveLinear)
	Exponential = Curve(api.BoostCurveExponential)
)

type PropertyValue struct {
	Property string
	Modifier Modifier
}

var _ Func = (*PropertyValue)(nil)

func (pv PropertyValue) Func() api.BoostFunc {
	return api.BoostFunc{
		PropertyValue: &api.PropertyValue{
			Property: pv.Property,
			Modifier: api.BoostModifier(pv.Modifier),
		},
	}
}

type Modifier api.BoostModifier

const (
	Log1P = Modifier(api.BoostModifierLog1P)
	SQRT  = Modifier(api.BoostModifierSQRT)
)

type Filter struct{ filter.Expr }

var _ Func = (*Filter)(nil)

func (f Filter) Func() api.BoostFunc {
	return api.BoostFunc{
		Filter: f.Expr.Expr(),
	}
}
