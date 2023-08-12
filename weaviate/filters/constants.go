package filters

// WhereOperator used in Where Argument builder
type WhereOperator string

// And where operator
const And WhereOperator = "And"

// Like where operator
const Like WhereOperator = "Like"

// Or where operator
const Or WhereOperator = "Or"

// Equal where operator
const Equal WhereOperator = "Equal"

// Not where operator
const Not WhereOperator = "Not"

// NotEqual where operator
const NotEqual WhereOperator = "NotEqual"

// GreaterThan where operator
const GreaterThan WhereOperator = "GreaterThan"

// GreaterThanEqual where operator
const GreaterThanEqual WhereOperator = "GreaterThanEqual"

// LessThan where operator
const LessThan WhereOperator = "LessThan"

// LessThanEqual where operator
const LessThanEqual WhereOperator = "LessThanEqual"

// WithinGeoRange where operator
const WithinGeoRange WhereOperator = "WithinGeoRange"

// IsNull where operator
const IsNull WhereOperator = "IsNull"

// ContainsAny where operator
const ContainsAny WhereOperator = "ContainsAny"

// ContainsAll where operator
const ContainsAll WhereOperator = "ContainsAll"

// GeoCoordinatesParameter parameters in where filter
type GeoCoordinatesParameter struct {
	Latitude, Longitude, MaxDistance float32
}
