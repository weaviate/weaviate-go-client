package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

func TestModules_EncodeDecode(t *testing.T) {
	var ms internal.Modules[string]
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
		if assert.Implements(t, (*internal.CustomModule[string])(nil), decoded) {
			assert.Subset(t, decoded.(internal.CustomModule[string]).Raw(), raw, "map is preserved")
		}
	})

	t.Run("custom key type roundtrip", func(t *testing.T) {
		var ms internal.Modules[str]
		ms.Register(*new(custom))

		five := custom{Value: 5}
		encoded, err := ms.Encode(five)
		require.NoError(t, err, "encode")

		decoded, err := ms.Decode("str-module", encoded)
		assert.NoError(t, err, "decode")
		assert.Equal(t, five, decoded, "decoded value")
	})
}

func TestModules_Find(t *testing.T) {
	var ms internal.Modules[string]
	ms.Register(*new(module))

	for _, tt := range []struct {
		name  string
		m     map[string]any
		found assert.BoolAssertionFunc
		want  string
	}{
		{
			name: "key present",
			m: map[string]any{
				"foo":         "bar",
				"hello":       "world",
				"test-module": 5,
			},
			found: assert.True,
			want:  "test-module",
		},
		{
			name: "key not present",
			m: map[string]any{
				"foo":   "bar",
				"hello": "world",
			},
			found: assert.False,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			key, ok := ms.Find(tt.m)

			if tt.found(t, ok, "found a match") {
				assert.Equal(t, tt.want, key, "module name")
			}
		})
	}
}

// module implements [internal.Module] for a simple struct.
type module struct {
	Value int `json:"value"`
}

var _ internal.Module[string] = (*module)(nil)

func (module) Name() string { return "test-module" }

// invalid implements [internal.Module] for an integer.
type invalid int

var _ internal.Module[string] = (*invalid)(nil)

func (invalid) Name() string { return "bad-module" }

type (
	custom struct{ Value int }
	str    string
)

var _ internal.Module[str] = (*custom)(nil)

func (custom) Name() str { return "str-module" }
