package types

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

type Vector api.Vector

// Vector converts itself to an api.Vector.
func (v Vector) Vector() api.Vector { return api.Vector(v) }

// Vectors returns a single target vector.
func (v Vector) Vectors() []api.TargetVector { return []api.TargetVector{{Vector: v.Vector()}} }

// Vectors is a map of named vectors.
// An empty string key is an alias for "default" vector.
type Vectors map[string]Vector
