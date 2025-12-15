package dev

import "fmt"

// Assert panics with a formated message if the check is false.
// Do not use this function to validate user input; assertions
// should only fail due to a error in a package's code.
func Assert(check bool, msg string, args ...any) {
	if !check {
		panic(fmt.Sprintf(msg, args...))
	}
}

// AssertType panics if v is not of type T.
// A nil intput is returned as typed nil without a type assertion.
// Do not use this function to validate user input; assertions
// should only fail due to a error in a package's code.
func AssertType[T any](v any) T {
	t, ok := v.(T)
	if v == nil {
		return t
	}
	Assert(ok, "value must be %T, got %T", *new(T), v)
	return t
}
