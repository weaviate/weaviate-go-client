package api

import (
	"testing"

	"github.com/stretchr/testify/require"
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
}
