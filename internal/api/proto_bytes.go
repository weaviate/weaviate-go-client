package api

// marshal encodes a 1-dimensional vector into a byte-slice.
func marshalSingle([]float32) []byte {
	return nil
}

// unmarshalSingle decodes a byte-slice into a 1-dimensional vector.
// It is safe to pass nil-slice -- the return is also a nil.
func unmarshalSingle([]byte) []float32 {
	return nil
}

// marshalMulti encodes a multi-dimensional vector into a byte-slice.
func marshalMulti([][]float32) []byte {
	return nil
}

// unmarshalMulti decodes a byte-slice into a multi-dimensional vector.
// It is safe to pass nil-slice -- the return is also a nil.
func unmarshalMulti([]byte) [][]float32 {
	return nil
}
