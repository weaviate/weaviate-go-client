package main

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/batch"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/example"
)

func main() {
	ctx := context.Background()
	host, apiKey := example.ConnectionParams()

	// Connect to a WCD cluster.
	c, err := weaviate.NewWeaviateCloud(ctx, host, apiKey)
	example.Catch(err)
	defer c.Close()

	h, err := c.Collections.Create(ctx, collections.Collection{
		Name: "Vollast",
	})
	defer c.Collections.Delete(ctx, "Vollast")
	example.Catch(err)

	log.Printf("Created collection %q", h.CollectionName())

	log.Print("Start the batch...")
	b, err := h.Batch(ctx, batch.WithRetryTimes(1))
	example.Catch(err)

	tasks := make([]*batch.Task, 1000)
	log.Printf("Insert %d objects in %q", len(tasks), h.CollectionName())
	for i := range tasks {
		t, err := b.Object(nil)
		if err == context.Canceled {
			log.Println("Context canceled after %d objects, exit earlier", i)
			break
		}
		tasks[i] = t
	}

	log.Print("All objects added to the stream, closing the batch...")
	err = b.Close()
	example.Catch(err)

	log.Print("Streaming done, collect all results")
	for _, t := range tasks {
		if err := t.Wait(); err != nil {
			log.Printf("%s failed: %v", t.ID(), err)
		}
	}

	n, err := h.Count(ctx)
	example.Catch(err)

	log.Printf("Collection %q has %d objects", h.CollectionName(), n)
}
