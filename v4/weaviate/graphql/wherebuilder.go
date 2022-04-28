package graphql

import (
	"fmt"
	"strings"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/filters"
)

type WhereArgumentBuilder struct {
	operands         []*filters.WhereBuilder
	operator         filters.WhereOperator
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
	valueGeoRange    *filters.GeoCoordinatesParameter
}

// WithOperator the operator to be used
func (b *WhereArgumentBuilder) WithOperator(operator filters.WhereOperator) *WhereArgumentBuilder {
	b.operator = operator
	return b
}

// WithPath the list of properties that should be looked for
func (b *WhereArgumentBuilder) WithPath(path []string) *WhereArgumentBuilder {
	b.path = path
	return b
}

// WithValueInt the int value in where filter
func (b *WhereArgumentBuilder) WithValueInt(valueInt int64) *WhereArgumentBuilder {
	b.withValueInt = true
	b.valueInt = valueInt
	return b
}

// WithValueNumber the number value in where filter
func (b *WhereArgumentBuilder) WithValueNumber(valueNumber float64) *WhereArgumentBuilder {
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
func (b *WhereArgumentBuilder) WithValueDate(valueDate time.Time) *WhereArgumentBuilder {
	b.withValueDate = true
	b.valueDate = valueDate
	return b
}

// WithValueGeoRange the string value in where filter
func (b *WhereArgumentBuilder) WithValueGeoRange(valueGeoRange *filters.GeoCoordinatesParameter) *WhereArgumentBuilder {
	b.valueGeoRange = valueGeoRange
	return b
}

// WithOperands the operands to be used
func (b *WhereArgumentBuilder) WithOperands(operands []*filters.WhereBuilder) *WhereArgumentBuilder {
	b.operands = operands
	return b
}

// Build build the given clause
func (b *WhereArgumentBuilder) build() string {
	clause := []string{}
	clause = append(clause, b.buildWhereFilter())
	if len(b.operands) > 0 {
		operands := make([]string, len(b.operands))
		for i := range b.operands {
			operands[i] = fmt.Sprintf("{%s}", b.operands[i].String())
		}
		clause = append(clause, fmt.Sprintf("operands:[%s]", strings.Join(operands, ",")))
	}
	return fmt.Sprintf("where:{%s}", strings.Join(clause, " "))
}

func (b *WhereArgumentBuilder) buildWhereFilter() string {
	whereFilter := &filters.WhereBuilder{}
	whereFilter.WithOperator(b.operator)
	whereFilter.WithPath(b.path)

	if b.withValueInt {
		whereFilter.WithValueInt(b.valueInt)
	}
	if b.withValueNumber {
		whereFilter.WithValueNumber(b.valueNumber)
	}
	if b.withValueBoolean {
		whereFilter.WithValueBoolean(b.valueBoolean)
	}
	if len(b.valueString) > 0 {
		whereFilter.WithValueString(b.valueString)
	}
	if len(b.valueText) > 0 {
		whereFilter.WithValueText(b.valueText)
	}
	if b.withValueDate {
		whereFilter.WithValueDate(b.valueDate)
	}
	if b.valueGeoRange != nil {
		whereFilter.WithValueGeoRange(b.valueGeoRange)
	}

	return whereFilter.String()
}
