package internal

// MakeMap returns nil if size is 0.
//
// This works well with custom types derived from map.
// Namely, the output of MakeMap can be assigned to a
// variable or field of that type without an explicit cast.
//
// More that anything, this simplifies the testing setup,
// as we no longer need to pass a make(map[K]V) in every
// test stub where we don't expect any results.
func MakeMap[K comparable, V any](size int) map[K]V {
	if size == 0 {
		return nil
	}
	return make(map[K]V, size)
}
