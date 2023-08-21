package graphql

import (
	"io"
)

type NearDepthArgumentBuilder struct {
	depth        string
	depthReader  io.Reader
	hasCertainty bool
	certainty    float32
	hasDistance  bool
	distance     float32
}

// WithDepth base64 encoded depth
func (b *NearDepthArgumentBuilder) WithDepth(depth string) *NearDepthArgumentBuilder {
	b.depth = depth
	return b
}

// WithReader the depth file
func (b *NearDepthArgumentBuilder) WithReader(depthReader io.Reader) *NearDepthArgumentBuilder {
	b.depthReader = depthReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearDepthArgumentBuilder) WithCertainty(certainty float32) *NearDepthArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearDepthArgumentBuilder) WithDistance(distance float32) *NearDepthArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

// Build build the given clause
func (b *NearDepthArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:  "nearDepth",
		mediaField: "depth",
		data:       b.depth,
		dataReader: b.depthReader,
	}
	if b.hasCertainty {
		builder.withCertainty(b.certainty)
	}
	if b.hasDistance {
		builder.withDistance(b.distance)
	}
	return builder.build()
}
