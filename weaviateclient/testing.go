package weaviateclient

import (
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/test"
	"time"
	"context"
)

func setupLocalWeaviate() error {
	err := test.SetupWeavaite()
	if err != nil {
		return err
	}
	return waitForWeaviate()
}

func waitForWeaviate() error {
	cfg := Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client := New(cfg)


	for i:=0;i<20;i++ {
		ctx, _ := context.WithTimeout(context.Background(), time.Second * 3)
		isReady, _ := client.Misc.ReadyChecker().Do(ctx)
		if isReady {
			return nil
		}
		fmt.Printf("Weaviate not yet up waiting another 3 seconds. Iteration: %v", i)
		time.Sleep(time.Second * 3)
	}
	return fmt.Errorf("Weaviate did not start in time")
}

func tearDownLocalWeaviate() error {
	return test.TearDownWeavaite()
}