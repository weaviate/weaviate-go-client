package graphql

import (
	"io"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type NearVideoArgumentBuilder struct {
	video         string
	videoReader   io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
	targetVectors []string
}

// WithVideo base64 encoded video
func (b *NearVideoArgumentBuilder) WithVideo(video string) *NearVideoArgumentBuilder {
	b.video = video
	return b
}

// WithReader the video file
func (b *NearVideoArgumentBuilder) WithReader(videoReader io.Reader) *NearVideoArgumentBuilder {
	b.videoReader = videoReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearVideoArgumentBuilder) WithCertainty(certainty float32) *NearVideoArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearVideoArgumentBuilder) WithDistance(distance float32) *NearVideoArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

// WithTargetVectors target vector name
func (b *NearVideoArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearVideoArgumentBuilder {
	b.targetVectors = targetVectors
	return b
}

// Build build the given clause
func (b *NearVideoArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:  "nearVideo",
		mediaField: "video",
		data:       b.video,
		dataReader: b.videoReader,
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

func (b *NearVideoArgumentBuilder) togrpc() *pb.NearVideoSearch {
	builder := &nearMediaArgumentBuilder{
		data:       b.video,
		dataReader: b.videoReader,
	}
	nearVideo := &pb.NearVideoSearch{
		Video:         builder.getContent(),
		TargetVectors: b.targetVectors,
	}
	if b.hasCertainty {
		certainty := float64(b.certainty)
		nearVideo.Certainty = &certainty
	}
	if b.hasDistance {
		distance := float64(b.distance)
		nearVideo.Distance = &distance
	}
	return nearVideo
}
