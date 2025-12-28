package api

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestVectorBytes(t *testing.T) {
	t.Run("trivial single", func(t *testing.T) {
		v := []float32{1, 2, 3}
		want := []byte{
			0x0, 0x0, 0x80, 0x3f,
			0x0, 0x0, 0x0, 0x40,
			0x0, 0x0, 0x40, 0x40,
		}
		got := marshalSingle(v)
		require.Equal(t, want, got, "bad vector bytes")
	})

	t.Run("trivial multi", func(t *testing.T) {
		v := [][]float32{{1, 2, 3}, {1, 2, 3}}
		want := []byte{
			0x3, 0x0, // inner array size, uint16(3)
			0x0, 0x0, 0x80, 0x3f, // first vector
			0x0, 0x0, 0x0, 0x40,
			0x0, 0x0, 0x40, 0x40,
			0x0, 0x0, 0x80, 0x3f, // second vector
			0x0, 0x0, 0x0, 0x40,
			0x0, 0x0, 0x40, 0x40,
		}
		got := marshalMulti(v)
		require.Equal(t, want, got, "bad vector bytes")
	})

	prng := testkit.NewPRNG(t)

	const (
		sizeof_fp32   = 4 // float32 size in bytes
		sizeof_uint16 = 2 // uint16 size in bytes

		// We limit the size of the vectors to produce readable
		// test output in case of an error.
		maxVectorDim = 16 // Largest single vector
		maxNestedDim = 4  // Largest number of nested vectors
	)

	t.Run("fuzz single", func(t *testing.T) {
		for range 1_000 {
			// Arrange: generate randomly sized vector
			dim := prng.Intn(maxVectorDim + 1)
			v := make([]float32, dim)
			for i := range v {
				v[i] = float32(prng.NormFloat64())
			}

			// Act, Assert: convert to []byte and back, verifying the properties hold:
			// - size of the []byte == sizeof(float32) * dim
			// - resulting vector is equal to the input

			b := marshalSingle(v)
			require.Len(t, b, sizeof_fp32*dim, "wrong vector bytes size")

			v2 := unmarshalSingle(b)
			require.Len(t, v2, len(v), "vectors size differs after roundtrip")
			if len(v) != 0 {
				require.Equal(t, v, v2, "vectors do not match after roundtrip")
			}
		}
	})

	t.Run("fuzz multi", func(t *testing.T) {
		for range 1_000 {
			// Arrange: generate randomly sized vector
			nested := prng.Intn(maxNestedDim + 1)
			dim := prng.Intn(maxVectorDim) + 1 // nested vectors cannot be empty

			vv := make([][]float32, nested)
			for ii := range vv {
				v := make([]float32, dim)
				for i := range v {
					v[i] = float32(prng.NormFloat64())
				}
				vv[ii] = v
			}

			// Act, Assert: convert to []byte and back, verifying the properties hold:
			// - size of the []byte == sizeof(float32) * dim
			// - resulting vector is equal to the input

			b := marshalMulti(vv)
			want := sizeof_uint16 + nested*sizeof_fp32*dim
			if nested == 0 {
				want = 0
			}
			require.Len(t, b, want, "%v - wrong vector bytes size", vv)

			vv2 := unmarshalMulti(b)
			require.Lenf(t, vv2, len(vv), "vectors size differs after roundtrip\n\tinitial: %v", vv)
			if len(vv) != 0 {
				require.Equalf(t, vv, vv2, "vectors do not match after roundtrip")
			}
		}
	})
}
