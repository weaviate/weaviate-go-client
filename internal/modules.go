package internal

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type Module interface {
	Name() string
}

// CustomModule stores data for a module that
// isn't registered with a [Modules] registry.
type CustomModule interface {
	Module
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
type Modules struct {
	registry sync.Map // Module registry.
}

// Register adds the module to the registry. The module must be either
// a struct or a map. Modules relies on value copy after map lookup,
// so m must be the value and not a pointer, e.g. *new(struct) or make(map).
//
// This function uses reflection because all modules are registerred at
// package init time: doing the checks right away and failing quickly
// prevents potential runtime bugs. This is also not on the hotpath.
func (ms *Modules) Register(m Module) {
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
	ms.registry.Store(m.Name(), m)
}

func (ms *Modules) Decode(name string, raw map[string]any) (Module, error) {
	m, ok := ms.registry.Load(name)
	if !ok {
		return &customModule{
			name: name,
			raw:  raw,
		}, nil
	}

	dev.AssertType[Module](m, "module")

	// Delegate to custom Decode hook if the module has one.
	if d, ok := m.(interface {
		Decode(map[string]any) (Module, error)
	}); ok {
		return d.Decode(raw)
	}

	// Decode using our mapstructure wrapper.
	if err := Decode(raw, &m); err != nil {
		return nil, err
	}

	return m.(Module), nil
}

func (*Modules) Encode(m Module) (map[string]any, error) {
	v := make(map[string]any)
	if err := Encode(m, v); err != nil {
		return nil, err
	}
	return v, nil
}

// customModule implements a simple [CustomModule].
type customModule struct {
	name string
	raw  map[string]any
}

func (cm *customModule) Name() string        { return cm.name }
func (cm *customModule) Raw() map[string]any { return cm.raw }
