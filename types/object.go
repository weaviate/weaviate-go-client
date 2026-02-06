package types

import (
	"time"

	"github.com/google/uuid"
)

type Object[T any] struct {
	Collection string
	UUID       uuid.UUID
	Properties T
	References any

	// Vectors attached to this object. Only vectors requested in a query will be populated.
	Vectors Vectors

	//  CreatedAt is the object insertion timestamp.
	// Nil if not requested via [query.ReturnMetadata], otherwise a non-zero time.
	CreatedAt *time.Time

	// LastUpdatedAt is the object's last update timestamp.
	// Nil if not requested via [query.ReturnMetadata], otherwise a non-zero time.
	LastUpdatedAt *time.Time
}
