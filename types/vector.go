package types

import "github.com/weaviate/weaviate-go-client/v6/internal/api"

type (
	Vector api.Vector

	// Vectors is a map of named vectors.
	// An empty string key is an alias for "default" vector.
	Vectors map[string]Vector
)
