package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

func TestModules(t *testing.T) {
	var ms internal.Modules
	ms.Register(*new(module))

	t.Run("non-struct module", func(t *testing.T) {
		assert.Panics(t, func() {
			ms.Register(*new(invalid))
		})
	})

	t.Run("pass module by pointer reference", func(t *testing.T) {
		assert.Panics(t, func() {
			ms.Register(new(module))
		})
	})

	t.Run("register a nil module", func(t *testing.T) {
		assert.Panics(t, func() {
			ms.Register(nil)
		})
	})

	t.Run("registered module roundtrip", func(t *testing.T) {
		five := module{5}
		encoded, err := ms.Encode(five)
		require.NoError(t, err, "encode")
		assert.Subset(t, encoded, map[string]any{"value": 5}, "encoded map")

		decoded, err := ms.Decode("test-module", encoded)
		assert.NoError(t, err, "decode")
		assert.Equal(t, five, decoded, "decoded value")
	})

	t.Run("decode unknown module", func(t *testing.T) {
		raw := map[string]any{"value": 5}

		decoded, err := ms.Decode("unknown", raw)

		assert.NoError(t, err, "decode")
		require.NotNil(t, decoded, "decoded value")
		assert.Equal(t, "unknown", decoded.Name(), "module name")
		if assert.Implements(t, (*internal.CustomModule)(nil), decoded) {
			assert.Subset(t, decoded.(internal.CustomModule).Raw(), raw, "map is preserved")
		}
	})

	t.Run("custom decode hook", func(t *testing.T) {
		helloworld := decoder{"hello": "world"}
		ms.Register(helloworld)

		foobar := decoder{"foo": "bar"}
		encoded, err := ms.Encode(foobar)
		require.NoError(t, err, "encode")
		require.EqualValues(t, foobar, encoded, "encoded to its actual contents")

		decoded, err := ms.Decode("decoder-module", encoded)
		require.NoError(t, err, "decode")
		require.Subset(t, decoded, helloworld, "contains custom values")
		require.NotSubset(t, decoded, foobar, "contains unexpected values")
	})
}

// module implements [internal.Module] for a simple struct.
type module struct {
	Value int `json:"value"`
}

var _ internal.Module = (*module)(nil)

func (module) Name() string { return "test-module" }

// invalid implements [internal.Module] for an integer.
type invalid int

var _ internal.Module = (*invalid)(nil)

func (invalid) Name() string { return "bad-module" }

// decoder always decodes to the same map set by v.
type decoder map[string]any

var _ internal.Module = (*decoder)(nil)

func (decoder) Name() string { return "decoder-module" }
func (d decoder) Decode(map[string]any) (internal.Module, error) {
	return d, nil
}
