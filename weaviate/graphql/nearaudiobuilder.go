package graphql

import (
	"io"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type NearAudioArgumentBuilder struct {
	audio         string
	audioReader   io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
	targetVectors []string
	targets       *MultiTargetArgumentBuilder
}

// WithAudio base64 encoded audio
func (b *NearAudioArgumentBuilder) WithAudio(audio string) *NearAudioArgumentBuilder {
	b.audio = audio
	return b
}

// WithReader the audio file
func (b *NearAudioArgumentBuilder) WithReader(audioReader io.Reader) *NearAudioArgumentBuilder {
	b.audioReader = audioReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearAudioArgumentBuilder) WithCertainty(certainty float32) *NearAudioArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearAudioArgumentBuilder) WithDistance(distance float32) *NearAudioArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

// WithTargetVectors target vector name
func (b *NearAudioArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearAudioArgumentBuilder {
	b.targetVectors = targetVectors
	return b
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (h *NearAudioArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *NearAudioArgumentBuilder {
	h.targets = targets
	return h
}

// Build build the given clause
func (b *NearAudioArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:  "nearAudio",
		mediaField: "audio",
		data:       b.audio,
		dataReader: b.audioReader,
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

func (b *NearAudioArgumentBuilder) togrpc() *pb.NearAudioSearch {
	builder := &nearMediaArgumentBuilder{
		data:       b.audio,
		dataReader: b.audioReader,
	}
	nearAudio := &pb.NearAudioSearch{
		Audio: builder.getContent(),
	}
	if b.hasCertainty {
		certainty := float64(b.certainty)
		nearAudio.Certainty = &certainty
	}
	if b.hasDistance {
		distance := float64(b.distance)
		nearAudio.Distance = &distance
	}
	if b.targets != nil {
		nearAudio.Targets = b.targets.togrpc()
	}
	if len(b.targetVectors) > 0 && b.targets == nil {
		nearAudio.Targets = &pb.Targets{TargetVectors: b.targetVectors}
	}
	return nearAudio
}
