package internal

// vector represents a vector that can be 1D (single vector) or 2D (multi-vector).
// It provides methods to check its type and convert it to different formats.
type Vector interface {
	IsMulti() bool

	// ToFloat32 converts to []float32 for single vectors.
	// Panics if called on a multi-vector.
	ToFloat32() []float32

	// ToFloat32Multi converts to [][]float32 for multi-vectors.
	// For single vectors, returns a slice containing one vector.
	ToFloat32Multi() [][]float32

	// ToBytes converts the vector to bytes for gRPC transmission.
	ToBytes() []byte
}
