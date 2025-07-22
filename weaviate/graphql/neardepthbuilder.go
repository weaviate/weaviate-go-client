package graphql

import (
	"io"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type NearDepthArgumentBuilder struct {
	depth         string
	depthReader   io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
	targetVectors []string
	targets       *MultiTargetArgumentBuilder
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

// WithTargetVectors target vector name
func (b *NearDepthArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearDepthArgumentBuilder {
	b.targetVectors = targetVectors
	return b
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (h *NearDepthArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *NearDepthArgumentBuilder {
	h.targets = targets
	return h
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
	if len(b.targetVectors) > 0 {
		builder.withTargetVectors(b.targetVectors...)
	}
	builder.withTargets(b.targets)
	return builder.build()
}

func (b *NearDepthArgumentBuilder) togrpc() *pb.NearDepthSearch {
	builder := &nearMediaArgumentBuilder{
		data:       b.depth,
		dataReader: b.depthReader,
	}
	nearDepth := &pb.NearDepthSearch{
		Depth: builder.getContent(),
	}
	if b.hasCertainty {
		certainty := float64(b.certainty)
		nearDepth.Certainty = &certainty
	}
	if b.hasDistance {
		distance := float64(b.distance)
		nearDepth.Distance = &distance
	}
	if b.targets != nil {
		nearDepth.Targets = b.targets.togrpc()
	} else if len(b.targetVectors) > 0 {
		nearDepth.Targets = &pb.Targets{TargetVectors: b.targetVectors}
	}
	return nearDepth
}
