package internal

import (
	"github.com/go-viper/mapstructure/v2"
)

// tagName is the tag mapstructure will use to match struct fields to keys in the properties map.
const tagName = "json"

// Decode is a thin wrapper around mapstructure.Decode
// that decodes map[string]any into a Go struct.
// It uses "json" tags instead of the default "mapstructure".
func Decode[T any](m map[string]any, dest *T) error {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: tagName,
		Result:  dest,
	})
	if err != nil {
		return err
	}
	if err := d.Decode(m); err != nil {
		return err
	}
	return nil
}

// Decode is a thin wrapper around mapstructure.Decode
// that encodes a Go struct into a map[string]any.
// It uses "json" tags instead of the default "mapstructure".
func Encode[T any](v *T, dest map[string]any) error {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: tagName,
		Result:  &dest,
	})
	if err != nil {
		return err
	}
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}
