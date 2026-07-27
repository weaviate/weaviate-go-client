package main

import (
	"context"
	"errors"
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

	log.Print("Queue https://www.youtube.com/watch?v=WmPGZuaOlq4...")

	c.Collections.Delete(ctx, "Vollast")
	h, err := c.Collections.Create(ctx, collections.Collection{
		Name: "Vollast",
	})
	defer c.Collections.Delete(ctx, "Vollast")
	example.Catch(err)

	log.Printf("Created collection %q", h.CollectionName())

	log.Print("Start the batch...")
	batchCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b, err := h.Batch(batchCtx, batch.WithRetryTimes(1))
	example.Catch(err)

	tasks := make([]*batch.Task, 0, 1093)
	log.Printf("Insert %d objects in %q", cap(tasks), h.CollectionName())
	for i := range cap(tasks) {
		if i == 1024 {
			t := tasks[i-1]
			log.Printf("Wait for task=%s...", t.ID())
			example.Catch(t.Wait())
			log.Print("... and cancel context hehehe")
			log.Printf(`Guess we won't see those %d objects ¯\_(ツ)_/¯`, cap(tasks)-i)
			cancel()
		}
		t, err := b.Object(nil)
		if err == context.Canceled {
			log.Printf("Context canceled after %d objects, exit earlier", i-1)
			break
		}
		tasks = append(tasks, t)
	}

	log.Printf("%d objects added to the stream, closing the batch...", len(tasks))
	err = b.Close()
	if err != nil && !errors.Is(err, context.Canceled) {
		example.Catch(err)
	}

	log.Printf("Streaming done, collect all results (expect %d to fail with context.Canceled)", len(tasks)-1024)
	for _, t := range tasks {
		if err := t.Wait(); err != nil {
			log.Printf("\t- %s failed: %q", t.ID(), err)
		}
	}

	n, err := h.Count(ctx)
	example.Catch(err)

	log.Printf("Collection %q has %d objects", h.CollectionName(), n)
}
