package collections

type Collection struct {
	Name       string
	Properties map[string]Property
	References map[string]ReferenceProperty
}

// DataType defines supported property data types.
type DataType string

const (
	DataTypeText           DataType = "text"
	DataTypeBool           DataType = "boolean"
	DataTypeInt            DataType = "int"
	DataTypeNumber         DataType = "number"
	DataTypeDate           DataType = "date"
	DataTypeObject         DataType = "object"
	DataTypeGeoCoordinates DataType = "geoCoordinates"
	DataTypeTextArray      DataType = "text[]"
	DataTypeBoolArray      DataType = "boolean[]"
	DataTypeIntArray       DataType = "number[]"
	DataTypeNumberArray    DataType = "date[]"
	DataTypeDateArray      DataType = "object[]"
	DataTypeObjectArray    DataType = "geoCoordinates[]"
)

type Property struct {
	Name string
	Type DataType
}

type ReferenceProperty struct {
	Name        string
	Collections []string
}
