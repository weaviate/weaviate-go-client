package testenv

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/test"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
)

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
		if err := test.SetupWeaviate(); err != nil {
			return err
		}
		if waitForWeaviateStartup {
			port, _, _ := testsuit.GetPortAndAuthPw()
			return waitForStartup(context.TODO(), fmt.Sprintf("localhost:%v", port), 1*time.Second)
		}
		return nil
	}
	return nil
}

func SetupLocalWeaviate() error {
	return SetupLocalWeaviateWaitForStartup(true)
}

func isExternalWeaviateRunning() bool {
	val := os.Getenv("EXTERNAL_WEAVIATE_RUNNING")
	val = strings.ToLower(val)
	isExternalWeaviateRunning := val == "true"

	port, _, _ := testsuit.GetPortAndAuthPw()
	err := checkReady(context.TODO(), fmt.Sprintf("localhost:%v", port))
	isRunning := err == nil

	return isExternalWeaviateRunning || isRunning
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
