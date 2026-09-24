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
// The caller can control how the map is decoded into T
// by implementing [MapDecoder] for T.
func Decode[T any](m map[string]any, dest *T) error {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:    tagName,
		Result:     dest,
		DecodeHook: decodeHook(reflect.TypeFor[T](), dest),
	})
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
// The caller can control how v is incoded into a map
// by implementing [MapEncoder] for the type of v.
func Encode(v any, dest map[string]any) error {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:    tagName,
		Result:     &dest,
		DecodeHook: encodeHook(reflect.TypeOf(v)),
	})
	if err != nil {
		return err
	}
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}

type MapDecoder interface {
	DecodeMap(map[string]any) error
}

type MapEncoder interface {
	EncodeMap() (map[string]any, error)
}

var mapStringAnyType = reflect.TypeFor[map[string]any]()

// encodeHook wraps [MapEncoder] in [mapstructure.DecodeHookFuncType].
// so that EncodeMap is only called with an appropriate source type.
func encodeHook(t reflect.Type) mapstructure.DecodeHookFuncType {
	return func(from, to reflect.Type, data any) (any, error) {
		if from != t {
			return data, nil
		}

		me, ok := data.(MapEncoder)
		if !ok {
			return data, nil
		}

		// Because we control encodeHook, we can expect that it
		// will only be used to encode a struct into a map[string]any.
		dev.Assert(
			to == mapStringAnyType,
			"EncodeHook must be used with %s, not %s",
			mapStringAnyType.Name(), to.Name(),
		)

		return me.EncodeMap()
	}
}

// decodeHook wraps [MapDecoder] in [mapstructure.DecodeHookFuncType],
// so that DecodeMap is only called with map[string]any data.
func decodeHook(t reflect.Type, dest any) mapstructure.DecodeHookFuncType {
	md, ok := dest.(MapDecoder)
	if !ok {
		return func(_, _ reflect.Type, data any) (any, error) {
			return data, nil
		}
	}
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

		if err := md.DecodeMap(m); err != nil {
			return data, err
		}
		return md, nil
	}
}
