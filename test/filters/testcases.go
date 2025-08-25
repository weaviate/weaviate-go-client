package filtertestcases

import (
	"fmt"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/filters"
)

type FilterTestCase struct {
	Name        string
	Where       *filters.WhereBuilder
	Property    string
	ExpectedIds []string
}

type AllFilterTestCases struct {
	Contains []FilterTestCase
	Equal    []FilterTestCase
	Greater  []FilterTestCase
	Less     []FilterTestCase
	Like     []FilterTestCase
	OrAnd    []FilterTestCase
}

func AllPropsTestCases() AllFilterTestCases {
	data := testsuit.AllPropertiesData()

	id1, id2, id3 := data.ID1, data.ID2, data.ID3
	invertIds := createInvertIds(data.IDs)

	containsTestCases := []FilterTestCase{
		// Contains operator with array types
		{
			Name: "contains all authors with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.ContainsAll).
				WithValueString("John", "Jenny", "Joseph"),
			Property:    "authors",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains any authors with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.ContainsAny).
				WithValueString("John", "Jenny", "Joseph"),
			Property:    "authors",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains none authors with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.ContainsNone).
				WithValueString("Joseph", "Missing"),
			Property:    "authors",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "contains all colors with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.ContainsAll).
				WithValueText("red", "blue", "green"),
			Property:    "colors",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains any colors with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.ContainsAny).
				WithValueText("red", "blue", "green"),
			Property:    "colors",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains none colors with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.ContainsNone).
				WithValueText("green", "missing"),
			Property:    "colors",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "contains all numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.ContainsAll).
				WithValueNumber(1.1, 2.2, 3.3),
			Property:    "numbers",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains any numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.ContainsAny).
				WithValueNumber(1.1, 2.2, 3.3),
			Property:    "numbers",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains none numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.ContainsNone).
				WithValueNumber(3.3, 0),
			Property:    "numbers",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "contains all ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.ContainsAll).
				WithValueInt(1, 2, 3),
			Property:    "ints",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains any ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.ContainsAny).
				WithValueInt(1, 2, 3),
			Property:    "ints",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains none ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.ContainsNone).
				WithValueInt(3, 0),
			Property:    "ints",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "contains all bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.ContainsAll).
				WithValueBoolean(true, false, true),
			Property:    "bools",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "contains any bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.ContainsAny).
				WithValueBoolean(true, false),
			Property:    "bools",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains none bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.ContainsNone).
				WithValueBoolean(false),
			Property:    "bools",
			ExpectedIds: []string{id3},
		},
		{
			Name: "contains all uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.ContainsAll).
				WithValueText(id1, id2, id3),
			Property:    "uuids",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains any uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.ContainsAny).
				WithValueText(id1, id2, id3),
			Property:    "uuids",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains none uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.ContainsNone).
				WithValueText(id3, "FFFFFFFF-FFFF-0000-0000-000000000000"),
			Property:    "uuids",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "contains all dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.ContainsAll).
				WithValueDate(
					mustGetTime("2009-11-01T23:00:00Z"),
					mustGetTime("2009-11-02T23:00:00Z"),
					mustGetTime("2009-11-03T23:00:00Z"),
				),
			Property:    "dates",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains any dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.ContainsAny).
				WithValueDate(
					mustGetTime("2009-11-01T23:00:00Z"),
					mustGetTime("2009-11-02T23:00:00Z"),
					mustGetTime("2009-11-03T23:00:00Z"),
				),
			Property:    "dates",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains none dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.ContainsNone).
				WithValueDate(
					mustGetTime("2009-11-03T23:00:00Z"),
					mustGetTime("1970-01-01T00:00:00Z"),
				),
			Property:    "dates",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "complex contains all ints and all numbers with AND on int array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"numbers"}).
						WithOperator(filters.ContainsAll).
						WithValueNumber(1.1, 2.2, 3.3),
					filters.Where().
						WithPath([]string{"ints"}).
						WithOperator(filters.ContainsAll).
						WithValueInt(1, 2, 3),
				}),
			Property:    "ints",
			ExpectedIds: []string{id1},
		},
		{
			Name: "complex contains any ints and all numbers and none texts with OR",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"numbers"}).
						WithOperator(filters.ContainsAll).
						WithValueNumber(1.1, 2.2, 3.3),
					filters.Where().
						WithPath([]string{"ints"}).
						WithOperator(filters.ContainsAny).
						WithValueInt(3),
					filters.Where().
						WithPath([]string{"authors"}).
						WithOperator(filters.ContainsNone).
						WithValueString("Jenny", "Missing"),
				}),
			Property:    "ints",
			ExpectedIds: []string{id1, id3},
		},
		// Contains operator with primitives
		{
			Name: "contains any author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.ContainsAny).
				WithValueString("John", "Jenny", "Joseph"),
			Property:    "author",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains any color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.ContainsAny).
				WithValueText("red", "blue", "green"),
			Property:    "color",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains any number with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.ContainsAny).
				WithValueNumber(1.1, 2.2, 3.3),
			Property:    "number",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains any int with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.ContainsAny).
				WithValueInt(1, 2, 3),
			Property:    "int",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains any bool with bool",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.ContainsAny).
				WithValueBoolean(true, false, true),
			Property:    "bool",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains any uuid with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.ContainsAny).
				WithValueText(id1, id2, id3),
			Property:    "uuid",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains any uuid with id",
			Where: filters.Where().
				WithPath([]string{"_id"}).
				WithOperator(filters.ContainsAny).
				WithValueText(id1, id2, id3),
			Property:    "uuid",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "contains any date with date",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.ContainsAny).
				WithValueDate(
					mustGetTime("2009-11-01T23:00:00Z"), mustGetTime("2009-11-02T23:00:00Z"), mustGetTime("2009-11-03T23:00:00Z"),
				),
			Property:    "date",
			ExpectedIds: []string{id1, id2, id3},
		},

		{
			Name: "contains all author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.ContainsAll).
				WithValueString("Jenny"),
			Property:    "author",
			ExpectedIds: []string{id2},
		},
		{
			Name: "contains all color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.ContainsAll).
				WithValueText("blue"),
			Property:    "color",
			ExpectedIds: []string{id2},
		},
		{
			Name: "contains all number with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.ContainsAll).
				WithValueNumber(2.2),
			Property:    "number",
			ExpectedIds: []string{id2},
		},
		{
			Name: "contains all int with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.ContainsAll).
				WithValueInt(2),
			Property:    "int",
			ExpectedIds: []string{id2},
		},
		{
			Name: "contains all bool with bool",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.ContainsAll).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{id2},
		},
		{
			Name: "contains all uuid with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.ContainsAll).
				WithValueText(id2),
			Property:    "uuid",
			ExpectedIds: []string{id2},
		},
		{
			Name: "contains all date with date",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.ContainsAll).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "date",
			ExpectedIds: []string{id2},
		},

		{
			Name: "contains none author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.ContainsNone).
				WithValueString("Jenny", "Joseph"),
			Property:    "author",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains none color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.ContainsNone).
				WithValueText("blue", "green"),
			Property:    "color",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains none number with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.ContainsNone).
				WithValueNumber(2.2, 3.3),
			Property:    "number",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains none int with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.ContainsNone).
				WithValueInt(2, 3),
			Property:    "int",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains none bool with bool",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.ContainsNone).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "contains none uuid with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.ContainsNone).
				WithValueText(id2, id3),
			Property:    "uuid",
			ExpectedIds: []string{id1},
		},
		{
			Name: "contains none date with date",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.ContainsNone).
				WithValueDate(
					mustGetTime("2009-11-02T23:00:00Z"),
					mustGetTime("2009-11-03T23:00:00Z"),
				),
			Property:    "date",
			ExpectedIds: []string{id1},
		},
	}

	equalTestCases := []FilterTestCase{
		// arrays
		{
			Name: "equal author with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.Equal).
				WithValueString("Jenny"),
			Property:    "authors",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "notEqual author with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.NotEqual).
				WithValueString("Jenny"),
			Property:    "authors",
			ExpectedIds: []string{id3},
		},
		{
			Name: "equal color with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.Equal).
				WithValueText("blue"),
			Property:    "colors",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "notEqual color with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.NotEqual).
				WithValueText("blue"),
			Property:    "colors",
			ExpectedIds: []string{id3},
		},
		{
			Name: "equal numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.Equal).
				WithValueNumber(2.2),
			Property:    "numbers",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "notEqual numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.NotEqual).
				WithValueNumber(2.2),
			Property:    "numbers",
			ExpectedIds: []string{id3},
		},
		{
			Name: "equal ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.Equal).
				WithValueInt(2),
			Property:    "ints",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "notEqual ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.NotEqual).
				WithValueInt(2),
			Property:    "ints",
			ExpectedIds: []string{id3},
		},
		{
			Name: "equal bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.Equal).
				WithValueBoolean(false),
			Property:    "bools",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "notEqual bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.NotEqual).
				WithValueBoolean(false),
			Property:    "bools",
			ExpectedIds: []string{id3},
		},
		{
			Name: "equal uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.Equal).
				WithValueText(id2),
			Property:    "uuids",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "notEqual uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.NotEqual).
				WithValueText(id2),
			Property:    "uuids",
			ExpectedIds: []string{id3},
		},
		{
			Name: "equal dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.Equal).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "dates",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "notEqual dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.NotEqual).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "dates",
			ExpectedIds: []string{id3},
		},

		// primitives
		{
			Name: "equal author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.Equal).
				WithValueString("Jenny"),
			Property:    "author",
			ExpectedIds: []string{id2},
		},
		{
			Name: "notEqual author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.NotEqual).
				WithValueString("Jenny"),
			Property:    "author",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "equal color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.Equal).
				WithValueText("blue"),
			Property:    "color",
			ExpectedIds: []string{id2},
		},
		{
			Name: "notEqual color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.NotEqual).
				WithValueText("blue"),
			Property:    "color",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "equal numbers with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.Equal).
				WithValueNumber(2.2),
			Property:    "number",
			ExpectedIds: []string{id2},
		},
		{
			Name: "notEqual numbers with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.NotEqual).
				WithValueNumber(2.2),
			Property:    "number",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "equal ints with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.Equal).
				WithValueInt(2),
			Property:    "int",
			ExpectedIds: []string{id2},
		},
		{
			Name: "notEqual ints with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.NotEqual).
				WithValueInt(2),
			Property:    "int",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "equal bools with bool",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.Equal).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{id2},
		},
		{
			Name: "notEqual bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.NotEqual).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "equal uuids with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.Equal).
				WithValueText(id2),
			Property:    "uuid",
			ExpectedIds: []string{id2},
		},
		{
			Name: "notEqual uuids with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.NotEqual).
				WithValueText(id2),
			Property:    "uuid",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "equal dates with dates",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.Equal).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "date",
			ExpectedIds: []string{id2},
		},
		{
			Name: "notEqual dates with dates",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.NotEqual).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "date",
			ExpectedIds: []string{id1, id3},
		},
	}

	greaterTestCases := []FilterTestCase{
		// arrays
		{
			Name: "greaterThanEqual author with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueString("John"),
			Property:    "authors",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "greaterThan author with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.GreaterThan).
				WithValueString("John"),
			Property:    "authors",
			ExpectedIds: []string{id1},
		},
		{
			Name: "greaterThanEqual color with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueText("blue"),
			Property:    "colors",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "greaterThan color with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.GreaterThan).
				WithValueText("blue"),
			Property:    "colors",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "greaterThanEqual numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueNumber(2.2),
			Property:    "numbers",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "greaterThan numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.GreaterThan).
				WithValueNumber(2.2),
			Property:    "numbers",
			ExpectedIds: []string{id1},
		},
		{
			Name: "greaterThanEqual ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueInt(2),
			Property:    "ints",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "greaterThan ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.GreaterThan).
				WithValueInt(2),
			Property:    "ints",
			ExpectedIds: []string{id1},
		},
		{
			Name: "greaterThanEqual bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueBoolean(false),
			Property:    "bools",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "greaterThan bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.GreaterThan).
				WithValueBoolean(false),
			Property:    "bools",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "greaterThanEqual uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueText(id2),
			Property:    "uuids",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "greaterThan uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.GreaterThan).
				WithValueText(id2),
			Property:    "uuids",
			ExpectedIds: []string{id1},
		},
		{
			Name: "greaterThanEqual dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "dates",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "greaterThan dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.GreaterThan).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "dates",
			ExpectedIds: []string{id1},
		},

		// primitives
		{
			Name: "greaterThanEqual author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueString("John"),
			Property:    "author",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "greaterThan author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.GreaterThan).
				WithValueString("John"),
			Property:    "author",
			ExpectedIds: []string{id3},
		},
		{
			Name: "greaterThanEqual color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueText("blue"),
			Property:    "color",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "greaterThan color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.GreaterThan).
				WithValueText("blue"),
			Property:    "color",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "greaterThanEqual numbers with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueNumber(2.2),
			Property:    "number",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "greaterThan numbers with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.GreaterThan).
				WithValueNumber(2.2),
			Property:    "number",
			ExpectedIds: []string{id3},
		},
		{
			Name: "greaterThanEqual ints with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueInt(2),
			Property:    "int",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "greaterThan ints with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.GreaterThan).
				WithValueInt(2),
			Property:    "int",
			ExpectedIds: []string{id3},
		},
		{
			Name: "greaterThanEqual bools with bool",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "greaterThan bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.GreaterThan).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "greaterThanEqual uuids with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueText(id2),
			Property:    "uuid",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "greaterThan uuids with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.GreaterThan).
				WithValueText(id2),
			Property:    "uuid",
			ExpectedIds: []string{id3},
		},
		{
			Name: "greaterThanEqual dates with dates",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.GreaterThanEqual).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "date",
			ExpectedIds: []string{id2, id3},
		},
		{
			Name: "greaterThan dates with dates",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.GreaterThan).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "date",
			ExpectedIds: []string{id3},
		},
	}

	lessTestCases := []FilterTestCase{
		// arrays
		{
			Name: "lessThanEqual author with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.LessThanEqual).
				WithValueString("John"),
			Property:    "authors",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThan author with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.LessThan).
				WithValueString("John"),
			Property:    "authors",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThanEqual color with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.LessThanEqual).
				WithValueText("blue"),
			Property:    "colors",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThan color with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.LessThan).
				WithValueText("blue"),
			Property:    "colors",
			ExpectedIds: []string{},
		},
		{
			Name: "lessThanEqual numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.LessThanEqual).
				WithValueNumber(2.2),
			Property:    "numbers",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThan numbers with number array",
			Where: filters.Where().
				WithPath([]string{"numbers"}).
				WithOperator(filters.LessThan).
				WithValueNumber(2.2),
			Property:    "numbers",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThanEqual ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.LessThanEqual).
				WithValueInt(2),
			Property:    "ints",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThan ints with int array",
			Where: filters.Where().
				WithPath([]string{"ints"}).
				WithOperator(filters.LessThan).
				WithValueInt(2),
			Property:    "ints",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThanEqual bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.LessThanEqual).
				WithValueBoolean(false),
			Property:    "bools",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThan bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bools"}).
				WithOperator(filters.LessThan).
				WithValueBoolean(false),
			Property:    "bools",
			ExpectedIds: []string{},
		},
		{
			Name: "lessThanEqual uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.LessThanEqual).
				WithValueText(id2),
			Property:    "uuids",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThan uuids with uuid array",
			Where: filters.Where().
				WithPath([]string{"uuids"}).
				WithOperator(filters.LessThan).
				WithValueText(id2),
			Property:    "uuids",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThanEqual dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.LessThanEqual).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "dates",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "lessThan dates with dates array",
			Where: filters.Where().
				WithPath([]string{"dates"}).
				WithOperator(filters.LessThan).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "dates",
			ExpectedIds: []string{id1, id2, id3},
		},

		// primitives
		{
			Name: "lessThanEqual author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.LessThanEqual).
				WithValueString("John"),
			Property:    "author",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThan author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.LessThan).
				WithValueString("John"),
			Property:    "author",
			ExpectedIds: []string{id2},
		},
		{
			Name: "lessThanEqual color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.LessThanEqual).
				WithValueText("blue"),
			Property:    "color",
			ExpectedIds: []string{id2},
		},
		{
			Name: "lessThan color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.LessThan).
				WithValueText("blue"),
			Property:    "color",
			ExpectedIds: []string{},
		},
		{
			Name: "lessThanEqual numbers with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.LessThanEqual).
				WithValueNumber(2.2),
			Property:    "number",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThan numbers with number",
			Where: filters.Where().
				WithPath([]string{"number"}).
				WithOperator(filters.LessThan).
				WithValueNumber(2.2),
			Property:    "number",
			ExpectedIds: []string{id1},
		},
		{
			Name: "lessThanEqual ints with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.LessThanEqual).
				WithValueInt(2),
			Property:    "int",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThan ints with int",
			Where: filters.Where().
				WithPath([]string{"int"}).
				WithOperator(filters.LessThan).
				WithValueInt(2),
			Property:    "int",
			ExpectedIds: []string{id1},
		},
		{
			Name: "lessThanEqual bools with bool",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.LessThanEqual).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{id2},
		},
		{
			Name: "lessThan bools with bool array",
			Where: filters.Where().
				WithPath([]string{"bool"}).
				WithOperator(filters.LessThan).
				WithValueBoolean(false),
			Property:    "bool",
			ExpectedIds: []string{},
		},
		{
			Name: "lessThanEqual uuids with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.LessThanEqual).
				WithValueText(id2),
			Property:    "uuid",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThan uuids with uuid",
			Where: filters.Where().
				WithPath([]string{"uuid"}).
				WithOperator(filters.LessThan).
				WithValueText(id2),
			Property:    "uuid",
			ExpectedIds: []string{id1},
		},
		{
			Name: "lessThanEqual dates with dates",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.LessThanEqual).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "date",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "lessThan dates with dates",
			Where: filters.Where().
				WithPath([]string{"date"}).
				WithOperator(filters.LessThan).
				WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
			Property:    "date",
			ExpectedIds: []string{id1},
		},
	}

	likeTestCases := []FilterTestCase{
		// arrays
		{
			Name: "like author with string array",
			Where: filters.Where().
				WithPath([]string{"authors"}).
				WithOperator(filters.Like).
				WithValueString("Jo*"),
			Property:    "authors",
			ExpectedIds: []string{id1, id2, id3},
		},
		{
			Name: "like color with text array",
			Where: filters.Where().
				WithPath([]string{"colors"}).
				WithOperator(filters.Like).
				WithValueText("*lu*"),
			Property:    "colors",
			ExpectedIds: []string{id1, id2},
		},

		// primitives
		{
			Name: "like author with string",
			Where: filters.Where().
				WithPath([]string{"author"}).
				WithOperator(filters.Like).
				WithValueString("Jo*"),
			Property:    "author",
			ExpectedIds: []string{id1, id3},
		},
		{
			Name: "like color with text",
			Where: filters.Where().
				WithPath([]string{"color"}).
				WithOperator(filters.Like).
				WithValueText("*lu*"),
			Property:    "color",
			ExpectedIds: []string{id2},
		},
	}

	orAndTestCases := []FilterTestCase{
		{
			Name: "OR author with string array",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"authors"}).
						WithOperator(filters.Equal).
						WithValueString("Jenny"),
					filters.Where().
						WithPath([]string{"author"}).
						WithOperator(filters.Equal).
						WithValueString("Jenny"),
				}),
			Property:    "authors",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "AND author with string array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"authors"}).
						WithOperator(filters.Equal).
						WithValueString("Jenny"),
					filters.Where().
						WithPath([]string{"author"}).
						WithOperator(filters.Equal).
						WithValueString("Jenny"),
				}),
			Property:    "authors",
			ExpectedIds: []string{id2},
		},
		{
			Name: "OR color with text array",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"colors"}).
						WithOperator(filters.Equal).
						WithValueText("blue"),
					filters.Where().
						WithPath([]string{"color"}).
						WithOperator(filters.Equal).
						WithValueText("blue"),
				}),
			Property:    "colors",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "AND color with text array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"colors"}).
						WithOperator(filters.Equal).
						WithValueText("blue"),
					filters.Where().
						WithPath([]string{"color"}).
						WithOperator(filters.Equal).
						WithValueText("blue"),
				}),
			Property:    "colors",
			ExpectedIds: []string{id2},
		},
		{
			Name: "OR numbers with number array",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"numbers"}).
						WithOperator(filters.Equal).
						WithValueNumber(2.2),
					filters.Where().
						WithPath([]string{"number"}).
						WithOperator(filters.Equal).
						WithValueNumber(2.2),
				}),
			Property:    "numbers",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "AND numbers with number array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"numbers"}).
						WithOperator(filters.Equal).
						WithValueNumber(2.2),
					filters.Where().
						WithPath([]string{"number"}).
						WithOperator(filters.Equal).
						WithValueNumber(2.2),
				}),
			Property:    "numbers",
			ExpectedIds: []string{id2},
		},
		{
			Name: "OR ints with int array",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"ints"}).
						WithOperator(filters.Equal).
						WithValueInt(2),
					filters.Where().
						WithPath([]string{"int"}).
						WithOperator(filters.Equal).
						WithValueInt(2),
				}),
			Property:    "ints",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "AND ints with int array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"ints"}).
						WithOperator(filters.Equal).
						WithValueInt(2),
					filters.Where().
						WithPath([]string{"int"}).
						WithOperator(filters.Equal).
						WithValueInt(2),
				}),
			Property:    "ints",
			ExpectedIds: []string{id2},
		},
		{
			Name: "OR bools with bool array",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"bools"}).
						WithOperator(filters.Equal).
						WithValueBoolean(false),
					filters.Where().
						WithPath([]string{"bool"}).
						WithOperator(filters.Equal).
						WithValueBoolean(false),
				}),
			Property:    "bools",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "AND bools with bool array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"bools"}).
						WithOperator(filters.Equal).
						WithValueBoolean(false),
					filters.Where().
						WithPath([]string{"bool"}).
						WithOperator(filters.Equal).
						WithValueBoolean(false),
				}),
			Property:    "bools",
			ExpectedIds: []string{id2},
		},
		{
			Name: "OR uuids with uuid array",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"uuids"}).
						WithOperator(filters.Equal).
						WithValueText(id2),
					filters.Where().
						WithPath([]string{"uuid"}).
						WithOperator(filters.Equal).
						WithValueText(id2),
				}),
			Property:    "uuids",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "AND uuids with uuid array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"uuids"}).
						WithOperator(filters.Equal).
						WithValueText(id2),
					filters.Where().
						WithPath([]string{"uuid"}).
						WithOperator(filters.Equal).
						WithValueText(id2),
				}),
			Property:    "uuids",
			ExpectedIds: []string{id2},
		},
		{
			Name: "OR dates with dates array",
			Where: filters.Where().
				WithOperator(filters.Or).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"dates"}).
						WithOperator(filters.Equal).
						WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
					filters.Where().
						WithPath([]string{"date"}).
						WithOperator(filters.Equal).
						WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
				}),
			Property:    "dates",
			ExpectedIds: []string{id1, id2},
		},
		{
			Name: "AND dates with dates array",
			Where: filters.Where().
				WithOperator(filters.And).
				WithOperands([]*filters.WhereBuilder{
					filters.Where().
						WithPath([]string{"dates"}).
						WithOperator(filters.Equal).
						WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
					filters.Where().
						WithPath([]string{"date"}).
						WithOperator(filters.Equal).
						WithValueDate(mustGetTime("2009-11-02T23:00:00Z")),
				}),
			Property:    "dates",
			ExpectedIds: []string{id2},
		},
	}

	return AllFilterTestCases{
		Contains: append(containsTestCases, createNotTestCases(containsTestCases, invertIds)...),
		Equal:    append(equalTestCases, createNotTestCases(equalTestCases, invertIds)...),
		Greater:  append(greaterTestCases, createNotTestCases(greaterTestCases, invertIds)...),
		Less:     append(lessTestCases, createNotTestCases(lessTestCases, invertIds)...),
		Like:     append(likeTestCases, createNotTestCases(likeTestCases, invertIds)...),
		OrAnd:    append(orAndTestCases, createNotTestCases(orAndTestCases, invertIds)...),
	}
}

func mustGetTime(date string) time.Time {
	parsed, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		panic(fmt.Sprintf("can't parse date: %v", date))
	}
	return parsed
}

func createInvertIds(allIds []string) func(ids []string) (inverted []string) {
	return func(ids []string) []string {
		existing := make(map[string]struct{}, len(ids))
		for i := range ids {
			existing[ids[i]] = struct{}{}
		}
		invIds := []string{}
		for i := range allIds {
			if _, ok := existing[allIds[i]]; !ok {
				invIds = append(invIds, allIds[i])
			}
		}
		return invIds
	}
}

func createNotTestCases(testcases []FilterTestCase, invertIds func([]string) []string) []FilterTestCase {
	not := make([]FilterTestCase, len(testcases))
	for i := range testcases {
		not[i] = FilterTestCase{
			Name: fmt.Sprintf("NOT %s", testcases[i].Name),
			Where: filters.Where().
				WithOperator(filters.Not).
				WithOperands([]*filters.WhereBuilder{testcases[i].Where}),
			Property:    testcases[i].Property,
			ExpectedIds: invertIds(testcases[i].ExpectedIds),
		}
	}
	return not
}
