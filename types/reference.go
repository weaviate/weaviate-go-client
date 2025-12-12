package types

// Reference represents a cross-reference to another object.
type Reference struct {
	// Collection the refernced object belongs to. Does not need to be set for single-target queries.
	Collection string

	// UUID of the referenced object.
	UUID string
}
