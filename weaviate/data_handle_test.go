package weaviate

import (
	"testing"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func TestDataClient_Insert_Simple(t *testing.T) {
	coll := &CollectionClient{name: "TestClass"}
	client := &DataClient{collection: coll}

	id, err :=
		client.Insert(
			t.Context(),
			WithProperties(map[string]any{"foo": "bar"}),
			WithID(uuid.NewString()),
			WithValidation(),
			WithVector(types.Vector{Single: []float32{0.1, 0.2, 0.3}}),
			WithVector(types.Vectors{"single": {Single: []float32{0.1, 0.2, 0.3}}}),
			WithVector([]types.Vector{{Name: "name", Single: []float32{0.1, 0.2, 0.3}}}),
		)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id == uuid.Nil {
		t.Fatalf("expected valid uuid, got nil")
	}
}
