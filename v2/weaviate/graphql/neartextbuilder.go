package graphql

import (
	"encoding/json"
	"fmt"
)

type NearTextArgumentBuilder struct {
	concepts      []string
	withCertainty bool
	certainty     float32
	moveTo        *MoveParameters
	moveAwayFrom  *MoveParameters
	withLimit     bool
	limit         int
}

// WithConcepts the result is based on
func (e *NearTextArgumentBuilder) WithConcepts(concepts []string) *NearTextArgumentBuilder {
	e.concepts = concepts
	return e
}

// WithCertainty that is minimally required for an object to be included in the result set
func (e *NearTextArgumentBuilder) WithCertainty(certainty float32) *NearTextArgumentBuilder {
	e.withCertainty = true
	e.certainty = certainty
	return e
}

// WithMoveTo specific concept
func (e *NearTextArgumentBuilder) WithMoveTo(parameters *MoveParameters) *NearTextArgumentBuilder {
	e.moveTo = parameters
	return e
}

// WithMoveAwayFrom specific concept
func (e *NearTextArgumentBuilder) WithMoveAwayFrom(parameters *MoveParameters) *NearTextArgumentBuilder {
	e.moveAwayFrom = parameters
	return e
}

// WithLimit of objects in result set
func (e *NearTextArgumentBuilder) WithLimit(limit int) *NearTextArgumentBuilder {
	e.withLimit = true
	e.limit = limit
	return e
}

// Build build the given clause
func (e *NearTextArgumentBuilder) build() string {
	concepts, _ := json.Marshal(e.concepts)
	clause := fmt.Sprintf("concepts: %v ", string(concepts))

	if e.withLimit {
		clause += fmt.Sprintf("limit: %v ", e.limit)
	}
	if e.withCertainty {
		clause += fmt.Sprintf("certainty: %v ", e.certainty)
	}
	if e.moveTo != nil {
		moveToConcepts, _ := json.Marshal(e.moveTo.Concepts)
		clause += fmt.Sprintf("moveTo: {concepts: %v force: %v} ", string(moveToConcepts), e.moveTo.Force)
	}
	if e.moveAwayFrom != nil {
		moveAwayFromConcepts, _ := json.Marshal(e.moveAwayFrom.Concepts)
		clause += fmt.Sprintf("moveAwayFrom: {concepts: %v force: %v} ", string(moveAwayFromConcepts), e.moveAwayFrom.Force)
	}

	return fmt.Sprintf("nearText:{%v} ", clause)
}
