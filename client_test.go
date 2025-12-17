package weaviate_test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/backup"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestClient(t *testing.T) {
	ctx := t.Context()

	c, err := weaviate.NewLocal(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// The disadvantage of the pattern below is that you can re-arrange the two
	// parameters, in which case WithLocal would overwrite the overrides.
	// 	weaviate.NewClient(ctx, weaviate.WithLocal(), weaviate.WithHTTPPort(8081))

	// Create a collection and get a handle
	h, err := c.Collections.Create(ctx, "Songs")
	if err != nil {
		t.Fatal(err)
	}

	// YOLO: write with ConsistencyLevel=ONE
	h = h.WithOptions(collections.WithConsistencyLevel(types.ConsistencyLevelOne))

	var single []float32
	var multi [][]float32

	h.Data.Insert(ctx, data.WithProperties(map[string]any{
		"album": "Killin Is My Business... And Business Is Good!",
		"title": "Rattlehead",
	}), data.WithVector(
		types.Vector{Name: "title", Single: single},
		types.Vector{Name: "lyrics", Multi: multi},
	), data.WithUUID(uuid.New()))

	h.Query.NearVector(ctx,
		types.Vector{Single: single},
		query.WithLimit(5),
		query.WithDistance(.34),
	)

	h.Query.NearVector(ctx, query.Average(
		types.Vector{Name: "title", Single: single},
		types.Vector{Name: "lyrics", Multi: multi},
	))

	res, _ := h.Query.NearVector(ctx, query.ManualWeights(
		query.Target(types.Vector{Name: "title", Single: single}, 0.4),
		query.Target(types.Vector{Name: "lyrics", Multi: multi}, 0.6),
	))

	fmt.Printf("NearVector: got %d objects\n", len(res.Objects))
	for i, obj := range res.Objects {
		fmt.Printf("#%d: id=%s\n\tproperties=%v\n", i, obj.UUID, obj.Properties)
	}

	grouped, _ := h.Query.NearVector.GroupBy(ctx,
		types.Vector{Single: single},
		"group by album",
		query.WithAutoLimit(2),
	)

	fmt.Printf("NearVector.GroupBy: got %d objects in %d groups\n", len(grouped.Objects), len(grouped.Groups))
	for key, group := range grouped.Groups {
		fmt.Printf("Group %q has %d objects \n", key, group.Size)
	}

	// Scan results differently
	type Song struct{ Title, Album, Lyrics string }

	songs := query.Scan[Song](res)

	for _, song := range songs {
		fmt.Printf("%q (%s) - %s\n",
			song.Properties.Title, song.Properties.Album, song.Properties.Lyrics)
	}

	albums, _ := query.ScanGrouped[Song](grouped)
	for _, album := range albums {
		fmt.Printf("album %q has %d songs:", album.Name, album.Size)
		for _, song := range album.Objects {
			fmt.Printf("\t%q - %s\n",
				song.Properties.Title, song.Properties.Lyrics)
		}
	}

	// Backups
	bak, err := c.Backup.Create(ctx, "bak-1", "filesystem",
		backup.WithIncludeCollections("Songs", "Artists"))
	if err != nil {
		log.Fatal(err)
	}

	backup.AwaitCompletion(ctx, bak, backup.WithPollingInterval(time.Minute))

	c.Backup.Restore(ctx, bak.ID, bak.Backend, backup.RestoreOptions{
		IncludeCollections: []string{"Songs", "Artists"},
	})
}
