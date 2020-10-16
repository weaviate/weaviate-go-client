package testenv

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/test"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient"
	"time"
)

func SetupLocalWeaviate() error {
	err := test.SetupWeavaite()
	if err != nil {
		return err
	}
	return WaitForWeaviate()
}

func WaitForWeaviate() error {
	cfg := weaviateclient.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	client := weaviateclient.New(cfg)


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

func TearDownLocalWeaviate() error {
	return test.TearDownWeavaite()
}
