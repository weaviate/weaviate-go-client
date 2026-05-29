package filter

import (
	"strings"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
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
	//   - For properties belonging to a referenced collection, it is
	//     the path to that property constructed via [Reference] methods.
	//
	// Example:
	//
	//	"album" 															// Property "album" of the root Songs collection
	//	Reference("performedBy").Property("given_name") 					// Property "given_name" of the referenced Artists collection (single-target reference)
	//	Reference("hasAwards").Collection("GrammyAwards").Property("year")	// Property "year" of the referenced GrammyAwards collection (multi-target reference)
	//	Reference("performedBy").Reference("bornIn").Property("population")	// performedBy -[Artists]-> bornIn -[Cities]-> population
	//
	// See [Len].
	Target string

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
		Target:   split(c.Target),
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

// Len turns the property target a len(property) expression.
// Do not use together with the output of [Reference.Count].
func Len(property string) string {
	path := split(property)
	if l := len(path); l > 0 {
		path[l-1] = "len(" + path[l-1] + ")"
	}
	return join(path...)
}

// Reference describes the path to the target property.
// The path may contain any arbitrary number of "hops",
// each pointing to another collection.
//
// Every path must terminate either in a collection property (see [Reference.Property]),
// or reference-count target (see [Reference.Count]).
type Reference string

// Reference adds another "hop" to the path.
func (r Reference) Reference(reference string) Reference { return Reference(r.Property(reference)) }

// Collection specifies target collection name for a multi-target reference.
func (r Reference) Collection(collection string) Reference {
	path := split(string(r))
	l := len(path)
	dev.Assert(l > 0, "len(reference path)")

	path[l-1] = api.MultiTargetReference(path[l-1], collection)
	return Reference(join(path...))
}

// Property appends a property target at the end of the reference path.
func (r Reference) Property(property string) string { return join(string(r), property) }

// Count converts the reference path into a reference-count target.
func (r Reference) Count() string {
	path := split(string(r))
	l := len(path)
	dev.Assert(l > 0, "len(reference path)")

	path[l-1] = "count(" + path[l-1] + ")"
	return join(path...)
}

// pathSep separates parts of the concatenated Target path.
const pathSep = "/"

func split(s string) []string { return strings.Split(s, pathSep) } // Returned slice is always non-empty.
func join(s ...string) string { return strings.Join(s, pathSep) }
