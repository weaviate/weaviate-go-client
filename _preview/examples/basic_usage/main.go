package main

import (
	"context"
	"fmt"
	"log"

	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Year   int    `json:"year"`
	Genre  string `json:"genre"`
}

func main() {
	ctx := context.Background()

	// Create client
	client, _ := weaviate.NewLocal(ctx)

	// Get collection handle
	songs := client.Collections.Use("Songs")

	// Example 1: Insert objects
	fmt.Println("=== Example 1: Insert objects ===")
	insertObjects(ctx, songs)

	// Example 2: Query with map-based results
	fmt.Println("\n=== Example 2: Query with map-based results ===")
	queryWithMaps(ctx, songs)

	// Example 3: Query with typed results
	fmt.Println("\n=== Example 3: Query with typed results ===")
	queryWithTypes(ctx, songs)
}

func insertObjects(ctx context.Context, songs *collections.Handle) {
	// Insert with struct and vector
	song1 := Song{
		Title:  "Bohemian Rhapsody",
		Artist: "Queen",
		Year:   1975,
		Genre:  "Rock",
	}

	m, err := data.Encode(&song1)
	log.Print(m)
	if err != nil {
		log.Fatal(err)
	}
	obj, err := songs.Data.Insert(ctx, &data.Object{
		Properties: m,
		Vectors:    []types.Vector{{Single: []float32{0.1, 0.2, 0.3}}},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted with ID: %s\n", obj.UUID)

	// Insert with map, no vector
	song2 := map[string]any{
		"title":  "Song 2",
		"artist": "Blur",
		"year":   1997,
		"genre":  "Rock",
	}

	obj, _ = songs.Data.Insert(ctx, &data.Object{Properties: song2})
	fmt.Printf("Inserted with ID: %s\n", obj.UUID)
}

func queryWithMaps(ctx context.Context, songs *collections.Handle) {
	queryVector := types.Vector{
		Single: []float32{0.1, 0.2, 0.3, 0.4},
	}

	result, _ := songs.Query.NearVector(ctx, queryVector, query.NearVector{
		Limit:      2,
		Similarity: query.Distance(0.5),
		Offset:     3,
	})

	for i, obj := range result.Objects {
		number := obj.Properties["number"].(string)
		fmt.Printf("%d. Number: %s (UUID: %s)\n", i+1, number, obj.UUID)

		if vec, ok := obj.Vectors["default"]; ok {
			fmt.Printf("   Vector: %v\n", vec)
		}
	}
}

func queryWithTypes(ctx context.Context, songs *collections.Handle) {
	queryVector := types.Vector{
		Single: []float32{0.1, 0.2, 0.3, 0.4},
	}

	// Get results as maps first
	result, _ := songs.Query.NearVector(ctx, queryVector, query.NearVector{
		Limit: 3,
	})

	// Demonstrates type-safe scanning (Song struct would need to match actual data)
	typedObjects, _ := query.Decode[Song](result)

	for i, song := range typedObjects {
		fmt.Printf("%d. Title: %s, Artist: %s, Year: %d, Genre: %s (UUID: %s)\n",
			i+1, song.Properties.Title, song.Properties.Artist,
			song.Properties.Year, song.Properties.Genre, song.UUID)
	}
}
