package testsuit

import (
	"testing"

	"github.com/launchdarkly/go-semver"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// AtLeastWeaviateVersion skips the test if the weaviate version is lower than the required version.
func AtLeastWeaviateVersion(t *testing.T, client *weaviate.Client, requiredVersion, message string) {
	meta, err := client.Misc().MetaGetter().Do(t.Context())
	require.Nil(t, err, "could not get weaviate meta information")

	runningVersion, err := semver.Parse(meta.Version)
	require.Nil(t, err, "could not parse WEAVIATE_VERSION env var")

	minVersion, err := semver.Parse(requiredVersion)
	require.Nil(t, err, "could not parse required version")

	if runningVersion.ComparePrecedence(minVersion) < 0 {
		t.Skip(message)
	}
}
