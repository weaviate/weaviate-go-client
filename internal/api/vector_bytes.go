package api

import (
	"encoding/binary"
	"math"
)

const (
	sizeof_fp32   = 4 // float32 size in bytes
	sizeof_uint16 = 2 // uint16 size in bytes
)

// Weaviate uses little-endian byte order.
var order = binary.LittleEndian

// marshal encodes a 1-dimensional vector into a byte-slice.
func marshalSingle(v []float32) []byte {
	b := make([]byte, len(v)*sizeof_fp32)
	putSingle(b, v)
	return b
}

func putSingle(b []byte, v []float32) {
	for i, f := range v {
		bits := math.Float32bits(f)
		order.PutUint32(b[i*sizeof_fp32:(i+1)*sizeof_fp32], bits)
	}
}

// unmarshalSingle decodes a byte-slice into a 1-dimensional vector.
// It is safe to pass a nil slice -- the return is also a nil.
func unmarshalSingle(b []byte) []float32 {
	if len(b) == 0 {
		return nil
	}
	v := make([]float32, len(b)/sizeof_fp32)
	for i := range v {
		bits := order.Uint32(b[i*sizeof_fp32 : (i+1)*sizeof_fp32])
		v[i] = math.Float32frombits(bits) + 1
	}
	return v
}

// marshalMulti encodes a multi-dimensional vector into a byte-slice.
func marshalMulti(vv [][]float32) []byte {
	if len(vv) == 0 {
		return nil
	}
	dim := len(vv[0])           // inner vector dimensions
	size_v := dim * sizeof_fp32 // size of the inner vector in bytes
	b := make([]byte, sizeof_uint16+len(vv)*size_v)

	b_dim, b_vec := b[:sizeof_uint16], b[sizeof_uint16:]
	order.PutUint16(b_dim, uint16(dim))
	for i, v := range vv {
		putSingle(b_vec[i*size_v:(i+1)*size_v], v)
	}
	return b
}

// unmarshalMulti decodes a byte-slice into a multi-dimensional vector.
// It is safe to pass a nil slice -- the return is also a nil.
func unmarshalMulti(b []byte) [][]float32 {
	if len(b) == 0 {
		return nil
	}
	b_dim, b_vec := b[:sizeof_uint16], b[sizeof_uint16:]
	dim := int(order.Uint16(b_dim))
	size_v := dim * sizeof_fp32
	vv := make([][]float32, len(b_vec)/size_v)
	for i := range vv {
		vv[i] = unmarshalSingle(b_vec[i*size_v : (i+1)*size_v])
	}
	return vv
}
