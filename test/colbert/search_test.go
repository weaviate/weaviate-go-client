package colbert

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"gonum.org/v1/hdf5"
)

func loadVectors(dataset_path string) [][][]float32 {

	vectors := loadHdf5Float32(dataset_path, "vectors")
	ids := loadHdf5Int32(dataset_path, "ids")

	data := make([][][]float32, ids[len(ids)-1]+1)

	id := ids[0]
	for j, vec := range vectors {
		if ids[j] != id {
			id = ids[j]
		}
		data[id] = append(data[id], vec)
	}

	fmt.Printf("Number of documents: %d\n", len(data))
	fmt.Printf("Vectors shape: %d x %d\n", len(vectors), len(vectors[0]))
	return data
}

func loadHdf5Float32(filename string, dataname string) [][]float32 {

	// Open HDF5 file
	file, err := hdf5.OpenFile(filename, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer file.Close()

	// Open dataset
	dataset, err := file.OpenDataset(dataname)
	if err != nil {
		log.Fatalf("Error opening loadHdf5Float32 dataset: %v", err)
	}
	defer dataset.Close()
	dataspace := dataset.Space()
	dims, _, _ := dataspace.SimpleExtentDims()

	byteSize := getHDF5ByteSize(dataset)

	if len(dims) != 2 {
		log.Fatal("expected 2 dimensions")
	}

	rows := dims[0]
	dimensions := dims[1]

	var chunkData [][]float32

	if byteSize == 4 {
		chunkData1D := make([]float32, rows*dimensions)
		dataset.Read(&chunkData1D)
		chunkData = convert1DChunk[float32](chunkData1D, int(dimensions), int(rows))
	} else if byteSize == 8 {
		chunkData1D := make([]float64, rows*dimensions)
		dataset.Read(&chunkData1D)
		chunkData = convert1DChunk[float64](chunkData1D, int(dimensions), int(rows))
	}

	return chunkData
}

func loadHdf5Queries(filename string, dataname string) [][][]float32 {

	// Open HDF5 file
	file, err := hdf5.OpenFile(filename, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer file.Close()

	// Open dataset
	dataset, err := file.OpenDataset(dataname)
	if err != nil {
		log.Fatalf("Error opening queries dataset: %v", err)
	}
	defer dataset.Close()
	dataspace := dataset.Space()
	dims, _, _ := dataspace.SimpleExtentDims()

	byteSize := getHDF5ByteSize(dataset)

	if len(dims) != 3 {
		log.Fatal("expected 3 dimensions")
	}

	num_queries := dims[0]
	num_token := dims[1]
	num_dim := dims[2]

	var chunkData [][][]float32

	if byteSize == 4 {
		chunkData1D := make([]float32, num_queries*num_token*num_dim)
		dataset.Read(&chunkData1D)
		chunkData = convert3DChunk[float32](chunkData1D, int(num_queries), int(num_token), int(num_dim))
	} else if byteSize == 8 {
		chunkData1D := make([]float64, num_queries*num_token*num_dim)
		dataset.Read(&chunkData1D)
		chunkData = convert3DChunk[float64](chunkData1D, int(num_queries), int(num_token), int(num_dim))
	}

	return chunkData
}

func cleanString(s string) string {

	s = strings.TrimRight(s, "\x00")
	s = strings.TrimSpace(s)
	return s
}

func loadHdf5Documents(filename string, dataname string) []string {

	file, err := hdf5.OpenFile(filename, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer file.Close()

	dataset, err := file.OpenDataset(dataname)
	if err != nil {
		log.Fatalf("Error opening loadHdf5Documents dataset: %v", err)
	}
	defer dataset.Close()

	space := dataset.Space()
	dims, _, err := space.SimpleExtentDims()
	_ = dims
	if err != nil {
		log.Fatalf("error getting dimensions: %v", err)
	}

	length := int(dims[0])
	data := make([]string, length)

	if err := dataset.Read(&data); err != nil {
		log.Fatalf("error reading dataset: %v", err)
	}

	// Truncate data when \x00 is found
	for i, s := range data {
		if strings.Contains(s, "\x00") {
			data[i] = s[:strings.Index(s, "\x00")]
		}
	}

	return data

}

func convert3DChunk[D float32 | float64](input []D, num_queries int, num_token int, num_dim int) [][][]float32 {
	chunkData := make([][][]float32, num_queries)
	for i := range chunkData {
		chunkData[i] = make([][]float32, num_token)
		for j := 0; j < num_token; j++ {
			chunkData[i][j] = make([]float32, num_dim)
			for k := 0; k < num_dim; k++ {
				chunkData[i][j][k] = float32(input[i*num_dim*num_token+j*num_dim+k])
			}
		}
	}
	return chunkData
}

func loadHdf5Int32(filename string, dataname string) []uint64 {

	// Open HDF5 file
	file, err := hdf5.OpenFile(filename, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer file.Close()

	// Open dataset
	dataset, err := file.OpenDataset(dataname)
	if err != nil {
		log.Fatalf("Error opening loadHdf5Float32 dataset: %v", err)
	}
	defer dataset.Close()
	dataspace := dataset.Space()
	dims, _, _ := dataspace.SimpleExtentDims()

	if len(dims) != 2 {
		log.Fatal("expected 2 dimensions")
	}

	rows := dims[0]
	dimensions := dims[1]

	var chunkData []uint64

	data_int32 := make([]int32, rows*dimensions)
	dataset.Read(&data_int32)
	chunkData = convertint32toUint64(data_int32)

	return chunkData
}

func convertint32toUint64(input []int32) []uint64 {
	chunkData := make([]uint64, len(input))
	for i := range chunkData {
		chunkData[i] = uint64(input[i])
	}
	return chunkData
}

func getHDF5ByteSize(dataset *hdf5.Dataset) uint {

	datatype, err := dataset.Datatype()
	if err != nil {
		log.Fatalf("Unabled to read datatype\n")
	}

	// log.WithFields(log.Fields{"size": datatype.Size()}).Printf("Parsing HDF5 byte format\n")
	byteSize := datatype.Size()
	if byteSize != 4 && byteSize != 8 {
		log.Fatalf("Unable to load dataset with byte size %d\n", byteSize)
	}
	return byteSize
}

func convert1DChunk[D float32 | float64](input []D, dimensions int, batchRows int) [][]float32 {
	chunkData := make([][]float32, batchRows)
	for i := range chunkData {
		chunkData[i] = make([]float32, dimensions)
		for j := 0; j < dimensions; j++ {
			chunkData[i][j] = float32(input[i*dimensions+j])
		}
	}
	return chunkData
}

func generateUUID(id uint32) strfmt.UUID {
	return strfmt.UUID(fmt.Sprintf("00000000-0000-0000-0000-%012d", id))
}

func loadHdf5Strings(filename string, dataname string) ([]string, error) {
	// Open HDF5 file
	file, err := hdf5.OpenFile(filename, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Open dataset
	dataset, err := file.OpenDataset(dataname)
	if err != nil {
		return nil, fmt.Errorf("error opening dataset: %v", err)
	}
	defer dataset.Close()

	// Get dataspace
	space := dataset.Space()
	dims, _, err := space.SimpleExtentDims()
	if err != nil {
		return nil, fmt.Errorf("error getting dimensions: %v", err)
	}

	// Prepare to read strings
	length := int(dims[0])
	data := make([]string, length)

	// Read each string element
	for i := 0; i < length; i++ {
		var s string
		// Create memory space for single element
		memspace, err := hdf5.CreateSimpleDataspace([]uint{1}, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating memory space: %v", err)
		}

		// Select hyperslab for single element in file space
		space.SelectHyperslab([]uint{uint(i)}, nil, []uint{1}, nil)

		// Read single element
		if err := dataset.ReadSubset(&s, memspace, space); err != nil {
			return nil, fmt.Errorf("error reading element %d: %v", i, err)
		}

		// Clean up the string
		data[i] = strings.TrimRight(s, "\x00")
		memspace.Close()
	}

	return data, nil
}

func TestVectorObjectSearch(t *testing.T) {
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
				hdf5_file: "/Users/rodrigo/Downloads/custom_Multivector_lotte-recreation-reduced_1000_-1.hdf5",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// clean up DB
				err := client.Schema().AllDeleter().Do(context.Background())
				require.NoError(t, err)

				// create class
				className := "NoVectorizer"
				vectorConfig := map[string]models.VectorConfig{
					"regular": {
						Vectorizer: map[string]interface{}{
							"none": nil,
						},
						VectorIndexType: "hnsw",
					},
					"colbert": {
						Vectorizer: map[string]interface{}{
							"none": nil,
						},
						VectorIndexConfig: map[string]interface{}{
							"multivector": map[string]interface{}{
								"enabled": true,
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
					},
					VectorConfig: vectorConfig,
				}
				err = client.Schema().ClassCreator().WithClass(class).Do(ctx)
				require.NoError(t, err)
				// perform batch import
				vectors := loadVectors(tt.hdf5_file)
				objects := make([]*models.Object, len(vectors))
				for i, vec := range vectors {
					objects[i] = &models.Object{
						ID:    generateUUID(uint32(i)),
						Class: className,
						Properties: map[string]interface{}{
							"name": "some name",
						},
						Vectors: models.Vectors{
							"colbert": vec,
						},
					}
				}
				chunk_size := 100
				for i := 0; i < len(objects); i += chunk_size {
					objects_chunk := objects[i:min(i+chunk_size, len(objects))]
					batchResponse, err := client.Batch().ObjectsBatcher().WithObjects(objects_chunk...).Do(ctx)
					require.NoError(t, err)
					assert.NotNil(t, batchResponse)
					assert.Equal(t, len(objects_chunk), len(batchResponse))
				}

				t.Run("get objects", func(t *testing.T) {

					// get object
					for _, obj := range objects {
						objs, err := client.Data().ObjectsGetter().
							WithClassName(className).WithID(obj.ID.String()).WithVector().
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, objs)
						require.Len(t, objs, 1)
						assert.Equal(t, obj.ID, objs[0].ID)
						assert.Equal(t, obj.Vectors["colbert"], objs[0].Vectors["colbert"])
					}
				})

				t.Run("vector search", func(t *testing.T) {

					_additional := graphql.Field{
						Name: "_additional", Fields: []graphql.Field{
							{Name: "distance"},
						},
					}

					// Test vector search
					for _, obj := range objects {

						nearVector := client.GraphQL().NearVectorArgBuilder().
							WithVector(obj.Vectors["colbert"]).WithTargetVectors("colbert")

						response, err := client.GraphQL().Get().
							WithClassName(className).
							WithFields(_additional).
							WithNearVector(nearVector).
							WithLimit(2).
							Do(ctx)
						require.NoError(t, err)
						require.NotEmpty(t, response)
						require.Len(t, response.Data["Get"].(map[string]interface{})[className].([]interface{}), 2)
					}
				})

				t.Run("near object search", func(t *testing.T) {

					_additional := graphql.Field{
						Name: "_additional", Fields: []graphql.Field{
							{Name: "distance"},
						},
					}

					for _, obj := range objects {

						nearObject := client.GraphQL().NearObjectArgBuilder().WithID(obj.ID.String()).WithTargetVectors("colbert")
						response, err := client.GraphQL().Get().
							WithClassName(className).
							WithFields(_additional).
							WithNearObject(nearObject).
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
