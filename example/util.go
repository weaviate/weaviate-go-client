package example

import (
	"fmt"
	"os"
)

// Catch panics for a non-nil error.
func Catch(err error) { Assert(err == nil, err) }

// Assert panics with message if check is false.
func Assert(check bool, a ...any) {
	if !check {
		panic(fmt.Sprint(a...))
	}
}

const (
	EnvHost   = "WEAVIATE_HOST"
	EnvAPIKey = "WEAVIATE_API_KEY"
)

// ConnectionParams reads host address and WCD API key from environment.
// Panics if either of these variables is not defined.
func ConnectionParams() (host string, apiKey string) {
	host, hasHost := os.LookupEnv(EnvHost)
	apiKey, hasKey := os.LookupEnv(EnvAPIKey)
	if !hasHost || !hasKey {
		panic(fmt.Sprintf("%q and %q must be defined", EnvHost, EnvAPIKey))
	}
	return
}
