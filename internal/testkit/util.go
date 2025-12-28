package testkit

import "time"

// Now is constant across the entire test run.
// To make Now comparable with time parsed from a string,
// Now is stripped from any monotonic clock reading and
// has it's Location set to nil.
//
// This is only important in test where we use stretchr/testify
// packages, which do not [compare time] correctly.
//
// [compare time]: https://github.com/stretchr/testify/issues/502
var Now = time.Now().UTC()

// Ptr is a helper for passing pointers to constants.
func Ptr[T any](v T) *T { return &v }
