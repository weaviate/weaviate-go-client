package graphql

import (
	"io"
)

type NearImuArgumentBuilder struct {
	imu           string
	imuReader     io.Reader
	hasCertainty  bool
	certainty     float32
	hasDistance   bool
	distance      float32
	targetVectors []string
	targets       *MultiTargetArgumentBuilder
}

// WithImu base64 encoded imu
func (b *NearImuArgumentBuilder) WithImu(imu string) *NearImuArgumentBuilder {
	b.imu = imu
	return b
}

// WithReader the imu file
func (b *NearImuArgumentBuilder) WithReader(imuReader io.Reader) *NearImuArgumentBuilder {
	b.imuReader = imuReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearImuArgumentBuilder) WithCertainty(certainty float32) *NearImuArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearImuArgumentBuilder) WithDistance(distance float32) *NearImuArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

// WithTargetVectors target vector name
func (b *NearImuArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearImuArgumentBuilder {
	b.targetVectors = targetVectors
	return b
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (h *NearImuArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *NearImuArgumentBuilder {
	h.targets = targets
	return h
}

// Build build the given clause
func (b *NearImuArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:     "nearIMU",
		mediaField:    "imu",
		data:          b.imu,
		dataReader:    b.imuReader,
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
	builder.withTargets((b.targets))
	return builder.build()
}
