package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearThermal(t *testing.T) {
	dirName := t.TempDir()
	fileName := dirName + "/thermal.file"

	t.Run("create dummy thermal file", func(t *testing.T) {
		file, err := os.Create(fileName)
		require.Nil(t, err)
		defer file.Close()

		written, err := file.Write([]byte("some thermal"))
		require.Nil(t, err)
		require.Greater(t, written, 0)
	})

	t.Run("from file", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearThermal := (&NearThermalArgumentBuilder{}).
			WithReader(file).
			build()

		expected := `nearThermal:{thermal: "c29tZSB0aGVybWFs"}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("from file with certainty", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearThermal := (&NearThermalArgumentBuilder{}).
			WithReader(file).
			WithCertainty(0.5).
			build()

		expected := `nearThermal:{thermal: "c29tZSB0aGVybWFs" certainty: 0.5}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("from file with distance", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearThermal := (&NearThermalArgumentBuilder{}).
			WithReader(file).
			WithDistance(0.5).
			build()

		expected := `nearThermal:{thermal: "c29tZSB0aGVybWFs" distance: 0.5}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("from base64", func(t *testing.T) {
		nearThermal := (&NearThermalArgumentBuilder{}).
			WithThermal("iVBORw0KGgoAAAANS").
			build()

		expected := `nearThermal:{thermal: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("from base64 with header", func(t *testing.T) {
		nearThermal := (&NearThermalArgumentBuilder{}).
			WithThermal("data:thermal/mp4;base64,iVBORw0KGgoAAAANS").
			build()

		expected := `nearThermal:{thermal: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("from base64 with certainty", func(t *testing.T) {
		nearThermal := (&NearThermalArgumentBuilder{}).
			WithThermal("iVBORw0KGgoAAAANS").
			WithCertainty(0.5).
			build()

		expected := `nearThermal:{thermal: "iVBORw0KGgoAAAANS" certainty: 0.5}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("from base64 with distance", func(t *testing.T) {
		nearThermal := (&NearThermalArgumentBuilder{}).
			WithThermal("iVBORw0KGgoAAAANS").
			WithDistance(0.5).
			build()

		expected := `nearThermal:{thermal: "iVBORw0KGgoAAAANS" distance: 0.5}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("empty", func(t *testing.T) {
		nearThermal := (&NearThermalArgumentBuilder{}).build()

		expected := `nearThermal:{}`
		assert.Equal(t, expected, nearThermal)
	})

	t.Run("from base64 with targetVectors", func(t *testing.T) {
		nearThermal := (&NearThermalArgumentBuilder{}).
			WithThermal("iVBORw0KGgoAAAANS").
			WithTargetVectors("targetVector").
			build()

		expected := `nearThermal:{thermal: "iVBORw0KGgoAAAANS" targetVectors: ["targetVector"]}`
		assert.Equal(t, expected, nearThermal)
	})
}
