package graphql

import (
	"io"
)

type NearAudioArgumentBuilder struct {
	audio         string
	audioReader   io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
	targetVectors []string
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
	if len(targetVectors) > 0 {
		b.targetVectors = targetVectors
	}
	return b
}

// Build build the given clause
func (b *NearAudioArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:     "nearAudio",
		mediaField:    "audio",
		data:          b.audio,
		dataReader:    b.audioReader,
		targetVectors: b.targetVectors,
	}
	if b.hasCertainty {
		builder.withCertainty(b.certainty)
	}
	if b.hasDistance {
		builder.withDistance(b.distance)
	}
	return builder.build()
}
