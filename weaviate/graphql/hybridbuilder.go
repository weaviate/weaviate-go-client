package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type HybridArgumentBuilder struct {
	query     string
	vector    []float32
	withAlpha bool
	alpha     float32
}

// WithQuery the search string
func (h *HybridArgumentBuilder) WithQuery(query string) *HybridArgumentBuilder {
	h.query = query
	return h
}

// WithVector the vector. Can be omitted
func (h *HybridArgumentBuilder) WithVector(vector []float32) *HybridArgumentBuilder {
	h.vector = vector
	return h
}

// WithAlpha the bias
func (h *HybridArgumentBuilder) WithAlpha(alpha float32) *HybridArgumentBuilder {
	h.withAlpha = true
	h.alpha = alpha
	return h
}

// Build build the given clause
func (h *HybridArgumentBuilder) build() string {
	clause := []string{}
	if h.query != "" {
		clause = append(clause, fmt.Sprintf("query: \"%s\"", h.query))
	}
	if len(h.vector) > 0 {
		vectorB, err := json.Marshal(h.vector)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal hybrid search vector: %s", err))
		}
		clause = append(clause, fmt.Sprintf("vector: %s", string(vectorB)))
	}
	if h.withAlpha {
		clause = append(clause, fmt.Sprintf("alpha: %v", h.alpha))
	}
	return fmt.Sprintf("hybrid:{%v}", strings.Join(clause, ", "))
}
