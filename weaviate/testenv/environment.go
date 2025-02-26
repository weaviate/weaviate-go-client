package testenv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
)

// EXTERNAL_WEAVIATE_RUNNING is the environment variable which controls the lifecycle
// of the Docker containers created for the test. If set to 'true', containers are
// preserved between test runs (useful for CI scenario, where all existing presets
// are started at once with ../../test/start_containers.sh.
// Otherwise containers are torn down on every test cleanup.
var EXTERNAL_WEAVIATE_RUNNING = os.Getenv("EXTERNAL_WEAVIATE_RUNNING")

var envExternalWeaviateRunning bool = strings.ToLower(EXTERNAL_WEAVIATE_RUNNING) == "true"

func waitForStartup(ctx context.Context, url string, interval time.Duration) error {
	t := time.NewTicker(interval)
	defer t.Stop()
	expired := ctx.Done()
	var lastErr error
	for {
		select {
		case <-t.C:
			lastErr = checkReady(ctx, url)
			if lastErr == nil {
				return nil
			}
		case <-expired:
			return fmt.Errorf("init context expired before remote was ready: %w", lastErr)
		}
	}
}

func checkReady(initCtx context.Context, url string) error {
	// spawn a new context (derived on the overall context) which is used to
	// consider an individual request timed out
	requestCtx, cancel := context.WithTimeout(initCtx, 500*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet,
		fmt.Sprintf("http://%s/v1/.well-known/ready", url), nil)
	if err != nil {
		return fmt.Errorf("create check ready request: %w", err)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send check ready request: %w", err)
	}

	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("not ready: status %d", res.StatusCode)
	}

	return nil
}

// SetupLocalWeaviate creates a local weaviate running on 8080 using docker compose
// Will only wait for it to be reachable if env `EXTERNAL_WEAVIATE_RUNNING` is set to True.
//
//	`EXTERNAL_WEAVIATE_RUNNING` should be set if all tests are supposed to be run in a test suit.
//	This prevents unnecessary starting and stopping of the docker-compose which prevents errors
//	due to syncing issues and speeds up the process
func SetupLocalWeaviateWaitForStartup(waitForWeaviateStartup bool) error {
	if !isExternalWeaviateRunning() {
		port, _, authEnabled := testsuit.GetPortAndAuthPw()
		if err := test.SetupWeaviate(authEnabled); err != nil {
			return err
		}
		if waitForWeaviateStartup {
			return waitForStartup(context.TODO(), fmt.Sprintf("localhost:%v", port), 1*time.Second)
		}
		return nil
	}
	return nil
}

func SetupLocalWeaviate() error {
	return SetupLocalWeaviateWaitForStartup(true)
}

// SetupLocalContainer is a test helper that starts a Weaviate instance from
// a docker-compose files specified by the preset. It returns a cleanup function,
// which will tear down ALL current containers if 'EXTERNAL_WEAVIATE_RUNNING=true'
// is set and do nothing otherwise.
//
// Usage:
//
//	container, stop := testenv.SetupLocalContainer(t, ctx, test.Basic, true)
//	t.Cleanup(stop)
//	client := testsuit.CreateTestClientForContainer(t, container)
func SetupLocalContainer(t *testing.T, ctx context.Context, preset test.Preset, waitForWeaviateStartup bool) (test.Container, func()) {
	t.Helper()

	container, start, stop := test.GetContainer(preset)

	err := start()
	if err == nil && waitForWeaviateStartup {
		err = waitForStartup(ctx, container.HTTPAddress(), 1*time.Second)
	}
	require.NoErrorf(t, err, "start container from %q", container.DockerComposeFile)

	mustStop := func() {
		if !envExternalWeaviateRunning {
			require.NoErrorf(t, stop(), "stop container from %q", container.DockerComposeFile)
		}
	}
	return container, mustStop
}

// isExternalWeaviateRunning checks if either EXTERNAL_WEAVIATE_RUNNING is set
// or a Weaviate container is already running.
func isExternalWeaviateRunning() bool {
	port, _, _ := testsuit.GetPortAndAuthPw()
	err := checkReady(context.TODO(), fmt.Sprintf("localhost:%v", port))
	isRunning := err == nil

	return envExternalWeaviateRunning || isRunning
}

// TearDownLocalWeaviate shuts down the locally started weaviate docker compose
// If `EXTERNAL_WEAVIATE_RUNNING` this function will not do anything
//
//	see SetupLocalWeaviate for more info.
func TearDownLocalWeaviate() error {
	if isExternalWeaviateRunning() {
		return nil
	}
	return TearDownLocalWeaviateForcefully()
}

func TearDownLocalWeaviateForcefully() error {
	err := test.TearDownWeaviate()
	time.Sleep(time.Second * 3) // Add some delay to make sure the command was executed before the program exits
	return err
}
