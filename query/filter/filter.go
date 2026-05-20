package filter

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

// Expr describes a filter expression. Expressions can be combined into [And] and [Or] groups.
type Expr interface {
	Expr() *api.FilterExpr
}

var (
	_ Expr = (*Cond)(nil)
	_ Expr = (*And)(nil)
	_ Expr = (*Or)(nil)
	_ Expr = (*Not)(nil)
)

type Operator api.FilterOperator

const (
	Equal            = Operator(api.FilterOperatorEqual)
	LessThan         = Operator(api.FilterOperatorLessThan)
	LessThanEqual    = Operator(api.FilterOperatorLessThanEqual)
	GreaterThan      = Operator(api.FilterOperatorGreaterThan)
	GreaterThanEqual = Operator(api.FilterOperatorGreaterThanEqual)
	Like             = Operator(api.FilterOperatorLike)
	IsNull           = Operator(api.FilterOperatorIsNull)
	ContainsAll      = Operator(api.FilterOperatorContainsAll)
	ContainsAny      = Operator(api.FilterOperatorContainsAny)
	ContainsNone     = Operator(api.FilterOperatorContainsNone)
)

type Cond struct {
	// Comparison operator.
	// See [Value] for discussion of values supported by different operators.
	Operator Operator

	// Target describes the left-hand side of the filter expression.
	// It is representsed as a path to the filter property:
	//
	//   - For regular properties, it is the name of the target property.
	//   - For properties belonging to a single-target reference, it is
	//     then name of the reference property followed by the name of
	//     the target property, arbitrarily nested.
	//   - For properties belonging to a multi-target reference, it is
	//     the name of the reference property followed by the name of
	//     the target collection and the property therein, arbitrarily nested.
	//
	// Example:
	//
	//	[]string{"album"} 									// Property "album" of the root Songs collection
	//	[]string{"performedBy", "given_name"} 				// Property "given_name" of the referenced Artists collection (single-target reference)
	//	[]string{"hasAwards", "GrammyAward", "year"} 		// Property "year" of the referenced Grammy collection (multi-target reference)
	//	[]string{"performedBy", "bornIn", "population"}  	// performedBy -[Artists]-> bornIn -[Cities]-> population
	//
	// See [Len] and [ReferenceCount].
	Target []string

	// Value represents the right-hand side of the filter expression.
	// Its type should match the data type of the [Target] property.
	//
	//	- [Like] operator requires a regular expression string.
	//	- [IsNull] operator does not need a value.
	//	- Contains- operators require a slice of values.
	// 	- Date properties should be compared to instances of [time.Time].
	Value any
}

func (c *Cond) Expr() *api.FilterExpr {
	if c == nil {
		return nil
	}
	return &api.FilterExpr{
		Operator: api.FilterOperator(c.Operator),
		Target:   c.Target,
		Value:    c.Value,
	}
}

// And is a group of sub-expressions joined with the AND operator.
// It can be the top level exression or combined with other sub-expressions.
type And []Expr

func (and And) Expr() *api.FilterExpr {
	if len(and) == 0 {
		return nil
	}
	exprs := make([]api.FilterExpr, 0, len(and))
	for _, ex := range and {
		if ex != nil {
			if expr := ex.Expr(); expr != nil {
				exprs = append(exprs, *expr)
			}
		}
	}
	if len(exprs) == 0 {
		return nil
	}
	return &api.FilterExpr{
		Operator: api.FilterOperatorAnd,
		Exprs:    exprs,
	}
}

// Or is a group of sub-expressions joined with the OR operator.
// It can be the top level exression or combined with other sub-expressions.
type Or []Expr

func (or Or) Expr() *api.FilterExpr {
	if len(or) == 0 {
		return nil
	}
	exprs := make([]api.FilterExpr, 0, len(or))
	for _, ex := range or {
		if ex != nil {
			if expr := ex.Expr(); expr != nil {
				exprs = append(exprs, *expr)
			}
		}
	}
	if len(exprs) == 0 {
		return nil
	}
	return &api.FilterExpr{
		Operator: api.FilterOperatorOr,
		Exprs:    exprs,
	}
}

// Not negates the expression.
type Not struct{ P Expr }

func (not Not) Expr() *api.FilterExpr {
	if not.P == nil {
		return nil
	}
	expr := not.P.Expr()
	if expr == nil {
		return nil
	}
	return &api.FilterExpr{
		Operator: api.FilterOperatorNot,
		Exprs:    []api.FilterExpr{*expr},
	}
}

var (
	UUID          = api.FieldUUID          // UUID property name.
	CreatedAt     = api.FieldCreatedAt     // Creation time property name.
	LastUpdatedAt = api.FieldLastUpdatedAt // Last update time property name.
)

// Len turns the last element of the property path into a len(property) expression.
func Len(path ...string) []string {
	if l := len(path); l > 0 {
		path[l-1] = "len(" + path[l-1] + ")"
	}
	return path
}

// ReferenceCount turns the last element of the property path into reference-count target.
func ReferenceCount(path ...string) []string {
	if l := len(path); l > 0 {
		path[l-1] = "count(" + path[l-1] + ")"
	}
	return path
}
