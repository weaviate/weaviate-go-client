package internal

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type Module[T ~string] interface {
	Name() T
}

// CustomModule stores data for a module that
// isn't registered with a [Modules] registry.
type CustomModule[T ~string] interface {
	Module[T]
	Raw() map[string]any
}

// Modules is a registry of modules keyed by their [Module.Name].
//
// In addition to module registration, it supports encoding the
// module as a map[string]any and decoding a map back into the type
// registerred for the key.
//
// N.B.: before being sent to the server, each module is ultimately
// encoded into a JSON string, and is decoded from a JSON string on read.
// The reason Encode/Decode do not work with [json.RawMessage] or a plain
// []byte is that internal/gen/rest/modules.go reads module configuration
// into interface{}, which json package specializes to map[string]any.
type Modules[T ~string] struct {
	registry sync.Map // Module registry.
}

// Register adds the module to the registry. The module must be either
// a struct or a map. Modules relies on value copy after map lookup,
// so m must be the value and not a pointer, e.g. *new(struct) or make(map).
//
// This function uses reflection because all modules are registerred at
// package init time: doing the checks right away and failing quickly
// prevents potential runtime bugs. This is also not on the hotpath.
func (ms *Modules[T]) Register(m Module[T]) {
	t := reflect.TypeOf(m)
	if t == nil {
		panic("register: module is nil")
	}
	switch t.Kind() {
	case reflect.Struct, reflect.Map: // ok
	case reflect.Pointer:
		panic(fmt.Sprintf("register %s: module must be passed by value", m.Name()))
	default:
		panic(fmt.Sprintf("register %s: module must be a struct or a map", m.Name()))
	}
	ms.registry.Store(string(m.Name()), m)
}

func (ms *Modules[T]) Decode(name string, raw map[string]any) (Module[T], error) {
	m, ok := ms.registry.Load(name)
	if !ok {
		return &customModule[T]{
			name: name,
			raw:  raw,
		}, nil
	}

	dev.AssertType[Module[T]](m, "module")

	// Decode using our mapstructure wrapper.
	if err := Decode(raw, &m); err != nil {
		return nil, err
	}

	return m.(Module[T]), nil
}

func (*Modules[T]) Encode(m Module[T]) (map[string]any, error) {
	v := make(map[string]any)
	if err := Encode(m, v); err != nil {
		return nil, err
	}
	return v, nil
}

// customModule implements a simple [CustomModule].
type customModule[T ~string] struct {
	name string
	raw  map[string]any
}

func (cm *customModule[T]) Name() T             { return T(cm.name) }
func (cm *customModule[T]) Raw() map[string]any { return cm.raw }
