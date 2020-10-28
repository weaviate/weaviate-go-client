package paragons

type ExploreFields string

const Certainty ExploreFields = "certainty"
const Beacon ExploreFields = "beacon"
const ClassName ExploreFields = "className"

type MoveParameters struct {
	Concepts []string
	Force float32
}

