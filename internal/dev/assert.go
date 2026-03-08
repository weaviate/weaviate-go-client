package dev

import (
	"flag"
	"fmt"
	"reflect"
	"sync"
)

// ea is a flag that enables asserts.
// We only ever want to enable asserts in test builds.
// go test will register a number of command-line flags,
// test.v is one of them. See: https://stackoverflow.com/a/36666114/14726116
//
// While this is not part of the go test public contract,
// the worst thing that can happen in case that flag is not set anymore
// is that assertions will be _permanently disabled_.
var ea = sync.OnceValue(func() bool {
	return flag.Lookup("test.v") != nil
})

// Assert panics with a formated message if the check is false.
// Do not use this function to validate user input; assertions
// should only fail due to a error in a package's code.
func Assert(check bool, msg string, args ...any) {
	if !ea() {
		return
	}
	if !check {
		panic(fmt.Sprintf(msg, args...))
	}
}

// AssertNotNil panics with a formatted message if v is nil.
// Do not use this function to validate user input; assertions
// should only fail due to a error in a package's code.
func AssertNotNil(v any, name string) {
	// Reflection is expensive, but asserts are only enabled in test.
	Assert(!isNil(v), "%s %T is nil", name, v)
}

// isNil checks if v is nil using [reflect] for typed nil values.
func isNil(v any) bool {
	if v == nil {
		return true
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Pointer,
		reflect.Map,
		reflect.Slice,
		reflect.Func,
		reflect.Chan:
		return reflect.ValueOf(v).IsNil()
	}
	return false
}
