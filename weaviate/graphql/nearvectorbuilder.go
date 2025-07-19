package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc/common"
	"github.com/weaviate/weaviate/entities/models"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type NearVectorArgumentBuilder struct {
	vector           models.Vector
	vectorsPerTarget map[string][]models.Vector
	withCertainty    bool
	certainty        float32
	withDistance     bool
	distance         float32
	targetVectors    []string
	targets          *MultiTargetArgumentBuilder
}

// WithVector sets the search vector to be used in query
func (b *NearVectorArgumentBuilder) WithVector(vector models.Vector) *NearVectorArgumentBuilder {
	b.vector = vector
	return b
}

// WithVectorPerTarget sets the search vector per target to be used in a multi target search query. This builder method takes
// precedence over WithVector. So if WithVectorPerTarget is used, WithVector will be ignored.
func (b *NearVectorArgumentBuilder) WithVectorPerTarget(vectorPerTarget map[string]models.Vector) *NearVectorArgumentBuilder {
	if len(vectorPerTarget) > 0 {
		vectorPerTargetTmp := make(map[string][]models.Vector)
		for k, v := range vectorPerTarget {
			vectorPerTargetTmp[k] = []models.Vector{v}
		}
		b.vectorsPerTarget = vectorPerTargetTmp
	}
	return b
}

// WithVectorsPerTarget sets the search vector per target to be used in a multi target search query. This builder method takes
// precedence over WithVector and WithVectorPerTarget. So if WithVectorsPerTarget is used, WithVector and WithVectorPerTarget will be ignored.
func (b *NearVectorArgumentBuilder) WithVectorsPerTarget(vectorPerTarget map[string][]models.Vector) *NearVectorArgumentBuilder {
	if len(vectorPerTarget) > 0 {
		b.vectorsPerTarget = vectorPerTarget
	}
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearVectorArgumentBuilder) WithCertainty(certainty float32) *NearVectorArgumentBuilder {
	b.withCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearVectorArgumentBuilder) WithDistance(distance float32) *NearVectorArgumentBuilder {
	b.withDistance = true
	b.distance = distance
	return b
}

// WithTargetVectors target vector name
func (b *NearVectorArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearVectorArgumentBuilder {
	if len(targetVectors) > 0 {
		b.targetVectors = targetVectors
	}
	return b
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder method takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (b *NearVectorArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *NearVectorArgumentBuilder {
	b.targets = targets
	return b
}

// Build build the given clause
func (b *NearVectorArgumentBuilder) build() string {
	clause := []string{}
	if b.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", b.certainty))
	}
	if b.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", b.distance))
	}

	if len(b.vectorsPerTarget) > 0 {
		vectorPerTarget := make([]string, 0, len(b.vectorsPerTarget))
		for k, v := range b.vectorsPerTarget {
			vBytes, err := json.Marshal(v)
			if err != nil {
				panic(fmt.Sprintf("could not marshal vector: %v", err))
			}
			vectorPerTarget = append(vectorPerTarget, fmt.Sprintf("%s: %v", k, string(vBytes)))
		}
		clause = append(clause, fmt.Sprintf("vectorPerTarget: {%s}", strings.Join(vectorPerTarget, ",")))
	}
	if !b.isVectorEmpty(b.vector) && len(b.vectorsPerTarget) == 0 {
		vectorB, err := json.Marshal(b.vector)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal nearVector search vector: %s", err))
		}
		clause = append(clause, fmt.Sprintf("vector: %s", string(vectorB)))
	}
	if b.targets != nil {
		clause = append(clause, fmt.Sprintf("targets: {%s}", b.targets.build()))
	}

	targetVectors := b.prepareTargetVectors(b.targetVectors)
	if len(targetVectors) > 0 {
		targetVectors, _ := json.Marshal(targetVectors)
		clause = append(clause, fmt.Sprintf("targetVectors: %s", targetVectors))
	}
	return fmt.Sprintf("nearVector:{%v}", strings.Join(clause, " "))
}

func (b *NearVectorArgumentBuilder) isVectorEmpty(vector models.Vector) bool {
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

// prepareTargetVectors adds appends the target name for each target vector associated with it.
// Example:
//
//	// For target vectors:
//	WithTargetVectors("v1", "v2").
//	WithVectorProTarget(map[string][][]float32{"v1": {{1,2,3}, {4,5,6}}})
//	// Outputs:
//	[]string{"v1", "v1", "v2"}
//
// The server requires that the target names be repeated for each target vector,
// and passing them once only is a mistake that the users can easily make.
// This way, the client provides some safeguard.
//
// Note, too, that in case the user fails to pass a value in TargetVectors,
// it will not be added to the query.
func (b NearVectorArgumentBuilder) prepareTargetVectors(targets []string) (out []string) {
	for _, target := range targets {
		if vectors, ok := b.vectorsPerTarget[target]; ok {
			for range vectors {
				out = append(out, target)
			}
		} else {
			out = append(out, target)
		}
	}
	return
}

func (b *NearVectorArgumentBuilder) togrpc() *pb.NearVector {
	nearVector := &pb.NearVector{}
	if !b.isVectorEmpty(b.vector) && len(b.vectorsPerTarget) == 0 {
		nearVector.Vectors = []*pb.Vectors{common.GetVector("", b.vector)}
	}
	if b.withCertainty {
		certainty := float64(b.certainty)
		nearVector.Certainty = &certainty
	}
	if b.withDistance {
		distance := float64(b.distance)
		nearVector.Distance = &distance
	}
	if len(b.vectorsPerTarget) > 0 {
		var targetVectors []string
		var vectorForTargets []*pb.VectorForTarget
		for targetVector, vecs := range b.vectorsPerTarget {
			for _, v := range vecs {
				vectorForTargets = append(vectorForTargets, &pb.VectorForTarget{
					Name:    targetVector,
					Vectors: []*pb.Vectors{common.GetVector("", v)},
				})
				targetVectors = append(targetVectors, targetVector)
			}
		}
		if b.targets != nil {
			nearVector.Targets = b.targets.togrpc()
			nearVector.Targets.TargetVectors = targetVectors
		}
		nearVector.VectorForTargets = vectorForTargets
	}
	if len(b.targetVectors) > 0 && b.targets == nil {
		nearVector.Targets = &pb.Targets{TargetVectors: b.targetVectors}
	}
	return nearVector
}
