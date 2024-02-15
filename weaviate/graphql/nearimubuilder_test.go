package graphql

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearImu(t *testing.T) {
	dirName := t.TempDir()
	fileName := dirName + "/imu.file"

	t.Run("create dummy imu file", func(t *testing.T) {
		file, err := os.Create(fileName)
		require.Nil(t, err)
		defer file.Close()

		written, err := file.Write([]byte("some imu"))
		require.Nil(t, err)
		require.Greater(t, written, 0)
	})

	t.Run("from file", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearImu := (&NearImuArgumentBuilder{}).
			WithReader(file).
			build()

		expected := `nearIMU:{imu: "c29tZSBpbXU="}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("from file with certainty", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearImu := (&NearImuArgumentBuilder{}).
			WithReader(file).
			WithCertainty(0.5).
			build()

		expected := `nearIMU:{imu: "c29tZSBpbXU=" certainty: 0.5}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("from file with distance", func(t *testing.T) {
		file, err := os.Open(fileName)
		require.Nil(t, err)
		require.NotNil(t, file)
		defer file.Close()

		nearImu := (&NearImuArgumentBuilder{}).
			WithReader(file).
			WithDistance(0.5).
			build()

		expected := `nearIMU:{imu: "c29tZSBpbXU=" distance: 0.5}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("from base64", func(t *testing.T) {
		nearImu := (&NearImuArgumentBuilder{}).
			WithImu("iVBORw0KGgoAAAANS").
			build()

		expected := `nearIMU:{imu: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("from base64 with header", func(t *testing.T) {
		nearImu := (&NearImuArgumentBuilder{}).
			WithImu("data:imu/mp4;base64,iVBORw0KGgoAAAANS").
			build()

		expected := `nearIMU:{imu: "iVBORw0KGgoAAAANS"}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("from base64 with certainty", func(t *testing.T) {
		nearImu := (&NearImuArgumentBuilder{}).
			WithImu("iVBORw0KGgoAAAANS").
			WithCertainty(0.5).
			build()

		expected := `nearIMU:{imu: "iVBORw0KGgoAAAANS" certainty: 0.5}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("from base64 with distance", func(t *testing.T) {
		nearImu := (&NearImuArgumentBuilder{}).
			WithImu("iVBORw0KGgoAAAANS").
			WithDistance(0.5).
			build()

		expected := `nearIMU:{imu: "iVBORw0KGgoAAAANS" distance: 0.5}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("empty", func(t *testing.T) {
		nearImu := (&NearImuArgumentBuilder{}).build()

		expected := `nearIMU:{}`
		assert.Equal(t, expected, nearImu)
	})

	t.Run("from base64 with targetVectors", func(t *testing.T) {
		nearImu := (&NearImuArgumentBuilder{}).
			WithImu("iVBORw0KGgoAAAANS").
			WithTargetVectors("targetVector").
			build()

		expected := `nearIMU:{imu: "iVBORw0KGgoAAAANS" targetVectors: ["targetVector"]}`
		assert.Equal(t, expected, nearImu)
	})
}
