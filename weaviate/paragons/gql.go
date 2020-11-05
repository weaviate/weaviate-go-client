package paragons

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

