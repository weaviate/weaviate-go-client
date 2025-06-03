package schema

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate/entities/models"
)

// VectorAdder builder to add named vectors to a collection.
// Named vectors can only be added to collections which already define at least 1 named vector.
type VectorAdder struct {
	connection   *connection.Connection
	className    string
	addedVectors map[string]models.VectorConfig

	classGetter  *ClassGetter
	classUpdater *ClassUpdater
}

// WithClassName sets the collection name for which new vectors will be created.
func (pc *VectorAdder) WithClassName(className string) *VectorAdder {
	pc.className = className
	return pc
}

// WithVectors accepts configurations to create new named vectors.
func (pc *VectorAdder) WithVectors(vectors map[string]models.VectorConfig) *VectorAdder {
	pc.addedVectors = vectors
	return pc
}

func (va *VectorAdder) Do(ctx context.Context) error {
	class, err := va.classGetter.WithClassName(va.className).Do(ctx)
	if err != nil {
		return err
	}

	for name, vector := range va.addedVectors {
		if _, ok := class.VectorConfig[name]; ok {
			continue // named vectors are immutable
		}
		class.VectorConfig[name] = vector
	}
	return va.classUpdater.WithClass(class).Do(ctx)
}
