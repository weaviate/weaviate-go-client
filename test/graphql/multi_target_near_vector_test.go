package graphql

import (
	"context"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestMultiTargetSearch(t *testing.T) {
	err := testenv.SetupLocalWeaviate()
	if err != nil {
		t.Fatalf("failed to setup weaviate: %s", err)
	}
	client := testsuit.CreateTestClient()
	ctx := context.TODO()

	class := &models.Class{
		Class: "MultiTargetTest",
		VectorConfig: map[string]models.VectorConfig{
			"first":  {Vectorizer: map[string]interface{}{"none": nil}},
			"second": {Vectorizer: map[string]interface{}{"none": nil}},
		},
	}
	require.Nil(t, client.Schema().ClassCreator().WithClass(class).Do(ctx))
	defer client.Schema().ClassDeleter().WithClassName(class.Class).Do(ctx)

	objs := []*models.Object{
		{Class: class.Class, Vectors: map[string]models.Vector{"first": []float32{1, 0}, "second": []float32{1, 0, 0}}, ID: strfmt.UUID(uuid.New().String())},
		{Class: class.Class, Vectors: map[string]models.Vector{"first": []float32{0, 1}, "second": []float32{0, 0, 1}}, ID: strfmt.UUID(uuid.New().String())},
	}

	_, err = client.Batch().ObjectsBatcher().WithObjects(objs...).Do(ctx)
	require.Nil(t, err)

	nv := client.GraphQL().NearVectorMultiTargetArgBuilder()

	resp, err := client.GraphQL().Get().WithNearVectorMultiTarget(nv.Sum("first", "second").WithVectorPerTarget(map[string][]float32{"first": {1, 0}, "second": {1, 0, 0}})).WithClassName(class.Class).Do(ctx)
	require.Nil(t, err)
	require.NotNil(t, resp)
}
