package internal

type CombinationMethod string

const (
	CombinationMethodSum           CombinationMethod = "SUM"
	CombinationMethodMax           CombinationMethod = "MAX"
	CombinationMethodAverage       CombinationMethod = "AVERAGE"
	CombinationMethodManualWeights CombinationMethod = "MANUAL_WEIGHTS"
	CombinationMethodRelativeScore CombinationMethod = "RELATIVE_SCORE"
)
