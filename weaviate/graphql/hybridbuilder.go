package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type FusionType string

// Ranked FusionType operator
const Ranked FusionType = "rankedFusion"

// RelativeScore FusionType operator
const RelativeScore FusionType = "relativeScoreFusion"

type HybridArgumentBuilder struct {
	query      string
	vector     []float32
	withAlpha  bool
	alpha      float32
	properties []string
	fusionType FusionType
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

// WithProperties The properties which are searched. Can be omitted.
func (h *HybridArgumentBuilder) WithProperties(properties []string) *HybridArgumentBuilder {
	h.properties = properties
	return h
}

// WithFusionType sets the fusion type to be used. Can be omitted.
func (h *HybridArgumentBuilder) WithFusionType(fusionType FusionType) *HybridArgumentBuilder {
	h.fusionType = fusionType
	return h
}

// Build build the given clause
func (h *HybridArgumentBuilder) build() string {
	clause := []string{}
	if h.query != "" {
		clause = append(clause, fmt.Sprintf("query: %q", h.query))
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

	if len(h.properties) > 0 {
		props, err := json.Marshal(h.properties)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal hybrid properties: %s", err))
		}
		clause = append(clause, fmt.Sprintf("properties: %v", string(props)))
	}

	if h.fusionType != "" {
		clause = append(clause, fmt.Sprintf("fusionType: %s", h.fusionType))
	}

	return fmt.Sprintf("hybrid:{%v}", strings.Join(clause, ", "))
}
