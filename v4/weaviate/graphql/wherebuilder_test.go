package graphql

import (
	"fmt"
	"testing"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
)

func TestWhereArgumentBuilder_build(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		builder *WhereArgumentBuilder
		want    string
	}{
		{
			name:    "with: path operator.And text",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(filters.And).WithValueText("txt"),
			want:    "where:{operator: And path: [\"id\"] valueText: \"txt\"}",
		},
		{
			name:    "with: path operator.Equal string",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(filters.Equal).WithValueString("txt"),
			want:    "where:{operator: Equal path: [\"id\"] valueString: \"txt\"}",
		},
		{
			name:    "with: path operator.GreaterThan int",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(filters.GreaterThan).WithValueInt(11),
			want:    "where:{operator: GreaterThan path: [\"id\"] valueInt: 11}",
		},
		{
			name:    "with: path operator.Or bool",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(filters.Or).WithValueBoolean(true),
			want:    "where:{operator: Or path: [\"id\"] valueBoolean: true}",
		},
		{
			name:    "with: path operator.GreaterThanEqual number",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(filters.GreaterThanEqual).WithValueNumber(22.1),
			want:    "where:{operator: GreaterThanEqual path: [\"id\"] valueNumber: 22.1}",
		},
		{
			name: "with: path operator.WithinGeoRange geo",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(filters.WithinGeoRange).
				WithValueGeoRange(&filters.GeoCoordinatesParameter{Latitude: 50.51, Longitude: 0.11, MaxDistance: 3000}),
			want: "where:{operator: WithinGeoRange path: [\"id\"] valueGeoRange: {geoCoordinates:{latitude:50.51,longitude:0.11},distance:{max:3000}}}",
		},
		{
			name: "with: path operator.Like date",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(filters.Like).
				WithValueDate(now),
			want: fmt.Sprintf("where:{operator: Like path: [\"id\"] valueDate: %s}", now.Format(time.RFC3339Nano)),
		},
		{
			name: "with: operands",
			builder: newWhereArgBuilder().WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					newWhereFilter().WithPath([]string{"wordCount"}).WithOperator(filters.LessThanEqual).WithValueInt(10),
					newWhereFilter().WithPath([]string{"word"}).WithOperator(filters.LessThan).WithValueString("word"),
				}),
			want: "where:{operator: And operands:[{operator: LessThanEqual path: [\"wordCount\"] valueInt: 10},{operator: LessThan path: [\"word\"] valueString: \"word\"}]}",
		},
		{
			name: "with: multiple path operator.Not date",
			builder: newWhereArgBuilder().WithPath([]string{"p1", "p2", "p3"}).WithOperator(filters.Not).
				WithValueDate(now),
			want: fmt.Sprintf("where:{operator: Not path: [\"p1\",\"p2\",\"p3\"] valueDate: %s}", now.Format(time.RFC3339Nano)),
		},
		{
			name: "with: operands with multiple path",
			builder: newWhereArgBuilder().WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					newWhereFilter().WithPath([]string{"wordCount"}).WithOperator(filters.LessThanEqual).WithValueInt(10),
					newWhereFilter().WithPath([]string{"w1", "w2", "w3"}).WithOperator(filters.LessThan).WithValueString("word"),
				}),
			want: "where:{operator: And operands:[{operator: LessThanEqual path: [\"wordCount\"] valueInt: 10},{operator: LessThan path: [\"w1\",\"w2\",\"w3\"] valueString: \"word\"}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.build(); got != tt.want {
				t.Errorf("WhereArgumentBuilder.build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newWhereArgBuilder() *WhereArgumentBuilder {
	return &WhereArgumentBuilder{}
}

func newWhereFilter() *filters.WhereBuilder {
	return &filters.WhereBuilder{}
}
