package internal

// NilZero returns the dereferenced value if the pointer is not nil and the zero value of type T otherwise.
func NilZero[T any](v *T) T {
	var zero T
	if v != nil {
		return *v
	}
	return zero
}
