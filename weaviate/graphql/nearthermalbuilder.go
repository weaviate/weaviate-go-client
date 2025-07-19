package graphql

import (
	"io"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type NearThermalArgumentBuilder struct {
	thermal       string
	thermalReader io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
	targetVectors []string
	targets       *MultiTargetArgumentBuilder
}

// WithThermal base64 encoded thermal
func (b *NearThermalArgumentBuilder) WithThermal(thermal string) *NearThermalArgumentBuilder {
	b.thermal = thermal
	return b
}

// WithReader the thermal file
func (b *NearThermalArgumentBuilder) WithReader(thermalReader io.Reader) *NearThermalArgumentBuilder {
	b.thermalReader = thermalReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearThermalArgumentBuilder) WithCertainty(certainty float32) *NearThermalArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearThermalArgumentBuilder) WithDistance(distance float32) *NearThermalArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

// WithTargetVectors target vector name
func (b *NearThermalArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearThermalArgumentBuilder {
	b.targetVectors = targetVectors
	return b
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (h *NearThermalArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *NearThermalArgumentBuilder {
	h.targets = targets
	return h
}

// Build build the given clause
func (b *NearThermalArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:     "nearThermal",
		mediaField:    "thermal",
		data:          b.thermal,
		dataReader:    b.thermalReader,
		targetVectors: b.targetVectors,
	}
	if b.hasCertainty {
		builder.withCertainty(b.certainty)
	}
	if b.hasDistance {
		builder.withDistance(b.distance)
	}
	if len(b.targetVectors) > 0 {
		builder.withTargetVectors(b.targetVectors...)
	}
	builder.withTargets(b.targets)
	return builder.build()
}

func (b *NearThermalArgumentBuilder) togrpc() *pb.NearThermalSearch {
	builder := &nearMediaArgumentBuilder{
		data:       b.thermal,
		dataReader: b.thermalReader,
	}
	nearThermal := &pb.NearThermalSearch{
		Thermal: builder.getContent(),
	}
	if b.hasCertainty {
		certainty := float64(b.certainty)
		nearThermal.Certainty = &certainty
	}
	if b.hasDistance {
		distance := float64(b.distance)
		nearThermal.Distance = &distance
	}
	if b.targets != nil {
		nearThermal.Targets = b.targets.togrpc()
	}
	if len(b.targetVectors) > 0 && b.targets == nil {
		nearThermal.Targets = &pb.Targets{TargetVectors: b.targetVectors}
	}
	return nearThermal
}
