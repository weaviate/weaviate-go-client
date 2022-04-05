package graphql

import (
	"fmt"
	"strings"
)

type WhereArgumentBuilder struct {
	operands         []*WhereFilterBuilder
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
func (b *WhereArgumentBuilder) WithOperator(operator WhereOperator) *WhereArgumentBuilder {
	b.operator = operator
	return b
}

// WithPath the list of properties that should be looked for
func (b *WhereArgumentBuilder) WithPath(path []string) *WhereArgumentBuilder {
	b.path = path
	return b
}

// WithValueInt the int value in where filter
func (b *WhereArgumentBuilder) WithValueInt(valueInt int) *WhereArgumentBuilder {
	b.withValueInt = true
	b.valueInt = valueInt
	return b
}

// WithValueNumber the number value in where filter
func (b *WhereArgumentBuilder) WithValueNumber(valueNumber float32) *WhereArgumentBuilder {
	b.withValueNumber = true
	b.valueNumber = valueNumber
	return b
}

// WithValueBoolean the boolean value in where filter
func (b *WhereArgumentBuilder) WithValueBoolean(valueBoolean bool) *WhereArgumentBuilder {
	b.withValueBoolean = true
	b.valueBoolean = valueBoolean
	return b
}

// WithValueString the string value in where filter
func (b *WhereArgumentBuilder) WithValueString(valueString string) *WhereArgumentBuilder {
	b.valueString = valueString
	return b
}

// WithValueText the string value in where filter
func (b *WhereArgumentBuilder) WithValueText(valueText string) *WhereArgumentBuilder {
	b.valueText = valueText
	return b
}

// WithValueDate the date value in where filter
func (b *WhereArgumentBuilder) WithValueDate(valueDate string) *WhereArgumentBuilder {
	b.withValueDate = true
	b.valueDate = valueDate
	return b
}

// WithValueGeoRange the string value in where filter
func (b *WhereArgumentBuilder) WithValueGeoRange(valueGeoRange *GeoCoordinatesParameter) *WhereArgumentBuilder {
	b.valueGeoRange = valueGeoRange
	return b
}

// WithOperands the operands to be used
func (b *WhereArgumentBuilder) WithOperands(operands []*WhereFilterBuilder) *WhereArgumentBuilder {
	b.operands = operands
	return b
}

// Build build the given clause
func (b *WhereArgumentBuilder) build() string {
	clause := []string{}
	whereFilter := &WhereFilterBuilder{
		operator:         b.operator,
		path:             b.path,
		withValueInt:     b.withValueInt,
		valueInt:         b.valueInt,
		withValueNumber:  b.withValueNumber,
		valueNumber:      b.valueNumber,
		withValueBoolean: b.withValueBoolean,
		valueBoolean:     b.valueBoolean,
		valueString:      b.valueString,
		valueText:        b.valueText,
		withValueDate:    b.withValueDate,
		valueDate:        b.valueDate,
		valueGeoRange:    b.valueGeoRange,
	}
	clause = append(clause, whereFilter.build())
	if len(b.operands) > 0 {
		operands := make([]string, len(b.operands))
		for i := range b.operands {
			operands[i] = fmt.Sprintf("{%s}", b.operands[i].build())
		}
		clause = append(clause, fmt.Sprintf("operands:[%s]", strings.Join(operands, ",")))
	}
	return fmt.Sprintf("where:{%s}", strings.Join(clause, " "))
}
