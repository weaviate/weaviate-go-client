package data

import "github.com/weaviate/weaviate-go-client/v6/internal"

// MustEncode is an Encode wrapper that panics instead of returning an error.
func MustEncode[T any](v *T) map[string]any {
	m, err := Encode(v)
	if err != nil {
		panic(err)
	}
	return m
}

// Encode converts a Go struct into map[string]any.
// Encode supports "json" tag for customizing how
// struct field names get converted to map keys.
func Encode[T any](v *T) (map[string]any, error) {
	return internal.Encode(v)
}
