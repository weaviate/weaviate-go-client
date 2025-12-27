package testkit

import "time"

// Now is constant across the entire test run.
// To make Now comparable with time parsed from a string,
// Now is stripped from any monotonic clock reading.
var Now = time.Now().Round(0)

// Ptr is a helper for passing pointers to constants.
func Ptr[T any](v T) *T { return &v }
