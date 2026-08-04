package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/aggregate"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/example"
	vv8 "github.com/weaviate/weaviate-go-client/v6/modules/weaviate"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

func main() {
	ctx := context.Background()
	host, apiKey := example.ConnectionParams()

	// Connect to a WCD cluster.
	c, err := weaviate.NewWeaviateCloud(ctx, host, apiKey)
	example.Catch(err)
	defer c.Close()

	CollectionName := "Songs"

	c.Collections.Delete(ctx, CollectionName)
	songs, err := c.Collections.Create(ctx, collections.Collection{
		Name: CollectionName,
		Properties: []collections.Property{
			{Name: "title", DataType: collections.DataTypeText},
			{Name: "album", DataType: collections.DataTypeText},
			{Name: "duration_sec", DataType: collections.DataTypeInt},
		},
		Vectors: map[string]collections.VectorConfig{
			"title_vec": {
				Vectorizer: vv8.Text2Vec{
					Properties: []string{"title"},
					Model:      vv8.SnowflakeArcticEmbedMv1_5,
					Dimensions: 256,
				},
			},
		},
	})
	example.Catch(err)

	// Insert some objects, logging the "before" and "after" counts.
	count, err := songs.Count(ctx)
	example.Catch(err)

	log.Printf("collection %s has %d objects", CollectionName, count)

	type Song struct {
		Title    string `json:"title"`
		Album    string `json:"album"`
		Duration int    `json:"duration_sec"`
	}

	for _, s := range []Song{
		{
			Title:    "Schism",
			Album:    "Lateralus",
			Duration: 406,
		},
		{
			Title:    "Parabola",
			Album:    "Lateralus",
			Duration: 363,
		},
		{
			Title:    "The Pot",
			Album:    "10,000 Days",
			Duration: 382,
		},
		{
			Title:    "Forty Six & 2",
			Album:    "Ænima",
			Duration: 376,
		},
	} {
		res, err := songs.Data.Insert(ctx, &data.Object{
			Properties: data.MustEncode(&s),
		})
		if err != nil {
			log.Fatal(err)
		}
		for id, msg := range res.Errors {
			fmt.Printf("\tinsert song %q (id=%s) failed with error %q\n", s.Title, id, msg)
		}
	}

	count, err = songs.Count(ctx)
	example.Catch(err)

	log.Printf("collection %s has %d objects", CollectionName, count)

	// Query some objects using NearText search.
	nt, err := songs.Query.NearText(ctx, query.NearText{
		Concepts:      []string{"forty", "six"},
		MoveAway:      &query.Move{Concepts: []string{"pot"}, Force: .34},
		Limit:         3,
		ReturnVectors: []string{"title_vec"},
	})
	example.Catch(err)

	log.Printf("NearText[forty, six] returned %d objects:", len(nt.Objects))
	for _, obj := range nt.Objects {
		fmt.Println("\t- ", obj.Properties["title"])
	}

	if len(nt.Objects) == 0 {
		return
	}

	// Fetch 3 most similar objects to the first result hit
	target := nt.Objects[0]
	nv, err := songs.Query.NearVector(ctx, query.NearVector{
		Target:           target.Vectors["title_vec"],
		Similarity:       query.Distance(0.56),
		AutoLimit:        2,
		Limit:            3,
		ReturnProperties: []string{"title", "album"},
		ReturnMetadata:   query.ReturnMetadata{Distance: true},
	})
	example.Catch(err)

	decoded := make([]query.Object[Song], len(nv.Objects))
	example.Catch(query.Decode(nv, &decoded))

	log.Printf("NearVector[max_distance=.56] returned these %d entries:", len(nv.Objects))
	for _, obj := range decoded {
		fmt.Printf("\t- %q distance=%f\n", obj.Properties.Title, *obj.Metadata.Distance)
	}

	grouped, err := songs.Aggregate.OverAll.GroupBy(ctx, aggregate.OverAll{
		Integer: []aggregate.Integer{
			{Property: "duration_sec", Min: true, Median: true, Max: true},
		},
	}, aggregate.GroupBy{Property: "album", Limit: 5})
	example.Catch(err)

	for _, group := range grouped.Groups {
		log.Printf("Group %q has %d objects (group_by=%q)", group.Value, len(group.Aggregations.Integer), group.Property)
		for name, duration := range group.Aggregations.Integer {
			fmt.Printf("\t- songs in %q have median duration of %v\n", name, time.Duration(*duration.Median)*time.Second)
		}
	}
}
