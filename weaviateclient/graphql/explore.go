package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"github.com/semi-technologies/weaviate/entities/models"
)

// Explore query builder
type Explore struct {
	connection rest
	fields []paragons.ExploreFields
	concepts []string

	withLimit bool
	limit int
	withCertainty bool
	certainty float32

	moveTo *paragons.MoveParameters
	moveAwayFrom *paragons.MoveParameters
}

// WithFields that should be included in the result set
func (e *Explore) WithFields(fields []paragons.ExploreFields) *Explore {
	e.fields = fields
	return e
}

// WithConcepts the result is based on
func (e *Explore) WithConcepts(concepts []string) *Explore {
	e.concepts = concepts
	return e
}

// WithLimit of objects in result set
func (e *Explore) WithLimit(limit int) *Explore {
	e.withLimit = true
	e.limit = limit
	return e
}

// WithCertainty that is minimally required for an object to be included in the result set
func (e *Explore) WithCertainty(certainty float32) *Explore {
	e.withCertainty = true
	e.certainty = certainty
	return e
}

// WithMoveTo specific concept
func (e *Explore) WithMoveTo(parameters *paragons.MoveParameters) *Explore {
	e.moveTo = parameters
	return e
}

// WithMoveAwayFrom specific concept
func (e *Explore) WithMoveAwayFrom(parameters *paragons.MoveParameters) *Explore {
	e.moveAwayFrom = parameters
	return e
}

// build query
func (e *Explore) build() string {
	fields := ""
	for _, field := range e.fields {
		fields += fmt.Sprintf("%v ", field)
	}

	filterClause := e.createFilterClause()

	query := fmt.Sprintf("{Explore(%v){%v}}", filterClause, fields)


	return query
}

// Do execute explore search
func (e *Explore) Do(ctx context.Context) (*models.GraphQLResponse, error) {
	return runGraphQLQuery(ctx, e.connection, e.build())
}

func (e *Explore) createFilterClause() string {
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

	return clause
}