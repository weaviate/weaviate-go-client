package types

import "github.com/google/uuid"

type Map map[string]any

type Properties interface {
	Map | any
}

type Object[P Properties] struct {
	UUID               uuid.UUID
	Properties         P
	Vectors            Vectors
	CreationTimeUnix   int64
	LastUpdateTimeUnix int64
}
