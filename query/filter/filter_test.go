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
		want *api.FilterExpr
	}{
		{
			name: "eq",
			expr: &filter.Cond{
				Target:   "size",
				Operator: filter.Equal,
				Value:    3,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorEqual,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "lt",
			expr: &filter.Cond{
				Target:   "size",
				Operator: filter.LessThan,
				Value:    3,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorLessThan,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "lte",
			expr: &filter.Cond{
				Target:   "size",
				Operator: filter.LessThanEqual,
				Value:    3,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorLessThanEqual,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "gt",
			expr: &filter.Cond{
				Target:   "size",
				Operator: filter.GreaterThan,
				Value:    3,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorGreaterThan,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "gte",
			expr: &filter.Cond{
				Target:   "size",
				Operator: filter.GreaterThanEqual,
				Value:    3,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorGreaterThanEqual,
				Target:   []string{"size"},
				Value:    3,
			},
		},
		{
			name: "like",
			expr: &filter.Cond{
				Target:   "model",
				Operator: filter.Like,
				Value:    "[0-9]+Roadster",
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorLike,
				Target:   []string{"model"},
				Value:    "[0-9]+Roadster",
			},
		},
		{
			name: "null",
			expr: &filter.Cond{
				Target:   "discount",
				Operator: filter.IsNull,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorIsNull,
				Target:   []string{"discount"},
			},
		},
		{
			name: "contains all",
			expr: &filter.Cond{
				Target:   "gears",
				Operator: filter.ContainsAll,
				Value:    []int{1, 2, 3},
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorContainsAll,
				Target:   []string{"gears"},
				Value:    []int{1, 2, 3},
			},
		},
		{
			name: "contains any",
			expr: &filter.Cond{
				Target:   "gears",
				Operator: filter.ContainsAny,
				Value:    []int{1, 2, 3},
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorContainsAny,
				Target:   []string{"gears"},
				Value:    []int{1, 2, 3},
			},
		},
		{
			name: "contains none",
			expr: &filter.Cond{
				Target:   "gears",
				Operator: filter.ContainsNone,
				Value:    []int{1, 2, 3},
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorContainsNone,
				Target:   []string{"gears"},
				Value:    []int{1, 2, 3},
			},
		},
		{
			name: "len(property)",
			expr: &filter.Cond{
				Target:   filter.Len("model"),
				Operator: filter.Equal,
				Value:    4,
			},
			want: &api.FilterExpr{
				Target:   []string{"len(model)"},
				Operator: api.FilterOperatorEqual,
				Value:    4,
			},
		},
		{
			name: "reference count",
			expr: &filter.Cond{
				Target:   filter.Reference("soldIn").Count(),
				Operator: filter.LessThan,
				Value:    10,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorLessThan,
				Target:   []string{"count(soldIn)"},
				Value:    10,
			},
		},
		{
			name: "and",
			expr: filter.And{
				&filter.Cond{
					Target:   "length",
					Operator: filter.Equal,
					Value:    2,
				},
				&filter.Cond{
					Target:   "width",
					Operator: filter.Equal,
					Value:    3,
				},
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorAnd,
				Exprs: []api.FilterExpr{
					{Operator: api.FilterOperatorEqual, Target: []string{"length"}, Value: 2},
					{Operator: api.FilterOperatorEqual, Target: []string{"width"}, Value: 3},
				},
			},
		},
		{
			name: "or",
			expr: filter.Or{
				&filter.Cond{
					Target:   "length",
					Operator: filter.Equal,
					Value:    2,
				},
				&filter.Cond{
					Target:   "width",
					Operator: filter.Equal,
					Value:    3,
				},
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorOr,
				Exprs: []api.FilterExpr{
					{Operator: api.FilterOperatorEqual, Target: []string{"length"}, Value: 2},
					{Operator: api.FilterOperatorEqual, Target: []string{"width"}, Value: 3},
				},
			},
		},
		{
			name: "not",
			expr: filter.Not{
				&filter.Cond{
					Target:   "length",
					Operator: filter.Equal,
					Value:    2,
				},
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorNot,
				Exprs: []api.FilterExpr{
					{Operator: api.FilterOperatorEqual, Target: []string{"length"}, Value: 2},
				},
			},
		},
		{
			name: "reference",
			expr: &filter.Cond{
				Target:   filter.Reference("ownedBy").Property("name"),
				Operator: filter.Like,
				Value:    ".*_doe",
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorLike,
				Target:   []string{"ownedBy", "name"},
				Value:    ".*_doe",
			},
		},
		{
			name: "len(property) in reference",
			expr: &filter.Cond{
				Target:   filter.Len(filter.Reference("ownedBy").Property("name")),
				Operator: filter.GreaterThan,
				Value:    12,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorGreaterThan,
				Target:   []string{"ownedBy", "len(name)"},
				Value:    12,
			},
		},
		{
			name: "reference count in reference",
			expr: &filter.Cond{
				Target:   filter.Reference("ownedBy").Reference("hasFriends").Count(),
				Operator: filter.LessThan,
				Value:    4,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorLessThan,
				Target:   []string{"ownedBy", "count(hasFriends)"},
				Value:    4,
			},
		},
		{
			name: "multi-target reference",
			expr: &filter.Cond{
				Operator: filter.GreaterThanEqual,
				Target:   filter.Reference("hasAwards").Collection("GrammyAwards").Property("year"),
				Value:    1974,
			},
			want: &api.FilterExpr{
				Operator: api.FilterOperatorGreaterThanEqual,
				Target:   []string{"hasAwards:GrammyAwards", "year"},
				Value:    1974,
			},
		},
		{
			name: "empty AND group",
			expr: filter.And{},
			want: nil,
		},
		{
			name: "empty OR group",
			expr: filter.Or{},
			want: nil,
		},
		{
			name: "nil filter in AND group",
			expr: filter.And{nil},
			want: nil,
		},
		{
			name: "nil filter in OR group",
			expr: filter.Or{nil},
			want: nil,
		},
		{
			name: "negating a nil filter",
			expr: filter.Not{nil},
			want: nil,
		},
		{
			name: "negating an empty AND group",
			expr: filter.Not{filter.And{}},
			want: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkExpr(t, tt.expr.Expr(), tt.want)
		})
	}
}

// checkExpr asserts that [filter.Expr] returns the expected operator,
// target, test value, and sub-expressions.
func checkExpr(t *testing.T, got, want *api.FilterExpr) {
	assert.Equal(t, want, got)
	if want == nil {
		return
	}
	assert.Len(t, got.Exprs, len(want.Exprs), "sub-expressions")
	for i := range got.Exprs {
		checkExpr(t, &got.Exprs[i], &want.Exprs[i])
	}
}
