package api

import (
	"errors"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/semver"
)

const (
	LatestSupportedVersion   = "v1.36"
	EarliestSupportedVersion = "v1.28"
)

// errVersionNotSupported is returned when the Weaviate server version
// is outside of the supported version range.
var errVersionNotSupported = errors.New("server version is not supported")

// isVersionSupported returns true if server version v lies within
// [EarliestSupportedVersion, LatestSupportedVersion] range.
func isVersionSupported(v string) bool {
	return semver.AfterMajorMinor(v, EarliestSupportedVersion) &&
		(semver.BeforeMajorMinor(v, LatestSupportedVersion) ||
			semver.EqualMajorMinor(v, LatestSupportedVersion))
}
