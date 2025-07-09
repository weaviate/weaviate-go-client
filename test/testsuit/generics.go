package testsuit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/weaviate/weaviate-go-client/v5/test"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"github.com/weaviate/weaviate/usecases/auth/authorization"
)

const (
	NoAuthPort         = 8080
	NoAuthGRPCPort     = 50051
	AzurePort          = 8081
	OktaCCPort         = 8082
	OktaUsersPort      = 8083
	WCSPort            = 8085
	WCSGRPCPort        = 50056
	NoWeaviatePort     = 8888
	NoWeaviateGRPCPort = 55555
	RbacPort           = 8089
)

const ENV_INTEGRATION_TESTS_AUTH = "INTEGRATION_TESTS_AUTH"

var (
	authEnabled  = os.Getenv(ENV_INTEGRATION_TESTS_AUTH) == "auth_enabled"
	openAIApiKey = os.Getenv("OPENAI_APIKEY")
)

func GetPortAndAuthPw() (int, int, bool) {
	port := NoAuthPort
	grpcPort := NoAuthGRPCPort
	if authEnabled {
		port = WCSPort
		grpcPort = WCSGRPCPort
	}
	return port, grpcPort, authEnabled
}

// CreateWeaviateTestSchemaFood creates a class for each semantic type (Pizza and Soup)
// and adds some primitive properties (name and description)
func CreateWeaviateTestSchemaFood(t *testing.T, client *weaviate.Client, opts ...schemaOptions) {
	createWeaviateTestSchemaFood(t, client, false, opts...)
}

func CreateWeaviateTestSchemaFoodDeprecated(t *testing.T, client *weaviate.Client) {
	createWeaviateTestSchemaFood(t, client, true)
}

func createWeaviateTestSchemaFood(t *testing.T, client *weaviate.Client, isDeprecated bool, opts ...schemaOptions) {
	classes := []*models.Class{
		{
			Class:               "Pizza",
			Description:         "A delicious religion like food and arguably the best export of Italy.",
			InvertedIndexConfig: &models.InvertedIndexConfig{IndexTimestamps: true},
		},
		{
			Class:       "Soup",
			Description: "Mostly water based brew of sustenance for humans.",
		},
		{
			Class:               "Risotto",
			Description:         "Risotto is a northern Italian rice dish cooked with broth.",
			InvertedIndexConfig: &models.InvertedIndexConfig{IndexTimestamps: true},
		},
	}

	for _, class := range classes {
		for _, opt := range opts {
			opt(class)
		}
		err := client.Schema().ClassCreator().WithClass(class).Do(context.Background())
		assert.Nil(t, err)
	}
	namePropertyDataType := []string{"text"}
	if isDeprecated {
		namePropertyDataType = []string{"string"}
	}
	nameProperty := &models.Property{
		DataType:    namePropertyDataType,
		Description: "name",
		Name:        "name",
	}
	bestBeforeProperty := &models.Property{
		DataType:    []string{"date"},
		Description: "You better eat this food before it expires",
		Name:        "best_before",
	}
	descriptionProperty := &models.Property{
		DataType:    []string{"text"},
		Description: "description",
		Name:        "description",
	}
	priceProperty := &models.Property{
		DataType:    []string{"number"},
		Description: "price",
		Name:        "price",
		ModuleConfig: map[string]interface{}{
			"text2vec-contextionary": map[string]interface{}{
				"skip": true,
			},
		},
	}

	propErrT1 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErrT1)
	propErrA1 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErrA1)
	propErrI1 := client.Schema().PropertyCreator().WithClassName("Risotto").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErrI1)
	propErrT2 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrT2)
	propErrA2 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrA2)
	propErrI2 := client.Schema().PropertyCreator().WithClassName("Risotto").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErrI2)
	propErrT3 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(priceProperty).Do(context.Background())
	assert.Nil(t, propErrT3)
	propErrA3 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(priceProperty).Do(context.Background())
	assert.Nil(t, propErrA3)
	propErrI3 := client.Schema().PropertyCreator().WithClassName("Risotto").WithProperty(priceProperty).Do(context.Background())
	assert.Nil(t, propErrI3)
	propErrT4 := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(bestBeforeProperty).Do(context.Background())
	assert.Nil(t, propErrT4)
	propErrA4 := client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(bestBeforeProperty).Do(context.Background())
	assert.Nil(t, propErrA4)
	propErrI4 := client.Schema().PropertyCreator().WithClassName("Risotto").WithProperty(bestBeforeProperty).Do(context.Background())
	assert.Nil(t, propErrI4)
}

func CreateWeaviateTestSchemaWithVectorizorlessClass(t *testing.T, client *weaviate.Client) {
	createWeaviateTestSchemaWithVectorizorlessClass(t, client, false)
}

func CreateWeaviateTestSchemaWithVectorizorlessClassDeprecated(t *testing.T, client *weaviate.Client) {
	createWeaviateTestSchemaWithVectorizorlessClass(t, client, true)
}

func createWeaviateTestSchemaWithVectorizorlessClass(t *testing.T, client *weaviate.Client, isDeprecated bool) {
	vectorizorlessClass := &models.Class{
		Class:       "Donut",
		Description: "A type of leavened fried dough commonly covered with glaze and sprinkles.",
		Vectorizer:  "none",
	}

	err := client.Schema().ClassCreator().WithClass(vectorizorlessClass).Do(context.Background())
	assert.Nil(t, err)

	namePropertyDataType := []string{"text"}
	if isDeprecated {
		namePropertyDataType = []string{"string"}
	}
	nameProperty := &models.Property{
		DataType:    namePropertyDataType,
		Description: "name",
		Name:        "name",
	}
	descriptionProperty := &models.Property{
		DataType:    []string{"text"},
		Description: "description",
		Name:        "description",
	}

	propErr1 := client.Schema().PropertyCreator().WithClassName("Donut").WithProperty(nameProperty).Do(context.Background())
	assert.Nil(t, propErr1)

	propErr2 := client.Schema().PropertyCreator().WithClassName("Donut").WithProperty(descriptionProperty).Do(context.Background())
	assert.Nil(t, propErr2)
}

func CreateWeaviateTestSchemaFoodWithReferenceProperty(t *testing.T, client *weaviate.Client) {
	createWeaviateTestSchemaFoodWithReferenceProperty(t, client, false)
}

func CreateWeaviateTestSchemaFoodWithReferencePropertyDeprecated(t *testing.T, client *weaviate.Client) {
	createWeaviateTestSchemaFoodWithReferenceProperty(t, client, true)
}

// CreateWeaviateTestSchemaFoodWithReferenceProperty create the testing schema with a reference field otherFoods on both classes
func createWeaviateTestSchemaFoodWithReferenceProperty(t *testing.T, client *weaviate.Client, isDeprecated bool) {
	createWeaviateTestSchemaFood(t, client, isDeprecated)
	referenceProperty := &models.Property{
		DataType:    []string{"Pizza", "Soup"},
		Description: "reference to other foods",
		Name:        "otherFoods",
	}
	err := client.Schema().PropertyCreator().WithClassName("Pizza").WithProperty(referenceProperty).Do(context.Background())
	assert.Nil(t, err)
	err = client.Schema().PropertyCreator().WithClassName("Soup").WithProperty(referenceProperty).Do(context.Background())
	assert.Nil(t, err)
}

// CleanUpWeaviate removes the schema and thereby all data
func CleanUpWeaviate(t *testing.T, client *weaviate.Client) {
	ctx := context.Background()

	// Clean up all classes and by that also all data
	err := client.Schema().AllDeleter().Do(ctx)
	assert.Nil(t, err)

	// Cleanup all roles except for the builtin ones.
	roles, err := client.Roles().AllGetter().Do(ctx)
	clientErr := &fault.WeaviateClientError{}
	if err != nil && errors.As(err, &clientErr) && clientErr.StatusCode != -1 {
		t.Logf("delete all roles: %v. This error can be ignored in the 'deprecated' test suite", err)
	}

	for _, role := range roles {
		if name := role.Name; !slices.Contains(authorization.BuiltInRoles, name) {
			client.Roles().Deleter().WithName(role.Name).Do(ctx)
		}
	}
}

// CreateTestClient running on localhost 8080
func CreateTestClient(enableGRPC bool) *weaviate.Client {
	port, grpcPort, authEnabled := GetPortAndAuthPw()

	headers := map[string]string{}
	if openAIApiKey != "" {
		headers["X-OpenAI-Api-Key"] = openAIApiKey
	}

	cfg := weaviate.Config{
		Host:    fmt.Sprintf("localhost:%v", port),
		Scheme:  "http",
		Headers: headers,
	}
	if enableGRPC {
		cfg.GrpcConfig = &grpc.Config{
			Host: fmt.Sprintf("localhost:%v", grpcPort),
		}
	}

	var client *weaviate.Client
	var err error
	if authEnabled {
		cfg.AuthConfig = auth.ApiKey{Value: "my-secret-key"}
		client, err = weaviate.NewClient(cfg)
		if err != nil {
			log.Printf("Error occurred during startup: %v", err)
		}
	} else {
		client = weaviate.New(cfg)
	}
	return client
}

// CreateTestClientForContainer is a test helper that configures an appropriate weaviate.Client for the container.
func CreateTestClientForContainer(t *testing.T, container test.Container) *weaviate.Client {
	t.Helper()

	cfg := weaviate.Config{
		Host:   container.HTTPAddress(),
		Scheme: "http",
	}

	if openAIApiKey != "" {
		cfg.Headers = map[string]string{"X-OpenAI-Api-Key": openAIApiKey}
	}

	if container.EnableGRPC() {
		cfg.GrpcConfig = &grpc.Config{Host: container.GRPCAddress()}
	}

	if container.APISecret != "" {
		cfg.AuthConfig = auth.ApiKey{Value: "my-secret-key"}
	}

	client, err := weaviate.NewClient(cfg)
	require.NoError(t, err, "create test client")
	return client
}

// ParseReferenceResponseToStruct from the interface typed property schema returned by the client
func ParseReferenceResponseToStruct(t *testing.T, reference interface{}) models.MultipleRef {
	referenceList := reference.([]interface{})
	out := make(models.MultipleRef, len(referenceList))
	for i, untyped := range referenceList {
		asMap, ok := untyped.(map[string]interface{})
		assert.True(t, ok)
		beacon, ok := asMap["beacon"]
		assert.True(t, ok)
		beaconString, ok := beacon.(string)
		assert.True(t, ok)
		out[i] = &models.SingleRef{
			Beacon: strfmt.URI(beaconString),
		}
	}
	return out
}

type schemaOptions func(*models.Class)

func WithReplication(cls *models.Class) {
	cls.ReplicationConfig = &models.ReplicationConfig{
		Factor: 2,
	}
}

// CreateTestSchemaAndData with a few pizzas and soups
func CreateTestSchemaAndData(t *testing.T, client *weaviate.Client, opts ...schemaOptions) {
	createTestSchemaAndData(t, client, false, opts...)
}

func CreateTestSchemaAndDataDeprecated(t *testing.T, client *weaviate.Client) {
	createTestSchemaAndData(t, client, true)
}

// CreateTestSchemaAndData with a few pizzas and soups
func createTestSchemaAndData(t *testing.T, client *weaviate.Client, isDeprecated bool, opts ...schemaOptions) {
	createWeaviateTestSchemaFood(t, client, isDeprecated, opts...)

	// Create pizzas
	menuPizza := []*models.Object{
		{
			Class: "Pizza",
			Properties: map[string]interface{}{
				"name":        "Quattro Formaggi",
				"description": "Pizza quattro formaggi Italian: [ˈkwattro forˈmaddʒi] (four cheese pizza) is a variety of pizza in Italian cuisine that is topped with a combination of four kinds of cheese, usually melted together, with (rossa, red) or without (bianca, white) tomato sauce. It is popular worldwide, including in Italy,[1] and is one of the iconic items from pizzerias's menus.",
				"price":       float32(1.1),
				"best_before": "2022-05-03T12:04:40+02:00",
			},
		},
		{
			Class: "Pizza",
			Properties: map[string]interface{}{
				"name":        "Frutti di Mare",
				"description": "Frutti di Mare is an Italian type of pizza that may be served with scampi, mussels or squid. It typically lacks cheese, with the seafood being served atop a tomato sauce.",
				"price":       float32(1.2),
				"best_before": "2022-05-05T07:16:30+02:00",
			},
		},
		{
			Class: "Pizza",
			// this uuid guarantees that it's the first for cursor tests (otherwise
			// they might be flaky if the randomly generated ids are sometimes higher
			// and sometimes lower
			ID: "00000000-0000-0000-0000-000000000000",
			Properties: map[string]interface{}{
				"name":        "Hawaii",
				"description": "Universally accepted to be the best pizza ever created.",
				"price":       float32(1.3),
			},
		},
		{
			ID:    "5b6a08ba-1d46-43aa-89cc-8b070790c6f2",
			Class: "Pizza",
			Properties: map[string]interface{}{
				"name":        "Doener",
				"description": "A innovation, some say revolution, in the pizza industry.",
				"price":       float32(1.4),
			},
		},
	}
	menuSoup := []*models.Object{
		{
			Class: "Soup",
			Properties: map[string]interface{}{
				"name":        "ChickenSoup",
				"description": "Used by humans when their inferior genetics are attacked by microscopic organisms.",
				"price":       float32(2.1),
			},
		},
		{
			Class: "Soup",
			Properties: map[string]interface{}{
				"name":        "Beautiful",
				"description": "Putting the game of letter soups to a whole new level.",
				"price":       float32(2.2),
			},
		},
	}
	menuRisotto := []*models.Object{
		{
			Class: "Risotto",
			Properties: map[string]interface{}{
				"name":        "Risi e bisi",
				"description": "A Veneto spring dish.",
				"price":       float32(3.1),
				"best_before": "2022-05-03T12:04:40+02:00",
			},
		},
		{
			Class: "Risotto",
			Properties: map[string]interface{}{
				"name":        "Risotto alla pilota",
				"description": "A specialty of Mantua, made with sausage, pork, and Parmesan cheese.",
				"price":       float32(3.2),
				"best_before": "2022-05-03T12:04:40+02:00",
			},
		},
		{
			Class: "Risotto",
			ID:    "696bf381-7f98-40a4-bcad-841780e00e0e",
			Properties: map[string]interface{}{
				"name":        "Risotto al nero di seppia",
				"description": "A specialty of the Veneto region, made with cuttlefish cooked with their ink-sacs intact, leaving the risotto black .",
				"price":       float32(3.3),
			},
		},
	}
	thingsBatcher := client.Batch().ObjectsBatcher()
	for _, pizza := range menuPizza {
		thingsBatcher.WithObject(pizza)
	}
	for _, soup := range menuSoup {
		thingsBatcher.WithObject(soup)
	}
	for _, risotto := range menuRisotto {
		thingsBatcher.WithObject(risotto)
	}

	_, thingsErr := thingsBatcher.Do(context.Background())
	assert.Nil(t, thingsErr)
}

func CreateWeaviateTestSchemaDocumentPassage(t *testing.T, client *weaviate.Client) {
	document := &models.Class{
		Class: "Document",
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"text"},
			},
		},
		InvertedIndexConfig: &models.InvertedIndexConfig{IndexTimestamps: true},
	}
	passage := &models.Class{
		Class: "Passage",
		Properties: []*models.Property{
			{
				Name:     "content",
				DataType: []string{"text"},
			},
			{
				Name:     "type",
				DataType: []string{"text"},
			},
			{
				Name:     "ofDocument",
				DataType: []string{"Document"},
			},
		},
	}
	err := client.Schema().ClassCreator().WithClass(document).Do(context.Background())
	assert.Nil(t, err)
	err = client.Schema().ClassCreator().WithClass(passage).Do(context.Background())
	assert.Nil(t, err)
}

func CreateTestDocumentAndPassageSchemaAndData(t *testing.T, client *weaviate.Client) {
	CreateWeaviateTestSchemaDocumentPassage(t, client)

	documentIDs := []string{
		"00000000-0000-0000-0000-00000000000a",
		"00000000-0000-0000-0000-00000000000b",
		"00000000-0000-0000-0000-00000000000c",
		"00000000-0000-0000-0000-00000000000d",
	}
	passageIDs := []string{
		"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
		"00000000-0000-0000-0000-000000000003",
		"00000000-0000-0000-0000-000000000004",
		"00000000-0000-0000-0000-000000000005",
		"00000000-0000-0000-0000-000000000006",
		"00000000-0000-0000-0000-000000000007",
		"00000000-0000-0000-0000-000000000008",
		"00000000-0000-0000-0000-000000000009",
		"00000000-0000-0000-0000-000000000010",
		"00000000-0000-0000-0000-000000000011",
		"00000000-0000-0000-0000-000000000012",
		"00000000-0000-0000-0000-000000000013",
		"00000000-0000-0000-0000-000000000014",
		"00000000-0000-0000-0000-000000000015",
		"00000000-0000-0000-0000-000000000016",
		"00000000-0000-0000-0000-000000000017",
		"00000000-0000-0000-0000-000000000018",
		"00000000-0000-0000-0000-000000000019",
		"00000000-0000-0000-0000-000000000020",
	}
	// Create documents
	documents := make([]*models.Object, len(documentIDs))
	for i, docID := range documentIDs {
		documents[i] = &models.Object{
			ID:    strfmt.UUID(docID),
			Class: "Document",
			Properties: map[string]interface{}{
				"title": fmt.Sprintf("Title of the document %v", i),
			},
		}
	}
	// Create passages
	passages := make([]*models.Object, len(passageIDs))
	for i, passageID := range passageIDs {
		passages[i] = &models.Object{
			ID:    strfmt.UUID(passageID),
			Class: "Passage",
			Properties: map[string]interface{}{
				"content": fmt.Sprintf("Passage content %v", i),
				"type":    "document-passage",
			},
		}
	}

	batcher := client.Batch().ObjectsBatcher()
	for _, document := range documents {
		batcher.WithObject(document)
	}
	for _, passage := range passages {
		batcher.WithObject(passage)
	}
	_, err := batcher.Do(context.Background())
	assert.Nil(t, err)

	createReferences := func(t *testing.T, client *weaviate.Client,
		document *models.Object, passages []*models.Object,
	) {
		ref := client.Data().ReferencePayloadBuilder().
			WithID(document.ID.String()).WithClassName(document.Class).Payload()
		for _, passage := range passages {
			err := client.Data().ReferenceCreator().
				WithID(passage.ID.String()).
				WithClassName(passage.Class).
				WithReferenceProperty("ofDocument").
				WithReference(ref).
				Do(context.TODO())
			assert.Nil(t, err)
		}
	}

	createReferences(t, client, documents[0], passages[:10])
	createReferences(t, client, documents[1], passages[10:14])
}

const (
	AllProperties_RefClass     = "RefClass"
	AllProperties_RefClass2    = "RefClass2"
	AllProperties_RefID1       = "a0000000-0000-0000-0000-000000000001"
	AllProperties_RefID2       = "a0000000-0000-0000-0000-000000000002"
	AllProperties_RefID3       = "a0000000-0000-0000-0000-000000000003"
	AllProperties_hasRefClass  = "hasRefClass"
	AllProperties_hasRefClass2 = "hasRefClass2"
)

func AllPropertiesSchemaCreate(t *testing.T, client *weaviate.Client, className string, withCrossRefs, withMultipleVectors bool) {
	class := &models.Class{
		Class: className,
		Properties: []*models.Property{
			{
				Name:     "color",
				DataType: []string{schema.DataTypeText.String()},
			},
			{
				Name:     "colors",
				DataType: []string{schema.DataTypeTextArray.String()},
			},
			{
				Name:     "author",
				DataType: []string{schema.DataTypeString.String()},
			},
			{
				Name:     "authors",
				DataType: []string{schema.DataTypeStringArray.String()},
			},
			{
				Name:     "number",
				DataType: []string{schema.DataTypeNumber.String()},
			},
			{
				Name:     "numbers",
				DataType: []string{schema.DataTypeNumberArray.String()},
			},
			{
				Name:     "int",
				DataType: []string{schema.DataTypeInt.String()},
			},
			{
				Name:     "ints",
				DataType: []string{schema.DataTypeIntArray.String()},
			},
			{
				Name:     "date",
				DataType: []string{schema.DataTypeDate.String()},
			},
			{
				Name:     "dates",
				DataType: []string{schema.DataTypeDateArray.String()},
			},
			{
				Name:     "bool",
				DataType: []string{schema.DataTypeBoolean.String()},
			},
			{
				Name:     "bools",
				DataType: []string{schema.DataTypeBooleanArray.String()},
			},
			{
				Name:     "uuid",
				DataType: []string{schema.DataTypeUUID.String()},
			},
			{
				Name:     "uuids",
				DataType: []string{schema.DataTypeUUIDArray.String()},
			},
		},
	}

	if withMultipleVectors {
		class.VectorConfig = map[string]models.VectorConfig{
			"author_and_colors": {
				Vectorizer: map[string]interface{}{
					"text2vec-contextionary": map[string]interface{}{
						"vectorizeClassName": false,
						"properties":         []interface{}{"author", "colors"},
					},
				},
				VectorIndexType: "hnsw",
			},
			"all_properties": {
				Vectorizer: map[string]interface{}{
					"text2vec-contextionary": map[string]interface{}{
						"vectorizeClassName": false,
					},
				},
				VectorIndexType: "flat",
			},
		}
	}

	if withCrossRefs {
		refIDs := []string{AllProperties_RefID1, AllProperties_RefID2, AllProperties_RefID3}
		refProperties := []string{"science-fiction", "novel", "fantasy"}
		createRefClass := func(t *testing.T, className string) {
			refClass := &models.Class{
				Class: className,
				Properties: []*models.Property{
					{
						Name:     "category",
						DataType: []string{schema.DataTypeText.String()},
					},
				},
			}
			err := client.Schema().ClassCreator().WithClass(refClass).Do(context.TODO())
			require.Nil(t, err)
		}
		createRefObjects := func(t *testing.T, className string) {
			refObjects := make([]*models.Object, len(refIDs))
			for i := range refObjects {
				refObjects[i] = &models.Object{
					Class: className,
					ID:    strfmt.UUID(refIDs[i]),
					Properties: map[string]interface{}{
						"category": refProperties[i],
					},
				}
			}

			resp, err := client.Batch().ObjectsBatcher().WithObjects(refObjects...).Do(context.TODO())
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp, len(refIDs))
		}

		for _, refClassName := range []string{AllProperties_RefClass, AllProperties_RefClass2} {
			createRefClass(t, refClassName)
			createRefObjects(t, refClassName)
		}

		refProps := []*models.Property{
			{
				Name:     AllProperties_hasRefClass,
				DataType: []string{AllProperties_RefClass},
			},
			{
				Name:     AllProperties_hasRefClass2,
				DataType: []string{AllProperties_RefClass, AllProperties_RefClass2},
			},
		}

		class.Properties = append(class.Properties, refProps...)
	}

	err := client.Schema().ClassCreator().WithClass(class).Do(context.TODO())
	require.Nil(t, err)
}

type AllProperties struct {
	ID1, ID2, ID3 string
	IDs           []string
	Authors       []string
	AuthorsArray  [][]string
	Colors        []string
	ColorssArray  [][]string
	Numbers       []float64
	NumbersArray  [][]float64
	Ints          []int64
	IntsArray     [][]int64
	Uuids         []string
	UuidsArray    [][]string
	Dates         []string
	DatesArray    [][]string
	Bools         []bool
	BoolsArray    [][]bool
}

func AllPropertiesData() AllProperties {
	id1 := "00000000-0000-0000-0000-000000000001"
	id2 := "00000000-0000-0000-0000-000000000002"
	id3 := "00000000-0000-0000-0000-000000000003"
	return AllProperties{
		ID1:     id1,
		ID2:     id2,
		ID3:     id3,
		IDs:     []string{id1, id2, id3},
		Authors: []string{"John", "Jenny", "Joseph"},
		AuthorsArray: [][]string{
			{"John", "Jenny", "Joseph"},
			{"John", "Jenny"},
			{"John"},
		},
		Colors: []string{"red", "blue", "green"},
		ColorssArray: [][]string{
			{"red", "blue", "green"},
			{"red", "blue"},
			{"red"},
		},
		Numbers: []float64{1.1, 2.2, 3.3},
		NumbersArray: [][]float64{
			{1.1, 2.2, 3.3},
			{1.1, 2.2},
			{1.1},
		},
		Ints: []int64{1, 2, 3},
		IntsArray: [][]int64{
			{1, 2, 3},
			{1, 2},
			{1},
		},
		Uuids: []string{id1, id2, id3},
		UuidsArray: [][]string{
			{id1, id2, id3},
			{id1, id2},
			{id1},
		},
		Dates: []string{"2009-11-01T23:00:00Z", "2009-11-02T23:00:00Z", "2009-11-03T23:00:00Z"},
		DatesArray: [][]string{
			{"2009-11-01T23:00:00Z", "2009-11-02T23:00:00Z", "2009-11-03T23:00:00Z"},
			{"2009-11-01T23:00:00Z", "2009-11-02T23:00:00Z"},
			{"2009-11-01T23:00:00Z"},
		},
		Bools: []bool{true, false, true},
		BoolsArray: [][]bool{
			{true, false, true},
			{true, false},
			{true},
		},
	}
}

func AllPropertiesDataAsMap() []map[string]interface{} {
	data := AllPropertiesData()
	id1 := data.ID1
	id2 := data.ID2
	id3 := data.ID3
	ids := []string{id1, id2, id3}

	authors := data.Authors
	authorsArray := data.AuthorsArray
	colors := data.Colors
	colorsArray := data.ColorssArray
	numbers := data.Numbers
	numbersArray := data.NumbersArray
	ints := data.Ints
	intsArray := data.IntsArray
	uuids := data.Uuids
	uuidsArray := data.UuidsArray
	dates := data.Dates
	datesArray := data.DatesArray
	bools := data.Bools
	boolsArray := data.BoolsArray
	properties := make([]map[string]interface{}, len(ids))
	for i := range ids {
		properties[i] = map[string]interface{}{
			"color":   colors[i],
			"colors":  colorsArray[i],
			"author":  authors[i],
			"authors": authorsArray[i],
			"number":  numbers[i],
			"numbers": numbersArray[i],
			"int":     ints[i],
			"ints":    intsArray[i],
			"uuid":    uuids[i],
			"uuids":   uuidsArray[i],
			"date":    dates[i],
			"dates":   datesArray[i],
			"bool":    bools[i],
			"bools":   boolsArray[i],
		}
	}
	return properties
}

func AllPropertiesDataWithCrossReferencesAsMap() []map[string]interface{} {
	return allPropertiesDataWithCrossReferencesAsMap(AllPropertiesDataAsMap())
}

func allPropertiesDataWithCrossReferencesAsMap(properties []map[string]interface{}) []map[string]interface{} {
	// add properties
	// cross references can be declared as []map[string]interface{} or []map[string]string
	for i := range properties {
		properties[i][AllProperties_hasRefClass] = []map[string]interface{}{
			{
				"beacon": fmt.Sprintf("weaviate://localhost/%s/%s", AllProperties_RefClass, AllProperties_RefID1),
			},
		}
		properties[i][AllProperties_hasRefClass2] = []map[string]string{
			{
				"beacon": fmt.Sprintf("weaviate://localhost/%s/%s", AllProperties_RefClass, AllProperties_RefID2),
			},
			{
				"beacon": fmt.Sprintf("weaviate://localhost/%s/%s", AllProperties_RefClass2, AllProperties_RefID3),
			},
		}
	}
	return properties
}

func AllPropertiesDataWithNestedObjectsAsMap() []map[string]interface{} {
	properties := AllPropertiesDataAsMap()
	for i := range properties {
		properties[i]["json"] = map[string]interface{}{
			"firstName":   "Stacey",
			"lastName":    "Spears",
			"proffession": "Accountant",
			"birthdate":   "2011-05-05T07:16:30+02:00",
			"phoneNumber": map[string]interface{}{
				"input":                  "020 1555444",
				"defaultCountry":         "nl",
				"internationalFormatted": "+31 20 1555444",
				"countryCode":            31,
				"national":               201555444,
				"nationalFormatted":      "020 1555444",
				"valid":                  true,
			},
			"location": map[string]interface{}{
				"latitude":  51.366667,
				"longitude": 5.9,
			},
		}
	}
	return properties
}

func AllPropertiesDataWithNestedArrayObjectsAsMap() []map[string]interface{} {
	properties := AllPropertiesDataWithNestedObjectsAsMap()
	for i := range properties {
		properties[i]["people"] = []interface{}{
			map[string]interface{}{
				"firstName":   "Robert",
				"lastName":    "Junior",
				"proffession": "Accountant",
				"birthdate":   "2002-05-05T07:16:30+02:00",
				"phoneNumber": map[string]interface{}{
					"input":          "020 1234567",
					"defaultCountry": "nl",
				},
				"location": map[string]interface{}{
					"latitude":  52.366667,
					"longitude": 4.9,
				},
			},
			map[string]interface{}{
				"firstName":   "Steven",
				"lastName":    "Spears",
				"proffession": "Accountant",
				"birthdate":   "2009-05-05T07:16:30+02:00",
				"phoneNumber": map[string]interface{}{
					"input":          "020 1555444",
					"defaultCountry": "nl",
				},
				"location": map[string]interface{}{
					"latitude":  51.366667,
					"longitude": 5.9,
				},
			},
		}
	}
	return properties
}

func AllPropertiesDataWithCrossReferencesWithNestedArrayObjectsAsMap() []map[string]interface{} {
	return allPropertiesDataWithCrossReferencesAsMap(AllPropertiesDataWithNestedArrayObjectsAsMap())
}

func allPropertiesObjects(className string, properties []map[string]interface{}) []*models.Object {
	data := AllPropertiesData()
	id1 := data.ID1
	id2 := data.ID2
	id3 := data.ID3
	ids := []string{id1, id2, id3}

	objects := make([]*models.Object, len(ids))
	for i, id := range ids {
		objects[i] = &models.Object{
			Class:      className,
			ID:         strfmt.UUID(id),
			Properties: properties[i],
		}
	}
	return objects
}

func AllPropertiesObjectsWithData(className string, properties []map[string]interface{}) []*models.Object {
	return allPropertiesObjects(className, properties)
}

func AllPropertiesObjects(className string) []*models.Object {
	return allPropertiesObjects(className, AllPropertiesDataAsMap())
}

func AllPropertiesObjectsWithNested(className string) []*models.Object {
	return allPropertiesObjects(className, AllPropertiesDataWithNestedObjectsAsMap())
}

func AllPropertiesObjectsWithNestedArray(className string) []*models.Object {
	return allPropertiesObjects(className, AllPropertiesDataWithNestedArrayObjectsAsMap())
}
