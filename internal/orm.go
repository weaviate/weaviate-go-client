package internal

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

// tagName is the tag mapstructure will use to match struct fields to keys in the properties map.
const tagName = "json"

// Decode is a thin wrapper around mapstructure.Decode
// that decodes map[string]any into a Go struct.
// It uses "json" tags instead of the default "mapstructure".
//
// Passing an optional [DecodeHookFunc] lets the caller control
// how the map is decoded into T. Leave as nil to use default logic.
func Decode[T any](m map[string]any, dest *T, h DecodeHookFunc) error {
	conf := &mapstructure.DecoderConfig{
		TagName: tagName,
		Result:  dest,
	}
	if h != nil {
		// We cannot assign the output directly to DecodeHook, because that
		// turns nil into a typed nil, that mapstructure cannot check with ==.
		conf.DecodeHook = h.hook(reflect.TypeFor[T]())
	}
	d, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return err
	}
	if err := d.Decode(m); err != nil {
		return err
	}
	return nil
}

// Encode is a thin wrapper around mapstructure.Decode
// that encodes a Go struct into a map[string]any.
// It uses "json" tags instead of the default "mapstructure".
//
// Passing an optional [EncodeHookFunc] lets the caller control
// how the value is encoded into a map. Leave as nil to use default logic.
func Encode(v any, dest map[string]any, h EncodeHookFunc) error {
	conf := &mapstructure.DecoderConfig{
		TagName: tagName,
		Result:  &dest,
	}
	if h != nil {
		// We cannot assign the output directly to DecodeHook, because that
		// turns nil into a typed nil, that mapstructure cannot check with ==.
		conf.DecodeHook = h.hook(reflect.TypeOf(v))
	}
	d, err := mapstructure.NewDecoder(conf)
	if err != nil {
		return err
	}
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}

// Hook enable custom encoding and decoding.
type Hook struct {
	Encode EncodeHookFunc
	Decode DecodeHookFunc
}

type (
	EncodeHookFunc func(from any) (map[string]any, error)
	DecodeHookFunc func(from map[string]any) (any, error)
)

var mapStringAnyType = reflect.TypeFor[map[string]any]()

// hook wraps [EncodeHookFunc] in [mapstructure.DecodeHookFuncType].
// so that f is only called with an appropriate source type.
func (f EncodeHookFunc) hook(t reflect.Type) mapstructure.DecodeHookFuncType {
	return func(from, to reflect.Type, data any) (any, error) {
		if from != t {
			return data, nil
		}

		// Because we control encodeHook, we can expect that it
		// will only be used to encode a struct into a map[string]any.
		dev.Assert(
			to == mapStringAnyType,
			"EncodeHook must be used with %s, not %s",
			mapStringAnyType.Name(), to.Name(),
		)

		return f(data)
	}
}

// hook wraps [DecodeHookFunc] in [mapstructure.DecodeHookFuncType],
// so that f is only called with map[string]any data.
func (f DecodeHookFunc) hook(t reflect.Type) mapstructure.DecodeHookFuncType {
	return func(from, to reflect.Type, data any) (any, error) {
		if to != t {
			return data, nil
		}

		// Because we control decodeHook, we can expect that it
		// will only be used to decode a map[string]any.
		dev.Assert(
			from == mapStringAnyType,
			"DecodeHook must be used with %s, not %s",
			mapStringAnyType.Name(), from.Name(),
		)

		m, ok := data.(map[string]any)
		if !ok {
			dev.Unreachable() // This would mean a bug in mapstructure package.
			return data, fmt.Errorf("expected data (%T) to be map[string]any", data)
		}
		return f(m)
	}
}
