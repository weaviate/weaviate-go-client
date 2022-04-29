package filters

import (
	"fmt"
	"testing"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v4/test/helpers"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/stretchr/testify/assert"
)

func TestWhereBuilder_BuildOperandsRecursively(t *testing.T) {
	nestedTwiceColor := (&WhereBuilder{}).
		WithOperator(Equal).
		WithPath([]string{"color"}).
		WithValueString("green")

	nestedTwicePrice := (&WhereBuilder{}).
		WithOperator(LessThan).
		WithPath([]string{"price"}).
		WithValueNumber(23.99)

	nestedOnceSize := (&WhereBuilder{}).
		WithOperator(Equal).
		WithPath([]string{"size"}).
		WithValueString("large").
		WithOperands([]*WhereBuilder{nestedTwiceColor, nestedTwicePrice})

	nestedOnceCountryOfOrigin := (&WhereBuilder{}).
		WithOperator(Equal).
		WithPath([]string{"countryOfOrigin"}).
		WithValueString("Taiwan")

	where := (&WhereBuilder{}).
		WithOperator(Equal).
		WithPath([]string{"id"}).
		WithValueString("123").
		WithOperands([]*WhereBuilder{nestedOnceSize, nestedOnceCountryOfOrigin})

	res := where.Build()

	expected := &models.WhereFilter{
		Operator:    "Equal",
		Path:        []string{"id"},
		ValueString: helpers.StringPointer("123"),
		Operands: []*models.WhereFilter{
			{
				Operator:    "Equal",
				Path:        []string{"size"},
				ValueString: helpers.StringPointer("large"),
				Operands: []*models.WhereFilter{
					{
						Operator:    "Equal",
						Path:        []string{"color"},
						ValueString: helpers.StringPointer("green"),
					},
					{
						Operator:    "LessThan",
						Path:        []string{"price"},
						ValueNumber: helpers.Float64Pointer(23.99),
					},
				},
			},
			{
				Operator:    "Equal",
				Path:        []string{"countryOfOrigin"},
				ValueString: helpers.StringPointer("Taiwan"),
			},
		},
	}

	assert.Equal(t, expected, res)
}

func TestWhereBuilder_String(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		builder *WhereBuilder
		want    string
	}{
		{
			name:    "with: path operator.And text",
			builder: Where().WithPath([]string{"id"}).WithOperator(And).WithValueText("txt"),
			want:    "where:{operator: And path: [\"id\"] valueText: \"txt\"}",
		},
		{
			name:    "with: path operator.Equal string",
			builder: Where().WithPath([]string{"id"}).WithOperator(Equal).WithValueString("txt"),
			want:    "where:{operator: Equal path: [\"id\"] valueString: \"txt\"}",
		},
		{
			name:    "with: path operator.GreaterThan int",
			builder: Where().WithPath([]string{"id"}).WithOperator(GreaterThan).WithValueInt(11),
			want:    "where:{operator: GreaterThan path: [\"id\"] valueInt: 11}",
		},
		{
			name:    "with: path operator.Or bool",
			builder: Where().WithPath([]string{"id"}).WithOperator(Or).WithValueBoolean(true),
			want:    "where:{operator: Or path: [\"id\"] valueBoolean: true}",
		},
		{
			name:    "with: path operator.GreaterThanEqual number",
			builder: Where().WithPath([]string{"id"}).WithOperator(GreaterThanEqual).WithValueNumber(22.1),
			want:    "where:{operator: GreaterThanEqual path: [\"id\"] valueNumber: 22.1}",
		},
		{
			name: "with: path operator.WithinGeoRange geo",
			builder: Where().WithPath([]string{"id"}).WithOperator(WithinGeoRange).
				WithValueGeoRange(&GeoCoordinatesParameter{Latitude: 50.51, Longitude: 0.11, MaxDistance: 3000}),
			want: "where:{operator: WithinGeoRange path: [\"id\"] valueGeoRange: {geoCoordinates:{latitude:50.51,longitude:0.11},distance:{max:3000}}}",
		},
		{
			name: "with: path operator.Like date",
			builder: Where().WithPath([]string{"id"}).WithOperator(Like).
				WithValueDate(now),
			want: fmt.Sprintf("where:{operator: Like path: [\"id\"] valueDate: %s}", now.Format(time.RFC3339Nano)),
		},
		{
			name: "with: operands",
			builder: Where().WithOperator(And).
				WithOperands([]*WhereBuilder{
					Where().WithPath([]string{"wordCount"}).WithOperator(LessThanEqual).WithValueInt(10),
					Where().WithPath([]string{"word"}).WithOperator(LessThan).WithValueString("word"),
				}),
			want: "where:{operator: And operands:[{operator: LessThanEqual path: [\"wordCount\"] valueInt: 10},{operator: LessThan path: [\"word\"] valueString: \"word\"}]}",
		},
		{
			name: "with: multiple path operator.Not date",
			builder: Where().WithPath([]string{"p1", "p2", "p3"}).WithOperator(Not).
				WithValueDate(now),
			want: fmt.Sprintf("where:{operator: Not path: [\"p1\",\"p2\",\"p3\"] valueDate: %s}", now.Format(time.RFC3339Nano)),
		},
		{
			name: "with: operands with multiple path",
			builder: Where().WithOperator(And).
				WithOperands([]*WhereBuilder{
					Where().WithPath([]string{"wordCount"}).WithOperator(LessThanEqual).WithValueInt(10),
					Where().WithPath([]string{"w1", "w2", "w3"}).WithOperator(LessThan).WithValueString("word"),
				}),
			want: "where:{operator: And operands:[{operator: LessThanEqual path: [\"wordCount\"] valueInt: 10},{operator: LessThan path: [\"w1\",\"w2\",\"w3\"] valueString: \"word\"}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.builder.String(); got != tt.want {
				t.Errorf("WhereArgumentBuilder.build() = %v, want %v", got, tt.want)
			}
		})
	}
}
