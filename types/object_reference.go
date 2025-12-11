package types

// ObjectReference represents a cross-reference to another object.
type ObjectReference struct {
	// TargetCollection is the collection of the referenced object.
	TargetCollection string

	// TargetID is the ID of the referenced object.
	TargetID string
}
