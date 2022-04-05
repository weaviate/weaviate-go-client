package graphql

import (
	"fmt"
	"strings"
)

type WhereFilterBuilder struct {
	operator         WhereOperator
	path             []string
	withValueInt     bool
	valueInt         int
	withValueNumber  bool
	valueNumber      float32
	withValueBoolean bool
	valueBoolean     bool
	valueString      string
	valueText        string
	withValueDate    bool
	valueDate        string
	valueGeoRange    *GeoCoordinatesParameter
}

// WithOperator the operator to be used
func (b *WhereFilterBuilder) WithOperator(operator WhereOperator) *WhereFilterBuilder {
	b.operator = operator
	return b
}

// WithPath the list of properties that should be looked for
func (b *WhereFilterBuilder) WithPath(path []string) *WhereFilterBuilder {
	b.path = path
	return b
}

// WithValueInt the int value in where filter
func (b *WhereFilterBuilder) WithValueInt(valueInt int) *WhereFilterBuilder {
	b.withValueInt = true
	b.valueInt = valueInt
	return b
}

// WithValueNumber the number value in where filter
func (b *WhereFilterBuilder) WithValueNumber(valueNumber float32) *WhereFilterBuilder {
	b.withValueNumber = true
	b.valueNumber = valueNumber
	return b
}

// WithValueBoolean the boolean value in where filter
func (b *WhereFilterBuilder) WithValueBoolean(valueBoolean bool) *WhereFilterBuilder {
	b.withValueBoolean = true
	b.valueBoolean = valueBoolean
	return b
}

// WithValueString the string value in where filter
func (b *WhereFilterBuilder) WithValueString(valueString string) *WhereFilterBuilder {
	b.valueString = valueString
	return b
}

// WithValueText the string value in where filter
func (b *WhereFilterBuilder) WithValueText(valueText string) *WhereFilterBuilder {
	b.valueText = valueText
	return b
}

// WithValueDate the date value in where filter
func (b *WhereFilterBuilder) WithValueDate(valueDate string) *WhereFilterBuilder {
	b.withValueDate = true
	b.valueDate = valueDate
	return b
}

// WithValueGeoRange the string value in where filter
func (b *WhereFilterBuilder) WithValueGeoRange(valueGeoRange *GeoCoordinatesParameter) *WhereFilterBuilder {
	b.valueGeoRange = valueGeoRange
	return b
}

// Build build the given clause
func (b *WhereFilterBuilder) build() string {
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
		clause = append(clause, fmt.Sprintf("valueDate: %s", b.valueDate))
	}
	if b.valueGeoRange != nil {
		clause = append(clause, fmt.Sprintf("valueGeoRange: {geoCoordinates:{latitude:%v,longitude:%v},distance:{max:%v}}",
			b.valueGeoRange.Latitude, b.valueGeoRange.Longitude, b.valueGeoRange.MaxDistance))
	}
	return fmt.Sprintf("%s", strings.Join(clause, " "))
}
