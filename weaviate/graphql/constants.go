package graphql

// ExploreFields used in an Explore GraphQL query
type ExploreFields string

const (
	// Certainty includes the certainty that a data object is related to the query concept
	Certainty ExploreFields = "certainty"

	// Distance includes the distance that a data object is related to the query concept
	Distance ExploreFields = "distance"

	// Beacon includes the beacon to the found objects
	Beacon ExploreFields = "beacon"

	// ClassName includes the class name of the found objects
	ClassName ExploreFields = "className"
)

// EmptyObjectStr string representation of an empty object
const EmptyObjectStr string = "{}"

// SortOrder used in Sort Argument builder
type SortOrder string

const (
	// Asc ascending sort order
	Asc SortOrder = "asc"

	// Desc descending sort order
	Desc SortOrder = "desc"
)
