package testkit

// Ptr is a helper for passing pointers to constants.
func Ptr[T any](v T) *T { return &v }
