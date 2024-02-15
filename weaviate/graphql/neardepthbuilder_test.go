package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearDepth(t *testing.T) {
	dirName := t.TempDir()
	fileName := dirName + "/depth.file"

	t.Run("create dummy depth file", func(t *testing.T) {
		file, err := os.Create(fileName)
		require.Nil(t, err)
		defer file.Close()

		written, err := file.Write([]byte("some depth"))
		require.Nil(t, err)
		require.Greater(t, written, 0)
	})

	t.Run("from file", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearDepth := (&NearDepthArgumentBuilder{}).
			WithReader(file).
			build()

		expected := `nearDepth:{depth: "c29tZSBkZXB0aA=="}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("from file with certainty", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearDepth := (&NearDepthArgumentBuilder{}).
			WithReader(file).
			WithCertainty(0.5).
			build()

		expected := `nearDepth:{depth: "c29tZSBkZXB0aA==" certainty: 0.5}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("from file with distance", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearDepth := (&NearDepthArgumentBuilder{}).
			WithReader(file).
			WithDistance(0.5).
			build()

		expected := `nearDepth:{depth: "c29tZSBkZXB0aA==" distance: 0.5}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("from base64", func(t *testing.T) {
		nearDepth := (&NearDepthArgumentBuilder{}).
			WithDepth("iVBORw0KGgoAAAANS").
			build()

		expected := `nearDepth:{depth: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("from base64 with header", func(t *testing.T) {
		nearDepth := (&NearDepthArgumentBuilder{}).
			WithDepth("data:depth/mp4;base64,iVBORw0KGgoAAAANS").
			build()

		expected := `nearDepth:{depth: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("from base64 with certainty", func(t *testing.T) {
		nearDepth := (&NearDepthArgumentBuilder{}).
			WithDepth("iVBORw0KGgoAAAANS").
			WithCertainty(0.5).
			build()

		expected := `nearDepth:{depth: "iVBORw0KGgoAAAANS" certainty: 0.5}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("from base64 with distance", func(t *testing.T) {
		nearDepth := (&NearDepthArgumentBuilder{}).
			WithDepth("iVBORw0KGgoAAAANS").
			WithDistance(0.5).
			build()

		expected := `nearDepth:{depth: "iVBORw0KGgoAAAANS" distance: 0.5}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("empty", func(t *testing.T) {
		nearDepth := (&NearDepthArgumentBuilder{}).build()

		expected := `nearDepth:{}`
		assert.Equal(t, expected, nearDepth)
	})

	t.Run("from base64 with targetVectors", func(t *testing.T) {
		nearDepth := (&NearDepthArgumentBuilder{}).
			WithDepth("iVBORw0KGgoAAAANS").
			WithTargetVectors("targetVector").
			build()

		expected := `nearDepth:{depth: "iVBORw0KGgoAAAANS" targetVectors: ["targetVector"]}`
		assert.Equal(t, expected, nearDepth)
	})
}
