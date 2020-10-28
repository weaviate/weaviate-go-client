package graphql

import (
	"encoding/json"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
)

type Explore struct {
	connection rest
	fields []paragons.ExploreFields
	concepts []string
}

func (e *Explore) WithFields(fields []paragons.ExploreFields) *Explore {
	e.fields = fields
	return e
}

func (e *Explore) WithConcepts(concepts []string) *Explore {
	e.concepts = concepts
	return e
}

func (e *Explore) build() string {
	fields := ""
	for _, field := range e.fields {
		fields += fmt.Sprintf("%v ", field)
	}

	filterClause := e.createFilterClause()

	query := fmt.Sprintf("{Explore(%v){%v}}", filterClause, fields)


	return query
}

func (e *Explore) createFilterClause() string {
	concepts, ok := json.Marshal(e.concepts)
	if ok != nil {
		return "Concepts not buildable"
	}

	clause := fmt.Sprintf("concepts: %v", string(concepts))


	return clause
}