package testenv

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/test"
)

// SetupLocalWeaviate creates a local weaviate running on 8080 using docker compose
// Will only wait for it to be reachable if env `EXTERNAL_WEAVIATE_RUNNING` is set to True.
//
//	`EXTERNAL_WEAVIATE_RUNNING` should be set if all tests are supposed to be run in a test suit.
//	This prevents unnecessary starting and stopping of the docker-compose which prevents errors
//	due to syncing issues and speeds up the process
func SetupLocalWeaviate() error {
	if !isExternalWeaviateRunning() {
		return test.SetupWeaviate()
	}
	return nil
}

func isExternalWeaviateRunning() bool {
	val := os.Getenv("EXTERNAL_WEAVIATE_RUNNING")
	val = strings.ToLower(val)
	fmt.Printf("\nEXTERNAL_WEAVIATE_RUNNING: %v\n", val)
	return val == "true"
}

// TearDownLocalWeaviate shuts down the locally started weaviate docker compose
// If `EXTERNAL_WEAVIATE_RUNNING` this function will not do anything
//
//	see SetupLocalWeaviate for more info.
func TearDownLocalWeaviate() error {
	if isExternalWeaviateRunning() {
		return nil
	}
	err := test.TearDownWeaviate()
	time.Sleep(time.Second * 3) // Add some delay to make sure the command was executed before the program exits
	return err
}
