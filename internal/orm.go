package internal

import "github.com/go-viper/mapstructure/v2"

// tagName is the tag mapstructure will use to match struct fields to keys in the properties map.
const tagName = "json"

// Decode is a thin wrapper around mapstructure.Decode
// that uses "json" tags instead of the default "mapstructure".
func Decode[P any](m map[string]any) (*P, error) {
	var out P
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:    tagName,
		ZeroFields: true,
		Result:     &out,
	})
	if err != nil {
		return nil, err
	}
	if err := d.Decode(m); err != nil {
		return nil, err
	}
	return &out, nil
}
