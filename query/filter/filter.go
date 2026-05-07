package filter

import (
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

// Expr describes a filter exression. Expressions can be combined into [And] and [Or] groups.
type Expr interface {
	Operator() api.FilterOperator
	Exprs() []Expr // For AND/OR

	Target() []string
	Value() any
}

// Property constructs the left-hand side of the filter expression.
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
//	filter.Property("album") // Property "album" of the root Songs collection
//	filter.Property("performedBy", "given_name") // Property "given_name" of the referenced Artists collection (single-target reference)
//	filter.Property("hasAwards", "GrammyAward", "year") // Property "year" of the referenced Grammy collection (multi-target reference)
//	filter.Property("performedBy", "bornIn", "population")  // performedBy -[Artists]-> bornIn -[Cities]-> population
func Property[T any](path ...string) Target[T] { return target[T](path) }

var (
	UUID          = Property[uuid.UUID](api.FieldUUID)
	CreatedAt     = Property[time.Time](api.FieldCreatedAt)
	LastUpdatedAt = Property[time.Time](api.FieldLastUpdatedAt)
)

// Len turns the last element of the property path into a len(property) expression.
func Len(path ...string) Target[int64] {
	if l := len(path); l > 0 {
		path[l-1] = "len(" + path[l-1] + ")"
	}
	return Property[int64](path...)
}

// ReferenceCount turns the last element of the property path into reference-count target.
func ReferenceCount(path ...string) Target[int64] {
	if l := len(path); l > 0 {
		path[l-1] = "count(" + path[l-1] + ")"
	}
	return Property[int64](path...)
}

// Target is the left-hand side of the filter expression.
type Target[T any] interface {
	Equal(T) Expr
	LessThan(T) Expr
	LessThanEqual(T) Expr
	GreaterThan(T) Expr
	GreaterThanEqual(T) Expr
	Like(string) Expr
	Null() Expr
	ContainsAny(...T) Expr
	ContainsAll(...T) Expr
	ContainsNone(...T) Expr
}

// Not negates the expression.
func Not(e Expr) Expr {
	return &expr{operator: api.FilterOperatorNot, exprs: []Expr{e}}
}

// And is a group of sub-expressions joined with the AND operator.
// It can be the top level exression or combined with other sub-expressions.
type And []Expr

var _ Expr = (And)(nil)

func (and And) Exprs() []Expr                { return and }
func (and And) Operator() api.FilterOperator { return api.FilterOperatorAnd }
func (and And) Target() []string             { return nil }
func (and And) Value() any                   { return nil }

// Or is a group of sub-expressions joined with the OR operator.
// It can be the top level exression or combined with other sub-expressions.
type Or []Expr

var _ Expr = (And)(nil)

func (or Or) Exprs() []Expr                { return or }
func (or Or) Operator() api.FilterOperator { return api.FilterOperatorOr }
func (or Or) Target() []string             { return nil }
func (or Or) Value() any                   { return nil }

// target implements [Target] for a property path.
type target[T any] []string

var _ Target[any] = (*target[any])(nil)

func (t target[T]) Equal(v T) Expr            { return t.expr(api.FilterOperatorEqual, v) }
func (t target[T]) LessThan(v T) Expr         { return t.expr(api.FilterOperatorLessThan, v) }
func (t target[T]) LessThanEqual(v T) Expr    { return t.expr(api.FilterOperatorLessThanEqual, v) }
func (t target[T]) GreaterThan(v T) Expr      { return t.expr(api.FilterOperatorGreaterThan, v) }
func (t target[T]) GreaterThanEqual(v T) Expr { return t.expr(api.FilterOperatorGreaterThanEqual, v) }
func (t target[T]) Like(v string) Expr        { return t.expr(api.FilterOperatorLike, v) }
func (t target[T]) Null() Expr                { return t.expr(api.FilterOperatorIsNull, nil) }
func (t target[T]) ContainsAll(vs ...T) Expr  { return t.expr(api.FilterOperatorContainsAll, vs) }
func (t target[T]) ContainsAny(vs ...T) Expr  { return t.expr(api.FilterOperatorContainsAny, vs) }
func (t target[T]) ContainsNone(vs ...T) Expr { return t.expr(api.FilterOperatorContainsNone, vs) }

func (t target[T]) expr(op api.FilterOperator, v any) *expr {
	return &expr{operator: op, target: t, value: v}
}

type expr struct {
	operator api.FilterOperator
	exprs    []Expr
	target   []string
	value    any
}

var _ Expr = (*expr)(nil)

func (e *expr) Operator() api.FilterOperator { return e.operator }
func (e *expr) Exprs() []Expr                { return e.exprs }
func (e *expr) Target() []string             { return e.target }
func (e *expr) Value() any                   { return e.value }
