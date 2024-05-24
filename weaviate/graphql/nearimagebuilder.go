package graphql

import (
	"io"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type NearImageArgumentBuilder struct {
	image         string
	imageReader   io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
	targetVectors []string
}

// WithImage base64 encoded image
func (b *NearImageArgumentBuilder) WithImage(image string) *NearImageArgumentBuilder {
	b.image = image
	return b
}

// WithReader the image file
func (b *NearImageArgumentBuilder) WithReader(imageReader io.Reader) *NearImageArgumentBuilder {
	b.imageReader = imageReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearImageArgumentBuilder) WithCertainty(certainty float32) *NearImageArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearImageArgumentBuilder) WithDistance(distance float32) *NearImageArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

// WithTargetVectors target vector name
func (b *NearImageArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearImageArgumentBuilder {
	b.targetVectors = targetVectors
	return b
}

// Build build the given clause
func (b *NearImageArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:  "nearImage",
		mediaField: "image",
		data:       b.image,
		dataReader: b.imageReader,
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

func (b *NearImageArgumentBuilder) togrpc() *pb.NearImageSearch {
	builder := &nearMediaArgumentBuilder{
		data:       b.image,
		dataReader: b.imageReader,
	}
	nearImage := &pb.NearImageSearch{
		Image:         builder.getContent(),
		TargetVectors: b.targetVectors,
	}
	if b.hasCertainty {
		certainty := float64(b.certainty)
		nearImage.Certainty = &certainty
	}
	if b.hasDistance {
		distance := float64(b.distance)
		nearImage.Distance = &distance
	}
	return nearImage
}
