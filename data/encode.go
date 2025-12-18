package data

import "github.com/weaviate/weaviate-go-client/v6/internal"

// Encode converts a Go struct into map[string]any.
// Encode supports "json" tag for customizing how
// struct field names get converted to map keys.
func Encode[T any](v *T) (map[string]any, error) {
	return internal.Encode(v)
}
