package filters

import (
	"fmt"
	"strings"
	"time"

	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

func Where() *WhereBuilder {
	return &WhereBuilder{}
}

type WhereBuilder struct {
	operands         []*WhereBuilder
	operator         WhereOperator
	path             []string
	withValueInt     bool
	valueInt         []int64
	withValueNumber  bool
	valueNumber      []float64
	withValueBoolean bool
	valueBoolean     []bool
	valueString      []string
	valueText        []string
	withValueDate    bool
	valueDate        []time.Time
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
func (b *WhereBuilder) WithValueInt(valueInt ...int64) *WhereBuilder {
	b.withValueInt = true
	b.valueInt = valueInt
	return b
}

// WithValueNumber the number value in where filter
func (b *WhereBuilder) WithValueNumber(valueNumber ...float64) *WhereBuilder {
	b.withValueNumber = true
	b.valueNumber = valueNumber
	return b
}

// WithValueBoolean the boolean value in where filter
func (b *WhereBuilder) WithValueBoolean(valueBoolean ...bool) *WhereBuilder {
	b.withValueBoolean = true
	b.valueBoolean = valueBoolean
	return b
}

// WithValueString the string value in where filter
func (b *WhereBuilder) WithValueString(valueString ...string) *WhereBuilder {
	b.valueString = valueString
	return b
}

// WithValueText the string value in where filter
func (b *WhereBuilder) WithValueText(valueText ...string) *WhereBuilder {
	b.valueText = valueText
	return b
}

// WithValueDate the date value in where filter
func (b *WhereBuilder) WithValueDate(valueDate ...time.Time) *WhereBuilder {
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
		if len(b.valueInt) == 1 && !isContainsOperator(b.operator) {
			whereFilter.ValueInt = &b.valueInt[0]
		} else {
			whereFilter.ValueIntArray = b.valueInt
		}
	}
	if b.withValueNumber {
		if len(b.valueNumber) == 1 && !isContainsOperator(b.operator) {
			whereFilter.ValueNumber = &b.valueNumber[0]
		} else {
			whereFilter.ValueNumberArray = b.valueNumber
		}
	}
	if b.withValueBoolean {
		if len(b.valueBoolean) == 1 && !isContainsOperator(b.operator) {
			whereFilter.ValueBoolean = &b.valueBoolean[0]
		} else {
			whereFilter.ValueBooleanArray = b.valueBoolean
		}
	}
	if len(b.valueString) > 0 {
		if len(b.valueString) == 1 && !isContainsOperator(b.operator) {
			whereFilter.ValueString = &b.valueString[0]
		} else {
			whereFilter.ValueStringArray = b.valueString
		}
	}
	if len(b.valueText) > 0 {
		if len(b.valueText) == 1 && !isContainsOperator(b.operator) {
			whereFilter.ValueText = &b.valueText[0]
		} else {
			whereFilter.ValueTextArray = b.valueText
		}
	}
	if b.withValueDate {
		formattedDates := formatDates(b.valueDate)
		if len(formattedDates) == 1 && !isContainsOperator(b.operator) {
			whereFilter.ValueDate = &formattedDates[0]
		} else {
			whereFilter.ValueDateArray = formattedDates
		}
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
		clause = append(clause, fmt.Sprintf("valueInt: %s", formatValues(b.valueInt, b.operator)))
	}
	if b.withValueNumber {
		clause = append(clause, fmt.Sprintf("valueNumber: %s", formatValues(b.valueNumber, b.operator)))
	}
	if b.withValueBoolean {
		clause = append(clause, fmt.Sprintf("valueBoolean: %s", formatValues(b.valueBoolean, b.operator)))
	}
	if len(b.valueString) > 0 {
		clause = append(clause, fmt.Sprintf("valueString: %s", formatValues(b.valueString, b.operator)))
	}
	if len(b.valueText) > 0 {
		clause = append(clause, fmt.Sprintf("valueText: %s", formatValues(b.valueText, b.operator)))
	}
	if b.withValueDate {
		clause = append(clause, fmt.Sprintf("valueDate: %s", formatValues(b.valueDate, b.operator)))
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

func pathToFilterTarget(path []string) *pb.FilterTarget {
	if len(path) > 2 {
		return &pb.FilterTarget{Target: &pb.FilterTarget_MultiTarget{MultiTarget: &pb.FilterReferenceMultiTarget{
			On:               path[0],
			TargetCollection: path[1],
			Target:           pathToFilterTarget(path[2:]),
		}}}
	}
	return &pb.FilterTarget{Target: &pb.FilterTarget_Property{Property: path[0]}}
}

func (b *WhereBuilder) ToGRPC() *pb.Filters {
	filters := &pb.Filters{}

	switch b.operator {
	case Or:
		filters.Operator = pb.Filters_OPERATOR_OR
		filters.Filters = operandsToGRPC(b.operands)

	case And:
		filters.Operator = pb.Filters_OPERATOR_AND
		filters.Filters = operandsToGRPC(b.operands)

	case Not:
		filters.Operator = pb.Filters_OPERATOR_NOT
		filters.Filters = operandsToGRPC(b.operands)

	default:
		filters.Target = pathToFilterTarget(b.path)

		switch b.operator {
		case Equal:
			filters.Operator = pb.Filters_OPERATOR_EQUAL
		case NotEqual:
			filters.Operator = pb.Filters_OPERATOR_NOT_EQUAL
		case GreaterThan:
			filters.Operator = pb.Filters_OPERATOR_GREATER_THAN
		case GreaterThanEqual:
			filters.Operator = pb.Filters_OPERATOR_GREATER_THAN_EQUAL
		case LessThan:
			filters.Operator = pb.Filters_OPERATOR_LESS_THAN
		case LessThanEqual:
			filters.Operator = pb.Filters_OPERATOR_LESS_THAN_EQUAL
		case Like:
			filters.Operator = pb.Filters_OPERATOR_LIKE
		case IsNull:
			filters.Operator = pb.Filters_OPERATOR_IS_NULL
		case WithinGeoRange:
			filters.Operator = pb.Filters_OPERATOR_WITHIN_GEO_RANGE
		case ContainsAny:
			filters.Operator = pb.Filters_OPERATOR_CONTAINS_ANY
		case ContainsAll:
			filters.Operator = pb.Filters_OPERATOR_CONTAINS_ALL
		case ContainsNone:
			filters.Operator = pb.Filters_OPERATOR_CONTAINS_NONE
		default:
		}

		if ln := len(b.valueInt); ln > 0 {
			if ln > 1 || isContainsOperator(b.operator) {
				filters.TestValue = &pb.Filters_ValueIntArray{ValueIntArray: &pb.IntArray{Values: b.valueInt}}
			} else {
				filters.TestValue = &pb.Filters_ValueInt{ValueInt: b.valueInt[0]}
			}
		}
		if ln := len(b.valueNumber); ln > 0 {
			if ln > 1 || isContainsOperator(b.operator) {
				filters.TestValue = &pb.Filters_ValueNumberArray{ValueNumberArray: &pb.NumberArray{Values: b.valueNumber}}
			} else {
				filters.TestValue = &pb.Filters_ValueNumber{ValueNumber: b.valueNumber[0]}
			}
		}
		if ln := len(b.valueBoolean); ln > 0 {
			if ln > 1 || isContainsOperator(b.operator) {
				filters.TestValue = &pb.Filters_ValueBooleanArray{ValueBooleanArray: &pb.BooleanArray{Values: b.valueBoolean}}
			} else {
				filters.TestValue = &pb.Filters_ValueBoolean{ValueBoolean: b.valueBoolean[0]}
			}
		}
		if ln := len(b.valueText); ln > 0 {
			if ln > 1 || isContainsOperator(b.operator) {
				filters.TestValue = &pb.Filters_ValueTextArray{ValueTextArray: &pb.TextArray{Values: b.valueText}}
			} else {
				filters.TestValue = &pb.Filters_ValueText{ValueText: b.valueText[0]}
			}
		} else if ln := len(b.valueString); ln > 0 {
			if ln > 1 || isContainsOperator(b.operator) {
				filters.TestValue = &pb.Filters_ValueTextArray{ValueTextArray: &pb.TextArray{Values: b.valueString}}
			} else {
				filters.TestValue = &pb.Filters_ValueText{ValueText: b.valueString[0]}
			}
		}
		if ln := len(b.valueDate); ln > 0 {
			formattedDates := formatDates(b.valueDate)
			if ln > 1 || isContainsOperator(b.operator) {
				filters.TestValue = &pb.Filters_ValueTextArray{ValueTextArray: &pb.TextArray{Values: formattedDates}}
			} else {
				filters.TestValue = &pb.Filters_ValueText{ValueText: formattedDates[0]}
			}
		}
		if b.valueGeoRange != nil {
			filters.TestValue = &pb.Filters_ValueGeo{ValueGeo: &pb.GeoCoordinatesFilter{
				Latitude:  b.valueGeoRange.Latitude,
				Longitude: b.valueGeoRange.Longitude,
				Distance:  b.valueGeoRange.MaxDistance,
			}}
		}
	}
	return filters
}

func operandsToGRPC(wheres []*WhereBuilder) []*pb.Filters {
	filters := make([]*pb.Filters, len(wheres))
	for i := range wheres {
		filters[i] = wheres[i].ToGRPC()
	}
	return filters
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

func formatValues[T any](values []T, operator WhereOperator) string {
	clause := make([]string, len(values))
	for i, value := range values {
		switch val := any(value).(type) {
		case string:
			clause[i] = fmt.Sprintf("%q", val)
		case time.Time:
			clause[i] = fmt.Sprintf("%q", val.Format(time.RFC3339Nano))
		default:
			clause[i] = fmt.Sprintf("%v", val)
		}
	}
	formattedValues := strings.Join(clause, ",")
	if len(values) > 1 || isContainsOperator(operator) {
		return fmt.Sprintf("[%s]", formattedValues)
	}
	return formattedValues
}

func formatDates(dates []time.Time) []string {
	formatted := make([]string, len(dates))
	for i := range dates {
		formatted[i] = dates[i].Format(time.RFC3339Nano)
	}
	return formatted
}

func isContainsOperator(operator WhereOperator) bool {
	switch operator {
	case ContainsAny, ContainsAll, ContainsNone:
		return true
	default:
		return false
	}
}
