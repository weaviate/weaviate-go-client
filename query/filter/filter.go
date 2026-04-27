package filter

import (
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

// NB(dyma): These need to be public interface methods, because we want the `query`
// package to be able to marshal filters independently
type Expr interface {
	Operator() api.FilterOperator // api.FilterOperator
	Exprs() []Expr                // For AND/OR

	Target() []string
	Value() any
}

func Property[T any](property string) Builder[T] { return target[T]{property} }

var (
	UUID          = Property[uuid.UUID](api.FieldUUID)
	CreatedAt     = Property[time.Time](api.FieldCreatedAt)
	LastUpdatedAt = Property[time.Time](api.FieldLastUpdatedAt)
)

type Builder[T any] interface {
	Equal(T) Expr
	NotEqual(T) Expr
	LessThan(T) Expr
	LessThanEqual(T) Expr
	GreaterThan(T) Expr
	GreaterThanEqual(T) Expr
	Like(string) Expr
	ContainsAny(...T) Expr
	ContainsAll(...T) Expr
	ContainsNone(...T) Expr
}

func Not(e Expr) Expr { return &expr{operator: api.FilterOperatorNot, exprs: []Expr{e}} }
func Null(target string) Expr {
	return &expr{operator: api.FilterOperatorIsNull, target: []string{target}}
}

type And []Expr

func (and And) Exprs() []Expr                { return and }
func (and And) Operator() api.FilterOperator { return api.FilterOperatorAnd }
func (and And) Target() []string             { return nil }
func (and And) Value() any                   { return nil }

var _ Expr = (And)(nil)

type Or []Expr

func (or Or) Exprs() []Expr                { return or }
func (or Or) Operator() api.FilterOperator { return api.FilterOperatorOr }
func (or Or) Target() []string             { return nil }
func (or Or) Value() any                   { return nil }

type target[T any] []string

func (t target[T]) Equal(v T) Expr            { return t.expr(api.FilterOperatorEqual, v) }
func (t target[T]) NotEqual(v T) Expr         { return Not(t.Equal(v)) }
func (t target[T]) LessThan(v T) Expr         { return t.expr(api.FilterOperatorLessThan, v) }
func (t target[T]) LessThanEqual(v T) Expr    { return t.expr(api.FilterOperatorLessThanEqual, v) }
func (t target[T]) GreaterThan(v T) Expr      { return t.expr(api.FilterOperatorGreaterThan, v) }
func (t target[T]) GreaterThanEqual(v T) Expr { return t.expr(api.FilterOperatorGreaterThanEqual, v) }
func (t target[T]) Like(v string) Expr        { return t.expr(api.FilterOperatorLike, v) }
func (t target[T]) ContainsAll(vs ...T) Expr  { return t.expr(api.FilterOperatorContainsAll, vs) }
func (t target[T]) ContainsAny(vs ...T) Expr  { return t.expr(api.FilterOperatorContainsAny, vs) }
func (t target[T]) ContainsNone(vs ...T) Expr { return t.expr(api.FilterOperatorContainsNone, vs) }
func (t target[T]) expr(op api.FilterOperator, v any) *expr {
	return &expr{operator: op, target: t, value: v}
}

var _ Builder[any] = (*target[any])(nil)

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
