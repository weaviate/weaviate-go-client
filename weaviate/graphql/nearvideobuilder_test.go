package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearVideo(t *testing.T) {
	dirName := t.TempDir()
	fileName := dirName + "/video.file"

	t.Run("create dummy video file", func(t *testing.T) {
		file, err := os.Create(fileName)
		require.Nil(t, err)
		defer file.Close()

		written, err := file.Write([]byte("some video"))
		require.Nil(t, err)
		require.Greater(t, written, 0)
	})

	t.Run("from file", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearVideo := (&NearVideoArgumentBuilder{}).
			WithReader(file).
			build()

		expected := `nearVideo:{video: "c29tZSB2aWRlbw=="}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("from file with certainty", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearVideo := (&NearVideoArgumentBuilder{}).
			WithReader(file).
			WithCertainty(0.5).
			build()

		expected := `nearVideo:{video: "c29tZSB2aWRlbw==" certainty: 0.5}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("from file with distance", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearVideo := (&NearVideoArgumentBuilder{}).
			WithReader(file).
			WithDistance(0.5).
			build()

		expected := `nearVideo:{video: "c29tZSB2aWRlbw==" distance: 0.5}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("from base64", func(t *testing.T) {
		nearVideo := (&NearVideoArgumentBuilder{}).
			WithVideo("iVBORw0KGgoAAAANS").
			build()

		expected := `nearVideo:{video: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("from base64 with header", func(t *testing.T) {
		nearVideo := (&NearVideoArgumentBuilder{}).
			WithVideo("data:video/mp4;base64,iVBORw0KGgoAAAANS").
			build()

		expected := `nearVideo:{video: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("from base64 with certainty", func(t *testing.T) {
		nearVideo := (&NearVideoArgumentBuilder{}).
			WithVideo("iVBORw0KGgoAAAANS").
			WithCertainty(0.5).
			build()

		expected := `nearVideo:{video: "iVBORw0KGgoAAAANS" certainty: 0.5}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("from base64 with distance", func(t *testing.T) {
		nearVideo := (&NearVideoArgumentBuilder{}).
			WithVideo("iVBORw0KGgoAAAANS").
			WithDistance(0.5).
			build()

		expected := `nearVideo:{video: "iVBORw0KGgoAAAANS" distance: 0.5}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("empty", func(t *testing.T) {
		nearVideo := (&NearVideoArgumentBuilder{}).build()

		expected := `nearVideo:{}`
		assert.Equal(t, expected, nearVideo)
	})

	t.Run("from base64 with targetVectors", func(t *testing.T) {
		nearVideo := (&NearVideoArgumentBuilder{}).
			WithVideo("iVBORw0KGgoAAAANS").
			WithTargetVectors("targetVector").
			build()

		expected := `nearVideo:{video: "iVBORw0KGgoAAAANS" targetVectors: ["targetVector"]}`
		assert.Equal(t, expected, nearVideo)
	})
}
