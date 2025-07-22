package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/crossref"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type NearObjectArgumentBuilder struct {
	id            string
	beacon        string
	withCertainty bool
	certainty     float32
	withDistance  bool
	distance      float32
	targetVectors []string
	targets       *MultiTargetArgumentBuilder
}

// WithID the id of the object
func (b *NearObjectArgumentBuilder) WithID(id string) *NearObjectArgumentBuilder {
	b.id = id
	return b
}

// WithBeacon the beacon of the object
func (b *NearObjectArgumentBuilder) WithBeacon(beacon string) *NearObjectArgumentBuilder {
	b.beacon = beacon
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearObjectArgumentBuilder) WithCertainty(certainty float32) *NearObjectArgumentBuilder {
	b.withCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearObjectArgumentBuilder) WithDistance(distance float32) *NearObjectArgumentBuilder {
	b.withDistance = true
	b.distance = distance
	return b
}

// WithTargetVectors target vector name
func (b *NearObjectArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearObjectArgumentBuilder {
	if len(targetVectors) > 0 {
		b.targetVectors = targetVectors
	}
	return b
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (b *NearObjectArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *NearObjectArgumentBuilder {
	b.targets = targets
	return b
}

// Build build the given clause
func (b *NearObjectArgumentBuilder) build() string {
	clause := []string{}
	if len(b.id) > 0 {
		clause = append(clause, fmt.Sprintf("id: \"%s\"", b.id))
	}
	if len(b.beacon) > 0 {
		clause = append(clause, fmt.Sprintf("beacon: \"%s\"", b.beacon))
	}
	if b.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", b.certainty))
	}
	if b.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", b.distance))
	}
	if b.targets != nil {
		clause = append(clause, fmt.Sprintf("targets:{%s}", b.targets.build()))
	}
	if len(b.targetVectors) > 0 && b.targets == nil {
		targetVectors, _ := json.Marshal(b.targetVectors)
		clause = append(clause, fmt.Sprintf("targetVectors: %s", targetVectors))
	}
	return fmt.Sprintf("nearObject:{%s}", strings.Join(clause, " "))
}

func (b *NearObjectArgumentBuilder) togrpc() *pb.NearObject {
	nearObject := &pb.NearObject{}
	id := b.id
	if len(b.beacon) > 0 {
		id = crossref.ExtractID(b.beacon)
	}
	nearObject.Id = id
	if b.withCertainty {
		certainty := float64(b.certainty)
		nearObject.Certainty = &certainty
	}
	if b.withDistance {
		distance := float64(b.distance)
		nearObject.Distance = &distance
	}
	if b.targets != nil {
		nearObject.Targets = b.targets.togrpc()
	}
	if len(b.targetVectors) > 0 && b.targets == nil {
		nearObject.Targets = &pb.Targets{TargetVectors: b.targetVectors}
	}
	return nearObject
}
