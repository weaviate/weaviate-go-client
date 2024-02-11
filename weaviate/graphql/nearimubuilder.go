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
	if len(targetVectors) > 0 {
		b.targetVectors = targetVectors
	}
	return b
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
	return builder.build()
}
