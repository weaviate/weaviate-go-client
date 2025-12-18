package weaviate

// Version reports the semver of this client.
// Usage:
//
//	$ version=$(git describe --tags --abbrev=0)
//	$ go build -ldflags "-X github.com/weaviate/weaviate-go-client/v6/weaviate.version=${version}" .
//
// Client also includes this value in the X-Weaviate-Client header for telemetry.
var version = "v6.0.0-alpha.1"

// Version reports the version of the package.
func Version() string {
	return version
}
