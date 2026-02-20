package weaviate

var version = "v6.0.0-alpha.1"

// Version reports the version of the package.
// This is sent in the X-Weaviate-Client header for telemetry.
func Version() string {
	return version
}
