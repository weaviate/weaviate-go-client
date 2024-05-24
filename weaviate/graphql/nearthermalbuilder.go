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
	return builder.build()
}

func (b *NearThermalArgumentBuilder) togrpc() *pb.NearThermalSearch {
	builder := &nearMediaArgumentBuilder{
		data:       b.thermal,
		dataReader: b.thermalReader,
	}
	nearThermal := &pb.NearThermalSearch{
		Thermal:       builder.getContent(),
		TargetVectors: b.targetVectors,
	}
	if b.hasCertainty {
		certainty := float64(b.certainty)
		nearThermal.Certainty = &certainty
	}
	if b.hasDistance {
		distance := float64(b.distance)
		nearThermal.Distance = &distance
	}
	return nearThermal
}
