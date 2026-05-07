package filter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

func TestFilter(t *testing.T) {
	for _, tt := range []struct {
		name string
		expr filter.Expr
		want api.Filter
	}{
		{
			name: "eq",
			expr: filter.Property[int]("size").Equal(3),
			want: api.Filter{
				Operator: api.FilterOperatorEqual,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "lt",
			expr: filter.Property[int]("size").LessThan(3),
			want: api.Filter{
				Operator: api.FilterOperatorLessThan,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "lte",
			expr: filter.Property[int]("size").LessThanEqual(3),
			want: api.Filter{
				Operator: api.FilterOperatorLessThanEqual,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "gt",
			expr: filter.Property[int]("size").GreaterThan(3),
			want: api.Filter{
				Operator: api.FilterOperatorGreaterThan,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "gte",
			expr: filter.Property[int]("size").GreaterThanEqual(3),
			want: api.Filter{
				Operator: api.FilterOperatorGreaterThanEqual,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "like",
			expr: filter.Property[string]("model").Like("[0-9]+Roadster"),
			want: api.Filter{
				Operator: api.FilterOperatorLike,
				Target:   []string{"model"},
				Value:    "[0-9]+Roadster",
			},
		},
		{
			name: "null",
			expr: filter.Property[string]("discount").Null(),
			want: api.Filter{
				Operator: api.FilterOperatorIsNull,
				Target:   []string{"discount"},
			},
		},
		{
			name: "contains all",
			expr: filter.Property[int]("gears").ContainsAll(1, 2, 3),
			want: api.Filter{
				Operator: api.FilterOperatorContainsAll,
				Target:   []string{"gears"},
				Value:    []int{1, 2, 3},
			},
		},
		{
			name: "contains any",
			expr: filter.Property[int]("gears").ContainsAny(1, 2, 3),
			want: api.Filter{
				Operator: api.FilterOperatorContainsAny,
				Target:   []string{"gears"},
				Value:    []int{1, 2, 3},
			},
		},
		{
			name: "contains none",
			expr: filter.Property[int]("gears").ContainsNone(1, 2, 3),
			want: api.Filter{
				Operator: api.FilterOperatorContainsNone,
				Target:   []string{"gears"},
				Value:    []int{1, 2, 3},
			},
		},
		{
			name: "len(property)",
			expr: filter.Len("model").Equal(4),
			want: api.Filter{
				Operator: api.FilterOperatorEqual,
				Target:   []string{"len(model)"},
				Value:    4,
			},
		},
		{
			name: "reference count",
			expr: filter.ReferenceCount("soldIn").LessThan(10),
			want: api.Filter{
				Operator: api.FilterOperatorLessThan,
				Target:   []string{"count(soldIn)"},
				Value:    10,
			},
		},
		{
			name: "and",
			expr: filter.And{
				filter.Property[int]("length").Equal(2),
				filter.Property[int]("width").Equal(3),
			},
			want: api.Filter{
				Operator: api.FilterOperatorAnd,
				Exprs: []api.Filter{
					{Operator: api.FilterOperatorEqual, Target: []string{"length"}, Value: 2},
					{Operator: api.FilterOperatorEqual, Target: []string{"width"}, Value: 3},
				},
			},
		},
		{
			name: "or",
			expr: filter.Or{
				filter.Property[int]("length").Equal(2),
				filter.Property[int]("width").Equal(3),
			},
			want: api.Filter{
				Operator: api.FilterOperatorOr,
				Exprs: []api.Filter{
					{Operator: api.FilterOperatorEqual, Target: []string{"length"}, Value: 2},
					{Operator: api.FilterOperatorEqual, Target: []string{"width"}, Value: 3},
				},
			},
		},
		{
			name: "not",
			expr: filter.Not(filter.Property[int]("length").Equal(2)),
			want: api.Filter{
				Operator: api.FilterOperatorNot,
				Exprs: []api.Filter{
					{Operator: api.FilterOperatorEqual, Target: []string{"length"}, Value: 2},
				},
			},
		},
		{
			name: "reference",
			expr: filter.Property[string]("ownedBy", "name").Like(".*_doe"),
			want: api.Filter{
				Operator: api.FilterOperatorLike,
				Target:   []string{"ownedBy", "name"},
				Value:    ".*_doe",
			},
		},
		{
			name: "len(property) in reference",
			expr: filter.Len("ownedBy", "name").GreaterThan(12),
			want: api.Filter{
				Operator: api.FilterOperatorGreaterThan,
				Target:   []string{"ownedBy", "len(name)"},
				Value:    12,
			},
		},
		{
			name: "reference count in reference",
			expr: filter.ReferenceCount("ownedBy", "hasFriends").LessThan(4),
			want: api.Filter{
				Operator: api.FilterOperatorLessThan,
				Target:   []string{"ownedBy", "count(hasFriends)"},
				Value:    4,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkExpr(t, tt.expr, tt.want)
		})
	}
}

// checkExpr asserts that [filter.Expr] returns the expected operator,
// target, test value, and sub-expressions.
func checkExpr(t *testing.T, e filter.Expr, want api.Filter) {
	assert.Equal(t, want.Operator, e.Operator(), "want %s, got %s", want.Operator, e.Operator())
	assert.Equal(t, want.Target, e.Target(), "target")
	assert.EqualValues(t, want.Value, e.Value(), "value")

	exprs := e.Exprs()
	assert.Len(t, exprs, len(want.Exprs), "sub-expressions")
	for i := range exprs {
		checkExpr(t, exprs[i], want.Exprs[i])
	}
}
