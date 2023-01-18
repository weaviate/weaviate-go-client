package filters

import (
	"fmt"
	"strings"
	"time"

	"github.com/weaviate/weaviate/entities/models"
)

func Where() *WhereBuilder {
	return &WhereBuilder{}
}

type WhereBuilder struct {
	operands         []*WhereBuilder
	operator         WhereOperator
	path             []string
	withValueInt     bool
	valueInt         int64
	withValueNumber  bool
	valueNumber      float64
	withValueBoolean bool
	valueBoolean     bool
	valueString      string
	valueText        string
	withValueDate    bool
	valueDate        time.Time
	valueGeoRange    *GeoCoordinatesParameter
}

// WithOperator the operator to be used
func (b *WhereBuilder) WithOperator(operator WhereOperator) *WhereBuilder {
	b.operator = operator
	return b
}

// WithPath the list of properties that should be looked for
func (b *WhereBuilder) WithPath(path []string) *WhereBuilder {
	b.path = path
	return b
}

func (b *WhereBuilder) WithOperands(operands []*WhereBuilder) *WhereBuilder {
	b.operands = operands
	return b
}

// WithValueInt the int value in where filter
func (b *WhereBuilder) WithValueInt(valueInt int64) *WhereBuilder {
	b.withValueInt = true
	b.valueInt = valueInt
	return b
}

// WithValueNumber the number value in where filter
func (b *WhereBuilder) WithValueNumber(valueNumber float64) *WhereBuilder {
	b.withValueNumber = true
	b.valueNumber = valueNumber
	return b
}

// WithValueBoolean the boolean value in where filter
func (b *WhereBuilder) WithValueBoolean(valueBoolean bool) *WhereBuilder {
	b.withValueBoolean = true
	b.valueBoolean = valueBoolean
	return b
}

// WithValueString the string value in where filter
func (b *WhereBuilder) WithValueString(valueString string) *WhereBuilder {
	b.valueString = valueString
	return b
}

// WithValueText the string value in where filter
func (b *WhereBuilder) WithValueText(valueText string) *WhereBuilder {
	b.valueText = valueText
	return b
}

// WithValueDate the date value in where filter
func (b *WhereBuilder) WithValueDate(valueDate time.Time) *WhereBuilder {
	b.withValueDate = true
	b.valueDate = valueDate
	return b
}

// WithValueGeoRange the string value in where filter
func (b *WhereBuilder) WithValueGeoRange(valueGeoRange *GeoCoordinatesParameter) *WhereBuilder {
	b.valueGeoRange = valueGeoRange
	return b
}

// Build creates a *models.WhereFilter from a *WhereBuilder
func (b *WhereBuilder) Build() *models.WhereFilter {
	whereFilter := &models.WhereFilter{
		Operator: string(b.operator),
		Path:     b.path,
	}

	if b.withValueInt {
		whereFilter.ValueInt = &b.valueInt
	}
	if b.withValueNumber {
		whereFilter.ValueNumber = &b.valueNumber
	}
	if b.withValueBoolean {
		whereFilter.ValueBoolean = &b.valueBoolean
	}
	if len(b.valueString) > 0 {
		whereFilter.ValueString = &b.valueString
	}
	if len(b.valueText) > 0 {
		whereFilter.ValueText = &b.valueText
	}
	if b.withValueDate {
		formattedDate := b.valueDate.Format(time.RFC3339Nano)
		whereFilter.ValueDate = &formattedDate
	}
	if b.valueGeoRange != nil {
		whereFilter.ValueGeoRange = buildWhereFilterGeoRange(b.valueGeoRange)
	}

	// recursively build operands
	for _, op := range b.operands {
		whereFilter.Operands = append(whereFilter.Operands, op.Build())
	}

	return whereFilter
}

// String formats the where builder as a string for GQL queries
func (b *WhereBuilder) String() string {
	return fmt.Sprintf("where:{%s}", b.string())
}

func (b *WhereBuilder) string() string {
	clause := []string{}
	if len(b.operator) > 0 {
		clause = append(clause, fmt.Sprintf("operator: %s", b.operator))
	}
	if len(b.path) > 0 {
		path := make([]string, len(b.path))
		for i := range b.path {
			path[i] = fmt.Sprintf("\"%s\"", b.path[i])
		}
		clause = append(clause, fmt.Sprintf("path: [%v]", strings.Join(path, ",")))
	}
	if b.withValueInt {
		clause = append(clause, fmt.Sprintf("valueInt: %v", b.valueInt))
	}
	if b.withValueNumber {
		clause = append(clause, fmt.Sprintf("valueNumber: %v", b.valueNumber))
	}
	if b.withValueBoolean {
		clause = append(clause, fmt.Sprintf("valueBoolean: %v", b.valueBoolean))
	}
	if len(b.valueString) > 0 {
		clause = append(clause, fmt.Sprintf("valueString: \"%s\"", b.valueString))
	}
	if len(b.valueText) > 0 {
		clause = append(clause, fmt.Sprintf("valueText: \"%s\"", b.valueText))
	}
	if b.withValueDate {
		clause = append(clause, fmt.Sprintf("valueDate: \"%s\"", b.valueDate.Format(time.RFC3339Nano)))
	}
	if b.valueGeoRange != nil {
		clause = append(clause, fmt.Sprintf("valueGeoRange: {geoCoordinates:{latitude:%v,longitude:%v},distance:{max:%v}}",
			b.valueGeoRange.Latitude, b.valueGeoRange.Longitude, b.valueGeoRange.MaxDistance))
	}
	if len(b.operands) > 0 {
		operands := make([]string, len(b.operands))
		for i := range b.operands {
			operands[i] = fmt.Sprintf("{%s}", b.operands[i].string())
		}
		clause = append(clause, fmt.Sprintf("operands:[%s]", strings.Join(operands, ",")))
	}
	return strings.Join(clause, " ")
}

func buildWhereFilterGeoRange(in *GeoCoordinatesParameter) *models.WhereFilterGeoRange {
	out := &models.WhereFilterGeoRange{
		Distance: &models.WhereFilterGeoRangeDistance{
			Max: float64(in.MaxDistance),
		},
		GeoCoordinates: &models.GeoCoordinates{
			Latitude:  &in.Latitude,
			Longitude: &in.Longitude,
		},
	}

	return out
}
