package testkit

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Now is constant across the entire test run.
// To make time comparable with time parsed from a string,
// Now is has no ns precision and always uses the local TZ.
//
// This is only important in test where we use stretchr/testify
// packages, which do not [compare time] correctly.
//
// [compare time]: https://github.com/stretchr/testify/issues/502
var Now = time.Date(6, time.Month(5), 4, 3, 2, 1, 0, time.Local)

// UUID is a stub UUID tests can use to verify the correct UUID is used.
var UUID = uuid.New()

// Ptr is a helper for passing pointers to constants.
func Ptr[T any](v T) *T { return &v }

// RequirePointer asserts that v is a pointer. If the assertion fails,
// the test t will fail immediately. Use RequirePointer as a pre-condition
// in unit tests to ensure the test cases are valid.
func RequirePointer(t *testing.T, v any, name string) {
	t.Helper()
	switch reflect.TypeOf(v).Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.Pointer:
		return
	default:
		require.FailNowf(t, "not a pointer", "%q must be a pointer", name)
	}
}
