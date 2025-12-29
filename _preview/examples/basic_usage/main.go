package main

import (
	"context"
	"fmt"

	weaviate "github.com/weaviate/weaviate-go-client/v6"
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
	client, _ := weaviate.NewClient("http", "localhost", 8080, 50051)

	client.Collections.Create(ctx, "Songs", collections.WithProperties(
		collections.Property{Name: "title", DataType: collections.DataTypeText},
		collections.Property{Name: "artist", DataType: collections.DataTypeText},
		collections.Property{Name: "year", DataType: collections.DataTypeInt},
		collections.Property{Name: "genre", DataType: collections.DataTypeText},
	))

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

	client.Collections.Delete(ctx, "Songs")
}

func insertObjects(ctx context.Context, songs *collections.Handle) {
	// Insert with struct and vector
	song1 := Song{
		Title:  "Bohemian Rhapsody",
		Artist: "Queen",
		Year:   1975,
		Genre:  "Rock",
	}

	id, _ := songs.Data.Insert(ctx,
		data.WithProperties(song1),
		data.WithVector(types.Vector{Single: []float32{0.1, 0.2, 0.3}}),
	)
	fmt.Printf("Inserted with ID: %s\n", id)

	// Insert with map, no vector
	song2 := map[string]any{
		"title":  "Song 2",
		"artist": "Blur",
		"year":   1997,
		"genre":  "Rock",
	}

	id, _ := songs.Data.Insert(ctx,
		data.WithProperties(song2),
	)
	fmt.Printf("Inserted with ID: %s\n", id)
}

func queryWithMaps(ctx context.Context, songs *collections.Handle) {
	queryVector := types.Vector{
		Single: []float32{0.1, 0.2, 0.3, 0.4},
	}

	result, _ := songs.Query.NearVector(ctx, queryVector,
		query.WithLimit(2),
		query.WithDistance(0.5),
		query.WithOffset(3),
	)

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
	result, _ := songs.Query.NearVector(ctx, queryVector,
		query.WithLimit(3),
	)

	// Demonstrates type-safe scanning (Song struct would need to match actual data)
	typedObjects := query.Scan[Song](result)

	for i, song := range typedObjects {
		fmt.Printf("%d. Title: %s, Artist: %s, Year: %d, Genre: %s (UUID: %s)\n",
			i+1, song.Properties.Title, song.Properties.Artist,
			song.Properties.Year, song.Properties.Genre, song.UUID)
	}

	grouped, _ := songs.Query.NearVector.GroupBy(ctx,
		types.Vector{Single: single},
		"group by album",
		query.WithAutoLimit(2),
	)

	albums := query.ScanGrouped[Song](grouped)
	for _, album := range albums.Groups {
		fmt.Printf("album %q has %d songs:", album.Name, album.Size)
		for _, song := range album.Objects {
			fmt.Printf("\t%q - %s\n",
				song.Properties.Title, song.Properties.Lyrics)
		}
	}
}
