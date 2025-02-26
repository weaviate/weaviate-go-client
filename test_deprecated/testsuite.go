package test_deprecated

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

func CreateTestClient() *weaviate.Client {
	cfg := weaviate.Config{
		Host:   "localhost:8089",
		Scheme: "http",
	}
	client := weaviate.New(cfg)
	client.WaitForWeavaite(60 * time.Second)
	return client
}

func command(app string, arguments []string) error {
	mydir, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command(app, arguments...)
	cmd.Dir = mydir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}

func isExternalWeaviateRunning() bool {
	val := os.Getenv("EXTERNAL_WEAVIATE_RUNNING")
	val = strings.ToLower(val)
	fmt.Printf("\nEXTERNAL_WEAVIATE_RUNNING: %v\n", val)
	return val == "true"
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
	err := TearDownWeaviateDeprecated()
	time.Sleep(time.Second * 3) // Add some delay to make sure the command was executed before the program exits
	return err
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
		return SetupWeaviateDeprecated()
	}
	return nil
}

// SetupWeaviateDeprecated run docker compose up
// for pre-v1.14 backwards compatibility tests
func SetupWeaviateDeprecated() error {
	app := "docker"
	arguments := []string{
		"compose",
		"-f",
		"docker-compose-deprecated-api-test.yml",
		"up",
		"-d",
	}
	return command(app, arguments)
}

// TearDownWeaviateDeprecated run docker-compose down
// for pre-v1.14 backwards compatibility tests
func TearDownWeaviateDeprecated() error {
	app := "docker"
	arguments := []string{
		"compose",
		"-f",
		"docker-compose-deprecated-api-test.yml",
		"down",
		"--remove-orphans",
	}
	return command(app, arguments)
}
