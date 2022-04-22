package graphql

// ExploreFields used in an Explore GraphQL query
type ExploreFields string

// Certainty includes the certainty that a data object is related to the query concept
const Certainty ExploreFields = "certainty"

// Beacon includes the beacon to the found objects
const Beacon ExploreFields = "beacon"

// ClassName includes the class name of the found objects
const ClassName ExploreFields = "className"

// MoveParameters to fine tune Explore queries
type MoveParameters struct {
	// Concepts that should be used as base for the movement operation
	Concepts []string
	// Force to be applied in the movement operation
	Force float32
}

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

// GeoCoordinatesParameter parameters in where filter
type GeoCoordinatesParameter struct {
	Latitude, Longitude, MaxDistance float32
}

// SortOrder used in Sort Argument builder
type SortOrder string

// Asc ascending sort order
const Asc SortOrder = "asc"

// Desc descending sort order
const Desc SortOrder = "desc"
