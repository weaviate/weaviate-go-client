package graphql

import (
	"io"
)

type NearThermalArgumentBuilder struct {
	thermal       string
	thermalReader io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
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

// Build build the given clause
func (b *NearThermalArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:  "nearThermal",
		mediaField: "thermal",
		data:       b.thermal,
		dataReader: b.thermalReader,
	}
	if b.hasCertainty {
		builder.withCertainty(b.certainty)
	}
	if b.hasDistance {
		builder.withDistance(b.distance)
	}
	return builder.build()
}
