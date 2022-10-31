package testenv

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/semi-technologies/weaviate-go-client/v5/test"
	"github.com/semi-technologies/weaviate-go-client/v5/weaviate"
)

// SetupLocalWeaviate creates a local weaviate running on 8080 using docker compose
// Will only wait for it to be reachable if env `EXTERNAL_WEAVIATE_RUNNING` is set to True.
//
//	`EXTERNAL_WEAVIATE_RUNNING` should be set if all tests are supposed to be run in a test suit.
//	This prevents unnecessary starting and stopping of the docker-compose which prevents errors
//	due to syncing issues and speeds up the process
func SetupLocalWeaviate() error {
	if !isExternalWeaviateRunning() {
		err := test.SetupWeaviate()
		if err != nil {
			return err
		}
	}
	return WaitForWeaviate()
}

// SetupLocalWeaviateDeprecated creates a local weaviate running on 8080 using docker compose
// for pre-v1.14 backwards compatibility tests.
// Will only wait for it to be reachable if env `EXTERNAL_WEAVIATE_RUNNING` is set to True.
//
//	`EXTERNAL_WEAVIATE_RUNNING` should be set if all tests are supposed to be run in a test suit.
//	This prevents unnecessary starting and stopping of the docker-compose which prevents errors
//	due to syncing issues and speeds up the process
func SetupLocalWeaviateDeprecated() error {
	if !isExternalWeaviateRunning() {
		err := test.SetupWeaviateDeprecated()
		if err != nil {
			return err
		}
	}
	return WaitForWeaviate()
}

func isExternalWeaviateRunning() bool {
	val := os.Getenv("EXTERNAL_WEAVIATE_RUNNING")
	val = strings.ToLower(val)
	fmt.Printf("\nEXTERNAL_WEAVIATE_RUNNING: %v\n", val)
	if val == "true" {
		return true
	}
	return false
}

// WaitForWeaviate waits until weaviate is started up and ready
// expects weaviat at localhost:8080
func WaitForWeaviate() error {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client := weaviate.New(cfg)

	for i := 0; i < 20; i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		isReady, _ := client.Misc().ReadyChecker().Do(ctx)
		if isReady {
			return nil
		}
		fmt.Printf("Weaviate not yet up waiting another 3 seconds. Iteration: %v\n", i)
		time.Sleep(time.Second * 3)
	}
	return fmt.Errorf("Weaviate did not start in time")
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

// TearDownLocalWeaviateDeprecated shuts down the locally started weaviate docker compose
// used for pre-v1.14 backwards compatibility tests.
// If `EXTERNAL_WEAVIATE_RUNNING` this function will not do anything
//
//	see SetupLocalWeaviate for more info.
func TearDownLocalWeaviateDeprecated() error {
	if isExternalWeaviateRunning() {
		return nil
	}
	err := test.TearDownWeaviateDeprecated()
	time.Sleep(time.Second * 3) // Add some delay to make sure the command was executed before the program exits
	return err
}
