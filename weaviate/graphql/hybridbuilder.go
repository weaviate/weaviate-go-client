package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc/common"
	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type FusionType string

// Ranked FusionType operator
const Ranked FusionType = "rankedFusion"

// RelativeScore FusionType operator
const RelativeScore FusionType = "relativeScoreFusion"

type HybridArgumentBuilder struct {
	query                 string
	vector                models.Vector
	withAlpha             bool
	alpha                 float32
	withMaxVectorDistance bool
	maxVectorDistance     float32
	properties            []string
	fusionType            FusionType
	targetVectors         []string
	targets               *MultiTargetArgumentBuilder
	searches              *HybridSearchesArgumentBuilder
	bm25SearchOperator    *BM25SearchOperatorBuilder
}

// WithQuery the search string
func (h *HybridArgumentBuilder) WithQuery(query string) *HybridArgumentBuilder {
	h.query = query
	return h
}

func (h *HybridArgumentBuilder) WithBM25SearchOperator(searchOperator BM25SearchOperatorBuilder) *HybridArgumentBuilder {
	h.bm25SearchOperator = &searchOperator
	return h
}

// WithVector the vector. Can be omitted
func (h *HybridArgumentBuilder) WithVector(vector models.Vector) *HybridArgumentBuilder {
	h.vector = vector
	return h
}

// WithAlpha the bias
func (h *HybridArgumentBuilder) WithAlpha(alpha float32) *HybridArgumentBuilder {
	h.withAlpha = true
	h.alpha = alpha
	return h
}

// WithMaxVectorDistance is the equivalent of 'distance' threshold in vector search.
func (s *HybridArgumentBuilder) WithMaxVectorDistance(d float32) *HybridArgumentBuilder {
	s.withMaxVectorDistance = true
	s.maxVectorDistance = d
	return s
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
	if !h.isVectorEmpty(h.vector) {
		vectorB, err := json.Marshal(h.vector)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal hybrid search vector: %s", err))
		}
		clause = append(clause, fmt.Sprintf("vector: %s", string(vectorB)))
	}
	if h.withAlpha {
		clause = append(clause, fmt.Sprintf("alpha: %v", h.alpha))
	}
	if h.withMaxVectorDistance {
		clause = append(clause, fmt.Sprintf("maxVectorDistance: %v", h.maxVectorDistance))
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

	if h.bm25SearchOperator != nil {
		clause = append(clause, fmt.Sprintf("bm25SearchOperator:%s", h.bm25SearchOperator.build()))
	}

	return fmt.Sprintf("hybrid:{%v}", strings.Join(clause, ", "))
}

func (h *HybridArgumentBuilder) isVectorEmpty(vector models.Vector) bool {
	switch v := vector.(type) {
	case []float32:
		return len(v) == 0
	case [][]float32:
		return len(v) == 0
	case models.C11yVector:
		return len(v) == 0
	default:
		return true
	}
}

func (h *HybridArgumentBuilder) togrpc() *pb.Hybrid {
	hybrid := &pb.Hybrid{
		Query: h.query,
	}
	if len(h.properties) > 0 {
		hybrid.Properties = h.properties
	}
	if h.withAlpha {
		hybrid.Alpha = h.alpha
	}
	if !h.isVectorEmpty(h.vector) {
		hybrid.Vectors = []*pb.Vectors{common.GetVector("", h.vector)}
	}
	if h.targets != nil {
		hybrid.Targets = h.targets.togrpc()
	}
	if len(h.targetVectors) > 0 && h.targets == nil {
		hybrid.Targets = &pb.Targets{TargetVectors: h.targetVectors}
	}
	switch h.fusionType {
	case Ranked:
		hybrid.FusionType = pb.Hybrid_FUSION_TYPE_RANKED
	case RelativeScore:
		hybrid.FusionType = pb.Hybrid_FUSION_TYPE_RELATIVE_SCORE
	default:
		hybrid.FusionType = pb.Hybrid_FUSION_TYPE_UNSPECIFIED
	}
	if h.searches != nil {
		if h.searches.nearText != nil {
			hybrid.NearText = h.searches.nearText.togrpc()
		}
		if h.searches.nearVector != nil {
			hybrid.NearVector = h.searches.nearVector.togrpc()
		}
	}
	if h.bm25SearchOperator != nil {
		hybrid.Bm25SearchOperator = h.bm25SearchOperator.togrpc()
	}
	return hybrid
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

func (s *HybridSearchesArgumentBuilder) build() string {
	var searches []string
	if s.nearText != nil {
		searches = append(searches, s.nearText.build())
	}
	if s.nearVector != nil {
		searches = append(searches, s.nearVector.build())
	}
	return strings.Join(searches, " ")
}
