package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type NearTextArgumentBuilder struct {
	concepts        []string
	withCertainty   bool
	certainty       float32
	moveTo          *MoveParameters
	moveAwayFrom    *MoveParameters
	withAutocorrect bool
	autocorrect     bool
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

// WithAutocorrect this is a setting enabling autocorrect of the concepts texts
func (e *NearTextArgumentBuilder) WithAutocorrect(autocorrect bool) *NearTextArgumentBuilder {
	e.withAutocorrect = true
	e.autocorrect = autocorrect
	return e
}

func (e *NearTextArgumentBuilder) buildMoveParam(name string, param *MoveParameters) string {
	moveToConcepts, _ := json.Marshal(param.Concepts)
	return fmt.Sprintf("%s: {concepts: %v force: %v}", name, string(moveToConcepts), param.Force)
}

// Build build the given clause
func (e *NearTextArgumentBuilder) build() string {
	clause := []string{}
	concepts, _ := json.Marshal(e.concepts)
	clause = append(clause, fmt.Sprintf("concepts: %v", string(concepts)))
	if e.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", e.certainty))
	}
	if e.moveTo != nil {
		clause = append(clause, e.buildMoveParam("moveTo", e.moveTo))
	}
	if e.moveAwayFrom != nil {
		clause = append(clause, e.buildMoveParam("moveAwayFrom", e.moveAwayFrom))
	}
	if e.withAutocorrect {
		clause = append(clause, fmt.Sprintf("autocorrect: %v", e.autocorrect))
	}
	return fmt.Sprintf("nearText:{%v}", strings.Join(clause, " "))
}
