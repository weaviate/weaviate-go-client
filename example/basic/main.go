package main

import (
	"context"
	"log"
	"os"

	"github.com/weaviate/weaviate-go-client/v6"
)

const (
	EnvHost   = "WEAVIATE_HOST"
	EnvAPIKey = "WEAVIATE_API_KEY"
)

func main() {
	host, hasHost := os.LookupEnv(EnvHost)
	apiKey, hasKey := os.LookupEnv(EnvAPIKey)
	if !hasHost || !hasKey {
		log.Printf("%q and %q must be defined. Skipping example.", EnvHost, EnvAPIKey)
		return
	}

	ctx := context.Background()

	c, err := weaviate.NewWeaviateCloud(ctx, host, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	if ok, err := c.Collections.Exists(ctx, "SampleProducts"); err != nil {
		log.Fatal(err)
	} else if !ok {
		log.Print("SampleProducts collection does not exist. Skipping example.")
		return
	}

	products := c.Collections.Use("SampleProducts")

	count, err := products.Count(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("collection SampleProducts has %d objects", count)

	for i := range 5 {
		obj, err := products.Data.Insert(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("object #%d=%q", i, obj.UUID)
	}

	count, err = products.Count(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("collection SampleProducts has %d objects", count)
}
