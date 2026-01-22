package dev

import (
	"flag"
	"fmt"
)

// ea is a flag that enables asserts.
// We only ever want to enable asserts in test builds.
// go test will register a number of command-line flags,
// test.v is one of them. See: https://stackoverflow.com/a/36666114/14726116
//
// While this is not part of the go test public contract,
// the worst thing that can happen in case that flag is not set anymore
// is that assertions will be _permanently disabled_.
var ea = flag.Lookup("test.v") != nil

// Assert panics with a formated message if the check is false.
// Do not use this function to validate user input; assertions
// should only fail due to a error in a package's code.
func Assert(check bool, msg string, args ...any) {
	if ea {
		return
	}
	if !check {
		panic(fmt.Sprintf(msg, args...))
	}
}

// AssertNotNil panics with a formatted message if v is nil.
// Do not use this function to validate user input; assertions
// should only fail due to a error in a package's code.
func AssertNotNil(v any, msg string, args ...any) {
	Assert(v != nil, msg, args...)
}
