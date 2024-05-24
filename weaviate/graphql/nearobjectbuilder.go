package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/crossref"

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
}

// WithID the id of the object
func (e *NearObjectArgumentBuilder) WithID(id string) *NearObjectArgumentBuilder {
	e.id = id
	return e
}

// WithBeacon the beacon of the object
func (e *NearObjectArgumentBuilder) WithBeacon(beacon string) *NearObjectArgumentBuilder {
	e.beacon = beacon
	return e
}

// WithCertainty that is minimally required for an object to be included in the result set
func (e *NearObjectArgumentBuilder) WithCertainty(certainty float32) *NearObjectArgumentBuilder {
	e.withCertainty = true
	e.certainty = certainty
	return e
}

// WithDistance that is minimally required for an object to be included in the result set
func (e *NearObjectArgumentBuilder) WithDistance(distance float32) *NearObjectArgumentBuilder {
	e.withDistance = true
	e.distance = distance
	return e
}

// WithTargetVectors target vector name
func (e *NearObjectArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearObjectArgumentBuilder {
	if len(targetVectors) > 0 {
		e.targetVectors = targetVectors
	}
	return e
}

// Build build the given clause
func (e *NearObjectArgumentBuilder) build() string {
	clause := []string{}
	if len(e.id) > 0 {
		clause = append(clause, fmt.Sprintf("id: \"%s\"", e.id))
	}
	if len(e.beacon) > 0 {
		clause = append(clause, fmt.Sprintf("beacon: \"%s\"", e.beacon))
	}
	if e.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", e.certainty))
	}
	if e.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", e.distance))
	}
	if len(e.targetVectors) > 0 {
		targetVectors, _ := json.Marshal(e.targetVectors)
		clause = append(clause, fmt.Sprintf("targetVectors: %s", targetVectors))
	}
	return fmt.Sprintf("nearObject:{%s}", strings.Join(clause, " "))
}

func (e *NearObjectArgumentBuilder) togrpc() *pb.NearObject {
	nearObject := &pb.NearObject{}
	id := e.id
	if len(e.beacon) > 0 {
		id = crossref.ExtractID(e.beacon)
	}
	nearObject.Id = id
	nearObject.TargetVectors = e.targetVectors
	if e.withCertainty {
		certainty := float64(e.certainty)
		nearObject.Certainty = &certainty
	}
	if e.withDistance {
		distance := float64(e.distance)
		nearObject.Distance = &distance
	}
	return nearObject
}
