package graphql

import (
	"context"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
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
		Class: "MultiTargetSearchTest",
		VectorConfig: map[string]models.VectorConfig{
			"first":  {Vectorizer: map[string]interface{}{"none": nil}, VectorIndexType: "hnsw"},
			"second": {Vectorizer: map[string]interface{}{"none": nil}, VectorIndexType: "hnsw"},
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

	outer := []struct {
		name string
		nvo  *graphql.NearMultiVectorArgumentBuilder
	}{
		{name: "Sum", nvo: client.GraphQL().NearVectorMultiTargetArgBuilder().Sum("first", "second")},
		{name: "Average", nvo: client.GraphQL().NearVectorMultiTargetArgBuilder().Average("first", "second")},
		{name: "Minimum", nvo: client.GraphQL().NearVectorMultiTargetArgBuilder().Minimum("first", "second")},
		{name: "Manual weights", nvo: client.GraphQL().NearVectorMultiTargetArgBuilder().ManualWeights(map[string]float32{"first": 1, "second": 1})},
		{name: "Relative score", nvo: client.GraphQL().NearVectorMultiTargetArgBuilder().RelativeScore(map[string]float32{"first": 1, "second": 1})},
		{name: "No", nvo: client.GraphQL().NearVectorMultiTargetArgBuilder()},
	}
	for _, to := range outer {
		inner := []struct {
			name string
			nvi  *graphql.NearMultiVectorArgumentBuilder
		}{
			{name: "with vector", nvi: to.nvo.WithVector([]float32{1, 0, 0})},
			{name: "with vector per target", nvi: to.nvo.WithVectorPerTarget(map[string][]float32{"first": {1, 0}, "second": {1, 0, 0}})},
		}
		for _, ti := range inner {
			t.Run(to.name+" combination "+ti.name, func(t *testing.T) {
				resp, err := client.GraphQL().Get().WithNearMultiVector(ti.nvi).WithClassName(class.Class).WithFields(graphql.Field{Name: "_additional", Fields: []graphql.Field{{Name: "id"}}}).Do(ctx)
				require.Nil(t, err)
				if resp.Errors != nil {
					errors := make([]string, len(resp.Errors))
					for i, e := range resp.Errors {
						errors[i] = e.Message
					}
					t.Fatalf("errors: %v", strings.Join(errors, ", "))
				}
				require.NotNil(t, resp.Data)
				require.Equal(t, objs[0].ID.String(), resp.Data["Get"].(map[string]interface{})[class.Class].([]interface{})[0].(map[string]interface{})["_additional"].(map[string]interface{})["id"].(string))
			})
		}
	}
}
