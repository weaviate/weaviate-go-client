package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/example"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/query/filter"
)

func main() {
	ctx := context.Background()
	host, apiKey := example.ConnectionParams()

	// Connect to a WCD cluster.
	c, err := weaviate.NewWeaviateCloud(ctx, host, apiKey)
	example.Catch(err)
	defer c.Close()

	CollectionName := "Notes"

	c.Collections.Delete(ctx, CollectionName)
	notes, err := c.Collections.Create(ctx, collections.Collection{
		Name: CollectionName,
		Properties: []collections.Property{
			{Name: "tags", DataType: collections.DataTypeTextArray},
			{Name: "flags", DataType: collections.DataTypeBoolArray},
			{Name: "digits", DataType: collections.DataTypeIntArray},
			{Name: "pies", DataType: collections.DataTypeNumberArray},
			{Name: "birthdays", DataType: collections.DataTypeDateArray},
			{Name: "friends", DataType: collections.DataTypeUUIDArray},
			{
				Name:     "things",
				DataType: collections.DataTypeObjectArray,
				NestedProperties: []collections.Property{
					{Name: "chair", DataType: collections.DataTypeText},
					{Name: "lucky_number", DataType: collections.DataTypeInt},
				},
			},
		},
	})
	example.Catch(err)

	type Entry struct {
		Tags      []string
		Flags     []bool
		Digits    []int
		Pies      []float64
		Birthdays []time.Time
		Friends   []uuid.UUID
		Things    []map[string]any
	}

	bdayAlice := time.Date(1973, 8, 1, 4, 20, 0, 0, time.UTC)
	bdayBob := time.Date(1971, 11, 17, 16, 20, 0, 0, time.UTC)
	uuidCindy := uuid.New()
	uuidDave := uuid.New()

	entry := Entry{
		Tags:      []string{"#rock", "#blues"},
		Flags:     []bool{true, true, false, true},
		Digits:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		Pies:      []float64{3.14, 3.1415},
		Birthdays: []time.Time{bdayAlice, bdayBob},
		Friends:   []uuid.UUID{uuidCindy, uuidDave},
		Things: []map[string]any{
			{"chair": "pink"},
			{"lucky_number": int64(13)},
		},
	}

	in, err := notes.Data.Insert(ctx, &data.Object{
		Properties: data.MustEncode(&entry),
	})
	example.Catch(err)

	out, err := notes.Query.OverAll(ctx, query.OverAll{
		Filter: filter.Cond{
			Target:   filter.UUID,
			Operator: filter.Equal,
			Value:    in.UUIDs[0],
		},
	})
	example.Catch(err)

	var results []query.Object[Entry]
	example.Catch(query.Decode(out, &results))

	example.Assert(len(results) == 1, "expect 1 object back")
	got := results[0].Properties

	example.Assert(len(got.Tags) == len(entry.Tags), "len(Tags)")
	for i := range got.Tags {
		example.Assert(got.Tags[i] == entry.Tags[i], got.Tags[i], " != ", entry.Tags[i], " (Tags)")
	}

	example.Assert(len(got.Flags) == len(entry.Flags), "len(Flags)")
	for i := range got.Flags {
		example.Assert(got.Flags[i] == entry.Flags[i], got.Flags[i], " != ", entry.Flags[i], " (Flags)")
	}

	example.Assert(len(got.Digits) == len(entry.Digits), "len(Digits)")
	for i := range got.Digits {
		example.Assert(got.Digits[i] == entry.Digits[i], got.Digits[i], " != ", entry.Digits[i], " (Digits)")
	}

	example.Assert(len(got.Pies) == len(entry.Pies), "len(Pies)")
	for i := range got.Pies {
		example.Assert(got.Pies[i] == entry.Pies[i], got.Pies[i], " != ", entry.Pies[i], " (Pies)")
	}

	example.Assert(len(got.Birthdays) == len(entry.Birthdays), "len(Birthdays)")
	for i := range got.Birthdays {
		example.Assert(got.Birthdays[i].Equal(entry.Birthdays[i]), got.Birthdays[i], " != ", entry.Birthdays[i], " (Birthdays)")
	}

	example.Assert(len(got.Friends) == len(entry.Friends), "len(Friends)")
	for i := range got.Friends {
		example.Assert(got.Friends[i] == entry.Friends[i], got.Friends[i], " != ", entry.Friends[i], " (Friends)")
	}

	example.Assert(len(got.Things) == len(entry.Things), "len(Things)")
	for i := range got.Things {
		for k := range got.Things[i] {
			example.Assert(got.Things[i][k] == entry.Things[i][k], got.Things[i][k], " != ", entry.Things[i][k], " (Things)")
		}
	}
}
