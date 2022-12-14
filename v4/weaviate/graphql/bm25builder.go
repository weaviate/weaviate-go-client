package graphql

import (
	"fmt"
	"strings"

	"encoding/json"
)

type BM25ArgumentBuilder struct {
	Properties []string
	Query      string
}

// WithQuery the search string
func (e *BM25ArgumentBuilder) WithQuery(query string) *BM25ArgumentBuilder {
	e.Query = query
	return e
}

// WithProperties the properties to search.  Leave blank for all
func (e *BM25ArgumentBuilder) WithProperties(properties ...string) *BM25ArgumentBuilder {
	e.Properties = properties
	return e
}

// Build build the given clause
func (e *BM25ArgumentBuilder) build() string {
	clause := []string{}

	clause = append(clause, fmt.Sprintf("query: \"%s\"", e.Query))
	if len(e.Properties) > 0 {
		propStr, err := json.Marshal(e.Properties)
		if err != nil {
			panic(err)
		}
		clause = append(clause, fmt.Sprintf("properties: %v", string(propStr)))
	}
	return fmt.Sprintf("bm25:{%v}", strings.Join(clause, ", "))
}
