package types

// Reference represents a cross-reference to another object.
type Reference struct {
	// Collection the refernced object belongs to.
	Collection string

	// UUID of the referenced object.
	UUID string
}
