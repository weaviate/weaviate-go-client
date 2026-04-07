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

	c, err := weaviate.NewWeaviateCloud(context.Background(), host, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	products := c.Collections.Use("SampleProducts")

	obj, err := products.Data.Insert(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Inserted object %q", obj.UUID)
}
