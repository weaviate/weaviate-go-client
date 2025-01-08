package colbert

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/backup"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

func TestJinjaAIVectorizer(t *testing.T) {
	ctx := context.Background()
	// err := testenv.SetupLocalWeaviate()
	// if err != nil {
	// 	require.Nil(t, err, "failed to start weaviate")
	// }

	client := testsuit.CreateTestClient(true)

	t.Run("searches with simple data ", func(t *testing.T) {

		// clean up DB
		err := client.Schema().AllDeleter().Do(context.Background())
		require.NoError(t, err)

		// create class
		className := "JinjaAIVectorizer"

		vectorConfig := map[string]models.VectorConfig{
			"regular": {
				Vectorizer: map[string]interface{}{
					"none": nil,
				},
				VectorIndexType: "hnsw",
			},
			"colbert": {
				Vectorizer: map[string]interface{}{
					"text2colbert-jinaai": map[string]interface{}{
						"properties":         []interface{}{"description"},
						"vectorizeClassName": false,
						"model":              "jina-colbert-v2",
					},
				},
				VectorIndexType: "hnsw",
			},
		}
		class := &models.Class{
			Class: className,
			Properties: []*models.Property{
				{
					Name:     "name",
					DataType: []string{schema.DataTypeText.String()},
				},
				{
					Name:     "description",
					DataType: []string{schema.DataTypeText.String()},
				},
			},
			VectorConfig: vectorConfig,
			ReplicationConfig: &models.ReplicationConfig{
				Factor: 1,
			},
		}
		err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
		require.NoError(t, err)
		// perform batch import
		objects := make([]*models.Object, 5)
		for i := range objects {
			objects[i] = &models.Object{
				ID:    generateUUID(uint32(i)),
				Class: className,
				Properties: map[string]interface{}{
					"name":        "some name",
					"description": "some description" + strconv.Itoa(i),
				},
			}
		}
		batchResponse, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(ctx)
		require.NoError(t, err)
		assert.NotNil(t, batchResponse)
		assert.Equal(t, len(objects), len(batchResponse))

		fields := []graphql.Field{
			{Name: "description"},
			{Name: "_additional", Fields: []graphql.Field{
				{Name: "id"},
				{Name: "distance"},
			}},
		}

		time.Sleep(10 * time.Second)

		t.Run("get objects and near vector", func(t *testing.T) {

			// get object
			for _, obj := range objects {
				objs, err := client.Data().ObjectsGetter().
					WithClassName(className).WithID(obj.ID.String()).WithVector().
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, objs)
				require.Len(t, objs, 1)
				assert.Equal(t, obj.ID, objs[0].ID)

				nearVector := client.GraphQL().NearVectorArgBuilder().
					WithVector(objs[0].Vectors["colbert"]).
					WithTargetVectors("colbert")

				response, err := client.GraphQL().Get().
					WithClassName(className).
					WithFields(fields...).
					WithNearVector(nearVector).
					WithLimit(2).
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, response)
				require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)

			}
		})

		t.Run("get near object", func(t *testing.T) {

			// Test near object search
			for _, obj := range objects {
				nearObject := client.GraphQL().NearObjectArgBuilder().WithID(obj.ID.String()).
					WithTargetVectors("colbert")
				response, err := client.GraphQL().Get().
					WithClassName(className).
					WithFields(fields...).
					WithNearObject(nearObject).
					WithLimit(2).
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, response)
				require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)
			}
		})

		t.Run("get near text", func(t *testing.T) {
			// Test near text search
			for _, obj := range objects {
				concepts := []string{"some description"}
				_ = obj
				nearText := client.GraphQL().NearTextArgBuilder().
					WithConcepts(concepts).
					WithTargetVectors("colbert")

				response, err := client.GraphQL().Get().
					WithClassName(className).
					WithFields(fields...).
					WithNearText(nearText).
					WithLimit(2).
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, response)
				require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)
			}
		})

		t.Run("get hybrid search", func(t *testing.T) {

			// Test near text search
			for _, obj := range objects {

				objs, err := client.Data().ObjectsGetter().
					WithClassName(className).WithID(obj.ID.String()).WithVector().
					Do(ctx)
				vector := objs[0].Vectors
				hybrid_query := client.GraphQL().HybridArgumentBuilder().
					WithAlpha(1.0).
					WithQuery("The description has ID:" + obj.ID.String()).
					WithTargetVectors("colbert").WithVector(vector)

				response, err := client.GraphQL().Get().
					WithClassName(className).
					WithFields(fields...).
					WithHybrid(hybrid_query).
					WithLimit(2).
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, response)
				require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)

				hybrid_query = client.GraphQL().HybridArgumentBuilder().
					WithAlpha(0.5).
					WithQuery("The description has ID:" + obj.ID.String()).
					WithTargetVectors("colbert").WithVector(vector)

				response, err = client.GraphQL().Get().
					WithClassName(className).
					WithFields(fields...).
					WithHybrid(hybrid_query).
					WithLimit(2).
					Do(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, response)
				require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)

			}
		})
	})

	err := testenv.TearDownLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to tear down weaviate")
	}
}

func TestJinjaAIVectorizerBatchImport(t *testing.T) {
	ctx := context.Background()
	err := testenv.SetupLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to start weaviate")
	}

	client := testsuit.CreateTestClient(true)

	t.Run("batch import", func(t *testing.T) {
		tests := []struct {
			name      string
			hdf5_file string
		}{
			{
				name:      "dataset with 1000 objects",
				hdf5_file: "/Users/rodrigo/Downloads/custom_Multivector_raw-data_lotte-recreation.hdf5",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// clean up DB
				err := client.Schema().AllDeleter().Do(context.Background())
				require.NoError(t, err)

				// create class
				className := "JinjaAIVectorizerBatchImport"

				vectorConfig := map[string]models.VectorConfig{
					"regular": {
						Vectorizer: map[string]interface{}{
							"none": nil,
						},
						VectorIndexType: "hnsw",
					},
					"colbert": {
						Vectorizer: map[string]interface{}{
							"text2colbert-jinaai": map[string]interface{}{
								"properties":         []interface{}{"description"},
								"vectorizeClassName": false,
								"model":              "jina-colbert-v2",
							},
						},
						VectorIndexType: "hnsw",
					},
				}
				class := &models.Class{
					Class: className,
					Properties: []*models.Property{
						{
							Name:     "name",
							DataType: []string{schema.DataTypeText.String()},
						},
						{
							Name:     "description",
							DataType: []string{schema.DataTypeText.String()},
						},
					},
					VectorConfig: vectorConfig,
					ReplicationConfig: &models.ReplicationConfig{
						Factor: 1,
					},
				}
				err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
				require.NoError(t, err)
				// perform batch import
				objects_strings := loadHdf5Documents(tt.hdf5_file, "documents")
				objects_strings = objects_strings[:100]
				objects := make([]*models.Object, len(objects_strings))

				for i := range objects {
					objects[i] = &models.Object{
						ID:    generateUUID(uint32(i)),
						Class: className,
						Properties: map[string]interface{}{
							"name":        "item " + strconv.Itoa(i),
							"description": objects_strings[i],
						},
					}
				}
				// Make batch import in several chunks
				chunk_size := 100
				for i := 0; i < len(objects); i += chunk_size {
					batchResponse, err := client.Batch().ObjectsBatcher().WithObjects(objects[i : i+chunk_size]...).Do(ctx)
					require.NoError(t, err)
					assert.NotNil(t, batchResponse)
					assert.Equal(t, len(objects), len(batchResponse))
				}

				fields := []graphql.Field{
					{Name: "description"},
					{Name: "_additional", Fields: []graphql.Field{
						{Name: "id"},
						{Name: "distance"},
					}},
				}

				time.Sleep(10 * time.Second)

				t.Run("get objects and near vector", func(t *testing.T) {

					// get object
					for _, obj := range objects {
						objs, err := client.Data().ObjectsGetter().
							WithClassName(className).WithID(obj.ID.String()).WithVector().
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, objs)
						require.Len(t, objs, 1)
						assert.Equal(t, obj.ID, objs[0].ID)

						nearVector := client.GraphQL().NearVectorArgBuilder().
							WithVector(objs[0].Vectors["colbert"]).
							WithTargetVectors("colbert")

						response, err := client.GraphQL().Get().
							WithClassName(className).
							WithFields(fields...).
							WithNearVector(nearVector).
							WithLimit(2).
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, response)
						require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)

					}
				})
				t.Run("get near object", func(t *testing.T) {

					// Test near object search
					for _, obj := range objects {
						nearObject := client.GraphQL().NearObjectArgBuilder().WithID(obj.ID.String()).
							WithTargetVectors("colbert")
						response, err := client.GraphQL().Get().
							WithClassName(className).
							WithFields(fields...).
							WithNearObject(nearObject).
							WithLimit(2).
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, response)
						require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)
					}
				})
				t.Run("get near text", func(t *testing.T) {

					queries_strings := loadHdf5Documents(tt.hdf5_file, "queries")
					queries_strings = queries_strings[:10]
					// Test near text search
					for _, query := range queries_strings {
						concepts := []string{query}
						nearText := client.GraphQL().NearTextArgBuilder().
							WithConcepts(concepts).
							WithTargetVectors("colbert")

						response, err := client.GraphQL().Get().
							WithClassName(className).
							WithFields(fields...).
							WithNearText(nearText).
							WithLimit(2).
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, response)
						require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)
					}
				})
				t.Run("get hybrid search with alpha 0.25", func(t *testing.T) {

					// Test near text search
					for _, obj := range objects {

						objs, err := client.Data().ObjectsGetter().
							WithClassName(className).WithID(obj.ID.String()).WithVector().
							Do(ctx)
						vector := objs[0].Vectors
						_ = vector
						hybrid_query := client.GraphQL().HybridArgumentBuilder().
							WithAlpha(0.25).
							WithQuery("The description has ID:" + obj.ID.String()).
							WithTargetVectors("colbert").WithVector(vector)

						response, err := client.GraphQL().Get().
							WithClassName(className).
							WithFields(fields...).
							WithHybrid(hybrid_query).
							WithLimit(2).
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, response)
						require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)
					}
				})

				t.Run("Backup and restore", func(t *testing.T) {

					backup_id := "multi-vector-backup" + strconv.Itoa(int(time.Now().Unix()))

					createResponse, err := client.Backup().Creator().
						WithBackupID(backup_id).
						WithBackend(backup.BACKEND_S3).
						Do(ctx)
					require.NoError(t, err)
					require.NotNil(t, createResponse.Status)
					require.Equal(t, "STARTED", *createResponse.Status)

					for {
						createStatusResponse, err := client.Backup().CreateStatusGetter().
							WithBackupID(backup_id).
							WithBackend(backup.BACKEND_S3).
							Do(ctx)
						require.NoError(t, err)
						require.NotNil(t, createStatusResponse.Status)
						if *createStatusResponse.Status == "SUCCESS" {
							break
						}
						time.Sleep(1 * time.Second)
					}

					fmt.Println("\nDeleting all schemas")
					err = client.Schema().AllDeleter().Do(ctx)
					require.NoError(t, err)
					fmt.Println("Deleted")

					fmt.Println("\nRestoring Backup")
					restoreResponse, err := client.Backup().Restorer().
						WithBackupID(backup_id).
						WithBackend(backup.BACKEND_S3).
						Do(ctx)
					require.NoError(t, err)
					require.NotNil(t, restoreResponse.Status)
					require.Equal(t, "STARTED", *restoreResponse.Status)

					for {
						restoreStatusResponse, err := client.Backup().CreateStatusGetter().
							WithBackupID(backup_id).
							WithBackend(backup.BACKEND_S3).
							Do(ctx)
						require.NoError(t, err)
						require.NotNil(t, restoreStatusResponse.Status)
						if *restoreStatusResponse.Status == "SUCCESS" {
							break
						}
						time.Sleep(1 * time.Second)
					}

					time.Sleep(5 * time.Second)

					// get object and NearVector after backup and restore
					for _, obj := range objects {
						objs, err := client.Data().ObjectsGetter().
							WithClassName(className).WithID(obj.ID.String()).WithVector().
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, objs)
						require.Len(t, objs, 1)
						assert.Equal(t, obj.ID, objs[0].ID)

						nearVector := client.GraphQL().NearVectorArgBuilder().
							WithVector(objs[0].Vectors["colbert"]).
							WithTargetVectors("colbert")

						response, err := client.GraphQL().Get().
							WithClassName(className).
							WithFields(fields...).
							WithNearVector(nearVector).
							WithLimit(2).
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, response)
						require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)

					}
				})

			})
		}
	})

	err = testenv.TearDownLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to tear down weaviate")
	}
}
