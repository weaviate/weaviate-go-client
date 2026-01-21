package types

import (
	"time"

	"github.com/google/uuid"
)

type Object[T any] struct {
	UUID          uuid.UUID
	Properties    T
	References    any
	Vectors       Vectors
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}
