// modules provides a registry for third-party modules integrated with Weaviate.
//
// To register a new module:
//
//	import "github.com/weaviate/weaviate-go-client/v6/modules"
//
//	modules.Register(*new(CustomModule))
//
// A module must be a struct and implement [Module] with a value receiver.
// The function panics if [Register] is passed a pointer and not a value,
// as the registry expects to copy the value before decoding into it.
package modules

import (
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

type (
	Module internal.Module
	Custom internal.CustomModule
)

var registry internal.Modules

func Register(m Module) {
	registry.Register(m)
}

// Encode creates a map of all public fields of the module struct.
func Encode(m Module) (map[string]any, error) {
	return registry.Encode(m)
}

// Decode creates a new instance of the module struct.
// The type depends on the type registered for the name.
func Decode(name string, raw map[string]any) (Module, error) {
	return registry.Decode(name, raw)
}
