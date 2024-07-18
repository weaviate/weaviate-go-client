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
	query         string
	vector        []float32
	withAlpha     bool
	alpha         float32
	properties    []string
	fusionType    FusionType
	targetVectors []string
	targets       *MultiTargetArgumentBuilder
	searches      *HybridSearchesArgumentBuilder
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

// WithTargetVectors sets the target vectors to be used with hybrid query.
func (h *HybridArgumentBuilder) WithTargetVectors(targetVectors ...string) *HybridArgumentBuilder {
	h.targetVectors = targetVectors
	return h
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (h *HybridArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *HybridArgumentBuilder {
	h.targets = targets
	return h
}

// WithSearches sets the searches to be used with hybrid.
func (h *HybridArgumentBuilder) WithSearches(searches *HybridSearchesArgumentBuilder) *HybridArgumentBuilder {
	h.searches = searches
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

	if h.targets != nil {
		clause = append(clause, fmt.Sprintf("targets:{%s}", h.targets.build()))
	}

	if len(h.targetVectors) > 0 && h.targets == nil {
		targetVectors, _ := json.Marshal(h.targetVectors)
		clause = append(clause, fmt.Sprintf("targetVectors: %s", targetVectors))
	}

	if h.searches != nil {
		clause = append(clause, fmt.Sprintf("searches:{%s}", h.searches.build()))
	}

	return fmt.Sprintf("hybrid:{%v}", strings.Join(clause, ", "))
}

type HybridSearchesArgumentBuilder struct {
	nearVector *NearVectorArgumentBuilder
	nearText   *NearTextArgumentBuilder
}

func (s *HybridSearchesArgumentBuilder) WithNearVector(nearVector *NearVectorArgumentBuilder) *HybridSearchesArgumentBuilder {
	s.nearVector = nearVector
	return s
}

func (s *HybridSearchesArgumentBuilder) WithNearText(nearText *NearTextArgumentBuilder) *HybridSearchesArgumentBuilder {
	s.nearText = nearText
	return s
}

func (h *HybridSearchesArgumentBuilder) build() string {
	searches := []string{}
	if h.nearText != nil {
		searches = append(searches, h.nearText.build())
	}
	if h.nearVector != nil {
		searches = append(searches, h.nearVector.build())
	}
	return strings.Join(searches, " ")
}
