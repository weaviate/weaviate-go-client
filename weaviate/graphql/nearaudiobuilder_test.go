package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearAudio(t *testing.T) {
	dirName := t.TempDir()
	fileName := dirName + "/audio.file"

	t.Run("create dummy audio file", func(t *testing.T) {
		file, err := os.Create(fileName)
		require.Nil(t, err)
		defer file.Close()

		written, err := file.Write([]byte("some audio"))
		require.Nil(t, err)
		require.Greater(t, written, 0)
	})

	t.Run("from file", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearAudio := (&NearAudioArgumentBuilder{}).
			WithReader(file).
			build()

		expected := `nearAudio:{audio: "c29tZSBhdWRpbw=="}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("from file with certainty", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearAudio := (&NearAudioArgumentBuilder{}).
			WithReader(file).
			WithCertainty(0.5).
			build()

		expected := `nearAudio:{audio: "c29tZSBhdWRpbw==" certainty: 0.5}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("from file with distance", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearAudio := (&NearAudioArgumentBuilder{}).
			WithReader(file).
			WithDistance(0.5).
			build()

		expected := `nearAudio:{audio: "c29tZSBhdWRpbw==" distance: 0.5}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("from base64", func(t *testing.T) {
		nearAudio := (&NearAudioArgumentBuilder{}).
			WithAudio("iVBORw0KGgoAAAANS").
			build()

		expected := `nearAudio:{audio: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("from base64 with header", func(t *testing.T) {
		nearAudio := (&NearAudioArgumentBuilder{}).
			WithAudio("data:audio/mp4;base64,iVBORw0KGgoAAAANS").
			build()

		expected := `nearAudio:{audio: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("from base64 with certainty", func(t *testing.T) {
		nearAudio := (&NearAudioArgumentBuilder{}).
			WithAudio("iVBORw0KGgoAAAANS").
			WithCertainty(0.5).
			build()

		expected := `nearAudio:{audio: "iVBORw0KGgoAAAANS" certainty: 0.5}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("from base64 with distance", func(t *testing.T) {
		nearAudio := (&NearAudioArgumentBuilder{}).
			WithAudio("iVBORw0KGgoAAAANS").
			WithDistance(0.5).
			build()

		expected := `nearAudio:{audio: "iVBORw0KGgoAAAANS" distance: 0.5}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("empty", func(t *testing.T) {
		nearAudio := (&NearAudioArgumentBuilder{}).build()

		expected := `nearAudio:{}`
		assert.Equal(t, expected, nearAudio)
	})

	t.Run("from base64 with targetVectors", func(t *testing.T) {
		nearAudio := (&NearAudioArgumentBuilder{}).
			WithAudio("data:audio/mp4;base64,iVBORw0KGgoAAAANS").
			WithTargetVectors("targetVector").
			build()

		expected := `nearAudio:{audio: "iVBORw0KGgoAAAANS" targetVectors: ["targetVector"]}`
		assert.Equal(t, expected, nearAudio)
	})
}
