package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearImage(t *testing.T) {
	dirName := t.TempDir()
	fileName := dirName + "/image.file"

	t.Run("create dummy image file", func(t *testing.T) {
		file, err := os.Create(fileName)
		require.Nil(t, err)
		defer file.Close()

		written, err := file.Write([]byte("some image"))
		require.Nil(t, err)
		require.Greater(t, written, 0)
	})

	t.Run("from file", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearImage := (&NearImageArgumentBuilder{}).
			WithReader(file).
			build()

		expected := `nearImage:{image: "c29tZSBpbWFnZQ=="}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("from file with certainty", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearImage := (&NearImageArgumentBuilder{}).
			WithReader(file).
			WithCertainty(0.5).
			build()

		expected := `nearImage:{image: "c29tZSBpbWFnZQ==" certainty: 0.5}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("from file with distance", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearImage := (&NearImageArgumentBuilder{}).
			WithReader(file).
			WithDistance(0.5).
			build()

		expected := `nearImage:{image: "c29tZSBpbWFnZQ==" distance: 0.5}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("from base64", func(t *testing.T) {
		nearImage := (&NearImageArgumentBuilder{}).
			WithImage("iVBORw0KGgoAAAANS").
			build()

		expected := `nearImage:{image: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("from base64 with header", func(t *testing.T) {
		nearImage := (&NearImageArgumentBuilder{}).
			WithImage("data:image/mp4;base64,iVBORw0KGgoAAAANS").
			build()

		expected := `nearImage:{image: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("from base64 with certainty", func(t *testing.T) {
		nearImage := (&NearImageArgumentBuilder{}).
			WithImage("iVBORw0KGgoAAAANS").
			WithCertainty(0.5).
			build()

		expected := `nearImage:{image: "iVBORw0KGgoAAAANS" certainty: 0.5}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("from base64 with distance", func(t *testing.T) {
		nearImage := (&NearImageArgumentBuilder{}).
			WithImage("iVBORw0KGgoAAAANS").
			WithDistance(0.5).
			build()

		expected := `nearImage:{image: "iVBORw0KGgoAAAANS" distance: 0.5}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("empty", func(t *testing.T) {
		nearImage := (&NearImageArgumentBuilder{}).build()

		expected := `nearImage:{}`
		assert.Equal(t, expected, nearImage)
	})

	t.Run("from base64 with targetVectors", func(t *testing.T) {
		nearImage := (&NearImageArgumentBuilder{}).
			WithImage("iVBORw0KGgoAAAANS").
			WithTargetVectors("targetVector").
			build()

		expected := `nearImage:{image: "iVBORw0KGgoAAAANS" targetVectors: ["targetVector"]}`
		assert.Equal(t, expected, nearImage)
	})
}
