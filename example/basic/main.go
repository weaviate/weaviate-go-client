package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/query"
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

	// Connect to a WCD cluster.
	c, err := weaviate.NewWeaviateCloud(ctx, host, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	// If SampleProducts collection does not exist, create it.
	canSearch := true
	if ok, err := c.Collections.Exists(ctx, "SampleProducts"); err != nil {
		log.Fatal(err)
	} else if !ok {
		c.Collections.Create(ctx, collections.Collection{
			Name: "SampleProducts",
			Properties: []collections.Property{
				{Name: "name", DataType: collections.DataTypeText},
				{Name: "description", DataType: collections.DataTypeText},
				{Name: "url", DataType: collections.DataTypeText},
			},
		})

		// Current client version does not support defining vector indices.
		// Without one, similarity search is not possible.
		canSearch = false
	}

	// Get a handle for SampleProducts collection.
	products := c.Collections.Use("SampleProducts")

	// Insert some objects, logging the "before" and "after" counts.
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
		fmt.Printf("\tobject #%d=%q\n", i, obj.UUID)
	}

	count, err = products.Count(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("collection SampleProducts has %d objects", count)

	if !canSearch {
		log.Print("collection SampleProducts has no vector indices")
		return
	}

	// Query some objects using NearText search.
	nt, err := products.Query.NearText(ctx, query.NearText{
		Concepts:      []string{"sneakers", "flipflops"},
		MoveAway:      &query.Move{Concepts: []string{"sandals"}, Force: .34},
		Limit:         3,
		ReturnVectors: []string{"text2vec_weaviate"}, // the default vector name
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("NearText[sneakers, flipflops] returned %d objects:", len(nt.Objects))
	for _, obj := range nt.Objects {
		fmt.Println("\t- ", obj.Properties["description"])
	}

	if len(nt.Objects) == 0 {
		return
	}

	// Fetch 3 most similar objects to the first result hit
	target := nt.Objects[0]
	nv, err := products.Query.NearVector(ctx, query.NearVector{
		Target:           target.Vectors["text2vec_weaviate"], // the default vector name
		Similarity:       query.Distance(0.56),
		AutoLimit:        2,
		Limit:            3,
		ReturnProperties: []string{"name", "url"},
		ReturnMetadata:   query.ReturnMetadata{Distance: true},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Scan results into our custom Go struct
	type Product struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoded := make([]query.Object[Product], len(nv.Objects))
	if err := query.Decode(nv, &decoded); err != nil {
		log.Fatal(err)
	}

	log.Print("NearVector[max_distance=.56] returned these 3 entries:")
	for _, obj := range decoded {
		fmt.Printf("\t- [%s](%s) distance=%f\n", obj.Properties.Name, obj.Properties.URL, *obj.Metadata.Distance)
	}
}
