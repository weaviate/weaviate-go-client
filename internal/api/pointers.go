package api

// ptr is a helper for passing pointers to constants.
func ptr[T any](v T) *T { return &v }

// nilPresent returns a pointer to v if present == true and nil otherwise.
func nilPresent[T any](v T, present bool) *T {
	if !present {
		return nil
	}
	return &v
}
