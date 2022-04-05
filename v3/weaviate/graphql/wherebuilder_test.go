package graphql

import (
	"testing"
)

func TestWhereArgumentBuilder_build(t *testing.T) {
	type fields struct {
		operands      []WhereFilterBuilder
		operator      WhereOperator
		path          []string
		valueInt      int
		valueNumber   float32
		valueBoolean  bool
		valueString   string
		valueText     string
		valueDate     string
		valueGeoRange *GeoCoordinatesParameter
	}
	tests := []struct {
		name    string
		builder *WhereArgumentBuilder
		want    string
	}{
		{
			name:    "with: path operator.And text",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(And).WithValueText("txt"),
			want:    "where:{operator: And path: [\"id\"] valueText: \"txt\"}",
		},
		{
			name:    "with: path operator.Equal string",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(Equal).WithValueString("txt"),
			want:    "where:{operator: Equal path: [\"id\"] valueString: \"txt\"}",
		},
		{
			name:    "with: path operator.GreaterThan int",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(GreaterThan).WithValueInt(11),
			want:    "where:{operator: GreaterThan path: [\"id\"] valueInt: 11}",
		},
		{
			name:    "with: path operator.Or bool",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(Or).WithValueBoolean(true),
			want:    "where:{operator: Or path: [\"id\"] valueBoolean: true}",
		},
		{
			name:    "with: path operator.GreaterThanEqual number",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(GreaterThanEqual).WithValueNumber(22.1),
			want:    "where:{operator: GreaterThanEqual path: [\"id\"] valueNumber: 22.1}",
		},
		{
			name: "with: path operator.WithinGeoRange geo",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(WithinGeoRange).
				WithValueGeoRange(&GeoCoordinatesParameter{Latitude: 50.51, Longitude: 0.11, MaxDistance: 3000}),
			want: "where:{operator: WithinGeoRange path: [\"id\"] valueGeoRange: {geoCoordinates:{latitude:50.51,longitude:0.11},distance:{max:3000}}}",
		},
		{
			name: "with: path operator.Like date",
			builder: newWhereArgBuilder().WithPath([]string{"id"}).WithOperator(Like).
				WithValueDate("2020-01-01T00:00:00-07:00"),
			want: "where:{operator: Like path: [\"id\"] valueDate: 2020-01-01T00:00:00-07:00}",
		},
		{
			name: "with: operands",
			builder: newWhereArgBuilder().WithOperator(And).
				WithOperands([]*WhereFilterBuilder{
					newWhereFilter().WithPath([]string{"wordCount"}).WithOperator(LessThanEqual).WithValueInt(10),
					newWhereFilter().WithPath([]string{"word"}).WithOperator(LessThan).WithValueString("word"),
				}),
			want: "where:{operator: And operands:[{operator: LessThanEqual path: [\"wordCount\"] valueInt: 10},{operator: LessThan path: [\"word\"] valueString: \"word\"}]}",
		},
		{
			name: "with: multiple path operator.Not date",
			builder: newWhereArgBuilder().WithPath([]string{"p1", "p2", "p3"}).WithOperator(Not).
				WithValueDate("2021-01-01T00:00:10-07:00"),
			want: "where:{operator: Not path: [\"p1\",\"p2\",\"p3\"] valueDate: 2021-01-01T00:00:10-07:00}",
		},
		{
			name: "with: operands with multiple path",
			builder: newWhereArgBuilder().WithOperator(And).
				WithOperands([]*WhereFilterBuilder{
					newWhereFilter().WithPath([]string{"wordCount"}).WithOperator(LessThanEqual).WithValueInt(10),
					newWhereFilter().WithPath([]string{"w1", "w2", "w3"}).WithOperator(LessThan).WithValueString("word"),
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

func newWhereFilter() *WhereFilterBuilder {
	return &WhereFilterBuilder{}
}
