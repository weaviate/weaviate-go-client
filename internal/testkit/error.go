package testkit

import (
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

// ErrWhaam is a stub error tests can use to verify some error is being propagated.
// The error message is an allusion to [Roy Lichtenstein's dyptich].
//
// [Roy Lichtenstein's dyptich]: https://en.wikipedia.org/wiki/Whaam!
var ErrWhaam = errors.New("Whaam!") // nolint:staticcheck

// Error adds an optional error check to a table-test case.
// A nil Error expects no error, so leaving the field unset
// works for all "happy" cases.
// A non-nil error expects an error to be present.
//
// Depending on the method used, Error will delegate the check
// to either testify/require or testify/assert packages.
//
// Example:
//
//	for _, tt := range []struct{
//		name string
//		act func() error
//		err testkit.Error
//	}{
//
//		{
//			name: "happy case",
//			act: func() error { return nil },
//		},
//		{
//			name: "expect error",
//			act: func() error { return nil },
//			err: testkit.ExpectError
//		},
//	}{
//
//		t.Run(tt.name, func(t *testing.T) {
//			err := tt.act()
//			if !tt.err.Assert(t, err) { // Use assertion result
//				t.Log("assertion failed")
//			}
//			tt.err.Require(t, err) // Fail immediately
//		})
//	}
type Error assert.ErrorAssertionFunc

// ExpectError checks that the error is not nil.
var ExpectError Error = assert.Error

// Assert returns the result of assert f if it is not nil and [assert.NoError] otherwise.
func (e Error) Assert(t *testing.T, err error, msgAndArgs ...any) bool {
	t.Helper()
	if e == nil {
		return assert.NoError(t, err, msgAndArgs...)
	}
	return e(t, err, msgAndArgs...)
}

// Require calls [require.NoError] if f is nil. Otherwise, uses the assertion
// and fails the test immediately if it returns false.
func (e Error) Require(t *testing.T, err error, msgAndArgs ...any) {
	t.Helper()
	if e == nil {
		require.NoError(t, err, msgAndArgs...)
	} else if !e(t, err, msgAndArgs...) {
		t.FailNow()
	}
}
