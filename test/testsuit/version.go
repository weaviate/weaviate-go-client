package testsuit

import (
	"testing"

	"github.com/launchdarkly/go-semver"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// AtLeastWeaviateVersion skips the test if the weaviate version is lower than the required version.
func AtLeastWeaviateVersion(t *testing.T, client *weaviate.Client, requiredVersion, message string) {
	if !ServerAtLeast(t, client, requiredVersion) {
		t.Skip(message)
	}
}

// ServerAtLeast reports whether the connected server is at least requiredVersion.
// Use it when a test needs to branch on server version; use AtLeastWeaviateVersion
// when it should skip.
func ServerAtLeast(t *testing.T, client *weaviate.Client, requiredVersion string) bool {
	meta, err := client.Misc().MetaGetter().Do(t.Context())
	require.Nil(t, err, "could not get weaviate meta information")

	runningVersion, err := semver.Parse(meta.Version)
	require.Nil(t, err, "could not parse weaviate server version %q", meta.Version)

	minVersion, err := semver.Parse(requiredVersion)
	require.Nil(t, err, "could not parse required version %q", requiredVersion)

	return runningVersion.ComparePrecedence(minVersion) >= 0
}
