package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
	"github.com/weaviate/weaviate/usecases/byteops"
)

type NearVectorArgumentBuilder struct {
	vector        []float32
	withCertainty bool
	certainty     float32
	withDistance  bool
	distance      float32
	targetVectors []string
}

// WithVector sets the search vector to be used in query
func (b *NearVectorArgumentBuilder) WithVector(vector []float32) *NearVectorArgumentBuilder {
	b.vector = vector
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

// Build build the given clause
func (b *NearVectorArgumentBuilder) build() string {
	clause := []string{}
	if b.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", b.certainty))
	}
	if b.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", b.distance))
	}
	if len(b.vector) != 0 {
		vectorB, err := json.Marshal(b.vector)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal nearVector search vector: %s", err))
		}
		clause = append(clause, fmt.Sprintf("vector: %s", string(vectorB)))
	}
	if len(b.targetVectors) > 0 {
		targetVectors, _ := json.Marshal(b.targetVectors)
		clause = append(clause, fmt.Sprintf("targetVectors: %s", targetVectors))
	}
	return fmt.Sprintf("nearVector:{%v}", strings.Join(clause, " "))
}

func (b *NearVectorArgumentBuilder) togrpc() *pb.NearVector {
	nearVector := &pb.NearVector{
		VectorBytes:   byteops.Float32ToByteVector(b.vector),
		TargetVectors: b.targetVectors,
	}
	if b.withCertainty {
		certainty := float64(b.certainty)
		nearVector.Certainty = &certainty
	}
	if b.withDistance {
		distance := float64(b.distance)
		nearVector.Distance = &distance
	}
	return nearVector
}
