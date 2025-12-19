package internal

import "github.com/go-viper/mapstructure/v2"

// tagName is the tag mapstructure will use to match struct fields to keys in the properties map.
const tagName = "json"

// Decode is a thin wrapper around mapstructure.Decode
// that decodes map[string]any into a Go struct.
// It uses "json" tags instead of the default "mapstructure".
func Decode[T any](m map[string]any) (*T, error) {
	var out T
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: tagName,
		Result:  &out,
	})
	if err != nil {
		return nil, err
	}
	if err := d.Decode(m); err != nil {
		return nil, err
	}
	return &out, nil
}

// Decode is a thin wrapper around mapstructure.Decode
// that encodes a Go struct into a map[string]any.
// It uses "json" tags instead of the default "mapstructure".
func Encode[T any](v *T) (map[string]any, error) {
	out := make(map[string]any)
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: tagName,
		Result:  &out,
	})
	if err != nil {
		return nil, err
	}
	if err := d.Decode(v); err != nil {
		return nil, err
	}
	return out, nil
}
