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
	Module internal.Module[string]
	Custom internal.CustomModule[string]
)

// Registry stores all modules registered with this package.
var Registry internal.Modules[string]

// Register a module to the package's [Registry].
func Register(m Module) {
	Registry.Register(m)
}
