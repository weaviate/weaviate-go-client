package graphql

import (
	"fmt"
	"strings"
)

type HybridArgumentBuilder struct {
	Query  string
	Vector []float32
	Alpha  float32
}

// WithQuery the search string
func (e *HybridArgumentBuilder) WithQuery(query string) *HybridArgumentBuilder {
	e.Query = query
	return e
}

// WithVector the vector.  Can be omitted
func (e *HybridArgumentBuilder) WithVector(vector ...float32) *HybridArgumentBuilder {
	e.Vector = vector
	return e
}

// WithAlpha the bias
func (e *HybridArgumentBuilder) WithAlpha(alpha float32) *HybridArgumentBuilder {
	e.Alpha = alpha
	return e
}

// Build build the given clause
func (e *HybridArgumentBuilder) build() string {
	clause := []string{}

	clause = append(clause, fmt.Sprintf("query: \"%s\"", e.Query))
	clause = append(clause, "vector: "+strings.Replace(fmt.Sprintf("%v", e.Vector), " ", ", ", -1))
	clause = append(clause, fmt.Sprintf("alpha: %v", e.Alpha))
	return fmt.Sprintf("hybrid:{%v}", strings.Join(clause, ", "))
}
