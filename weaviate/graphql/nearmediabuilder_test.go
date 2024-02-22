package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearMedia(t *testing.T) {
	dirName := t.TempDir()
	fileName := dirName + "/media.file"

	t.Run("create dummy media file", func(t *testing.T) {
		file, err := os.Create(fileName)
		require.Nil(t, err)
		defer file.Close()

		written, err := file.Write([]byte("some media"))
		require.Nil(t, err)
		require.Greater(t, written, 0)
	})

	t.Run("from file", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
			dataReader: file,
		}).
			build()

		expected := `nearMedia:{media: "c29tZSBtZWRpYQ=="}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("from file with certainty", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
			dataReader: file,
		}).
			withCertainty(0.5).
			build()

		expected := `nearMedia:{media: "c29tZSBtZWRpYQ==" certainty: 0.5}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("from file with distance", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
			dataReader: file,
		}).
			withDistance(0.5).
			build()

		expected := `nearMedia:{media: "c29tZSBtZWRpYQ==" distance: 0.5}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("from base64", func(t *testing.T) {
		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
			data:       "iVBORw0KGgoAAAANS",
		}).
			build()

		expected := `nearMedia:{media: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("from base64 with header", func(t *testing.T) {
		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
			data:       "data:media/mp4;base64,iVBORw0KGgoAAAANS",
		}).
			build()

		expected := `nearMedia:{media: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("from base64 with certainty", func(t *testing.T) {
		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
			data:       "iVBORw0KGgoAAAANS",
		}).
			withCertainty(0.5).
			build()

		expected := `nearMedia:{media: "iVBORw0KGgoAAAANS" certainty: 0.5}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("from base64 with distance", func(t *testing.T) {
		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
			data:       "iVBORw0KGgoAAAANS",
		}).
			withDistance(0.5).
			build()

		expected := `nearMedia:{media: "iVBORw0KGgoAAAANS" distance: 0.5}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("empty", func(t *testing.T) {
		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:  "nearMedia",
			mediaField: "media",
		}).build()

		expected := `nearMedia:{}`
		assert.Equal(t, expected, nearMedia)
	})

	t.Run("from base64 with targetVectors", func(t *testing.T) {
		nearMedia := (&nearMediaArgumentBuilder{
			mediaName:     "nearMedia",
			mediaField:    "media",
			data:          "iVBORw0KGgoAAAANS",
			targetVectors: []string{"targetVector"},
		}).
			build()

		expected := `nearMedia:{media: "iVBORw0KGgoAAAANS" targetVectors: ["targetVector"]}`
		assert.Equal(t, expected, nearMedia)
	})
}
