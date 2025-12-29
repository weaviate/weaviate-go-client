package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

type Song struct {
	Title  string   `json:"title"`
	Artist string   `json:"artist"`
	Year   int      `json:"year"`
	Genre  []string `json:"genre"`
}

func main() {
	ctx := context.Background()

	// Create client
	client, _ := weaviate.NewLocal(ctx, nil)

	if err := client.Collections.DeleteAll(ctx); err != nil {
		log.Fatal(err)
	}

	// Get collection handle
	songs, err := client.Collections.Create(ctx, collections.Collection{
		Name: "Songs",
		Properties: []collections.Property{
			{Name: "title", DataType: collections.DataTypeText},
			{Name: "artist", DataType: collections.DataTypeText},
			{Name: "year", DataType: collections.DataTypeInt},
			{Name: "genre", DataType: collections.DataTypeTextArray},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Insert objects
	fmt.Println("=== Example 1: Insert objects ===")
	if err := insertObjects(ctx, songs); err != nil {
		log.Fatal(err)
	}

	// Example 2: Query with map-based results
	fmt.Println("\n=== Example 2: Query with map-based results ===")
	if err := queryWithMaps(ctx, songs); err != nil {
		log.Fatal(err)
	}

	// Example 3: Query with typed results
	fmt.Println("\n=== Example 3: Query with typed results ===")
	if err := queryWithTypes(ctx, songs); err != nil {
		log.Fatal(err)
	}
}

func insertObjects(ctx context.Context, songs *collections.Handle) error {
	// Insert with struct and vector
	song1 := Song{
		Title:  "Bohemian Rhapsody",
		Artist: "Queen",
		Year:   1975,
		Genre:  []string{"Rock", "Opera"},
	}

	obj, err := songs.Data.Insert(ctx, &data.Object{
		Properties: data.MustEncode(&song1),
		// This synax works, but we can't provide our own vectors
		// because the collection only has a default vectorizer.
		// Vectors: []types.Vector{
		// 	{Single: []float32{0.1, 0.2, 0.3}},
		// },
	})
	if err != nil {
		return err
	}
	fmt.Printf("Inserted with ID: %s\n", obj.UUID)

	// Insert with map, no vector
	song2 := map[string]any{
		"title":  "Song 2",
		"artist": "Blur",
		"year":   1997,
		"genre":  []string{"Rock"},
	}

	obj, err = songs.Data.Insert(ctx, &data.Object{Properties: song2})
	if err != nil {
		return err
	}
	fmt.Printf("Inserted with ID: %s\n", obj.UUID)
	return nil
}

func queryWithMaps(ctx context.Context, songs *collections.Handle) error {
	// model2vec-text2vec (default vectorizer) generates vectors with 128 dimensions.
	v := make([]float32, 128)
	for i := range v {
		v[i] = rand.Float32()
	}

	result, err := songs.Query.NearVector(ctx, query.NearVector{
		Target:     &types.Vector{Single: v},
		Limit:      2,
		Similarity: query.Distance(0.5),
		Offset:     3,
	})
	if err != nil {
		return err
	}

	for i, obj := range result.Objects {
		number := obj.Properties["number"].(string)
		fmt.Printf("%d. Number: %s (UUID: %s)\n", i+1, number, obj.UUID)

		if vec, ok := obj.Vectors["default"]; ok {
			fmt.Printf("   Vector: %v\n", vec)
		}
	}
	return nil
}

func queryWithTypes(ctx context.Context, songs *collections.Handle) error {
	// model2vec-text2vec (default vectorizer) generates vectors with 128 dimensions.
	v := make([]float32, 128)
	for i := range v {
		v[i] = rand.Float32()
	}
	// Get results as maps first

	result, err := songs.Query.NearVector(ctx, query.NearVector{
		Target: &types.Vector{Single: v},
		Limit:  3,
	})
	if err != nil {
		return err
	}

	// Demonstrates type-safe scanning (Song struct would need to match actual data)
	typedObjects, err := query.Decode[Song](result)
	if err != nil {
		return err
	}

	for i, song := range typedObjects {
		fmt.Printf("%d. Title: %s, Artist: %s, Year: %d, Genre: %s (UUID: %s)\n",
			i+1, song.Properties.Title, song.Properties.Artist,
			song.Properties.Year, song.Properties.Genre, song.UUID)
	}
	return nil
}
