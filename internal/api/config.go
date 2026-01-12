package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type (
	Collection struct {
		Name          string
		Description   string
		Properties    []Property
		References    []ReferenceProperty
		Sharding      *ShardingConfig
		Replication   *ReplicationConfig
		InvertedIndex *InvertedIndexConfig
		MultiTenancy  *MultiTenancyConfig
	}
	Property struct {
		Name              string
		Description       string
		DataType          DataType
		NestedProperties  NestedProperties
		Tokenization      Tokenization
		IndexFilterable   bool
		IndexRangeFilters bool
		IndexSearchable   bool
	}
	ReferenceProperty struct {
		Name        string
		Collections []string // Collections that can be referenced.
	}
	ShardingConfig struct {
		DesiredCount        int
		DesiredVirtualCount int
		VirtualPerPhysical  int
	}
	ReplicationConfig struct {
		AsyncEnabled     bool             // Enable asynchronous replication.
		Factor           int              // Number of times a collection is replicated.
		DeletionStrategy DeletionStrategy // Conflict resolution strategy for deleted objects.
	}
	InvertedIndexConfig struct {
		IndexNullState         bool            // Index each object with the null state.
		IndexPropertyLength    bool            // Index length of properties.
		IndexTimestamps        bool            // Index each object by its internal timestamps.
		UsingBlockMaxWAND      bool            // Toggle UsingBlockMaxWAND usage for BM25 search.
		CleanupIntervalSeconds int32           // Asynchronous index cleanup internal.
		BM25                   *BM25Config     // Tuning parameters for the BM25 algorithm.
		Stopwords              *StopwordConfig // Fine-grained control over stopword list usage.
	}
	BM25Config         rest.BM25Config
	StopwordConfig     rest.StopwordConfig
	MultiTenancyConfig rest.MultiTenancyConfig
)

// DataType defines supported property data types.
type DataType string

const (
	DataTypeText           DataType = "text"
	DataTypeBool           DataType = "boolean"
	DataTypeInt            DataType = "int"
	DataTypeNumber         DataType = "number"
	DataTypeDate           DataType = "date"
	DataTypeObject         DataType = "object"
	DataTypeGeoCoordinates DataType = "geoCoordinates"
	DataTypeTextArray      DataType = "text[]"
	DataTypeBoolArray      DataType = "boolean[]"
	DataTypeIntArray       DataType = "number[]"
	DataTypeNumberArray    DataType = "date[]"
	DataTypeDateArray      DataType = "object[]"
	DataTypeObjectArray    DataType = "geoCoordinates[]"
)

// knownDataTypes are a set of all data types defined in the Weaviate server.
// A property whose data type is not in knownDataTypes is assumed to be a reference.
var knownDataTypes = NewSet([]DataType{
	DataTypeText,
	DataTypeBool,
	DataTypeInt,
	DataTypeNumber,
	DataTypeDate,
	DataTypeObject,
	DataTypeGeoCoordinates,
	DataTypeTextArray,
	DataTypeBoolArray,
	DataTypeIntArray,
	DataTypeNumberArray,
	DataTypeDateArray,
	DataTypeObjectArray,
})

type Tokenization string

const (
	TokenizationWord       Tokenization = Tokenization(rest.PropertyTokenizationWord)
	TokenizationWhitespace Tokenization = Tokenization(rest.PropertyTokenizationWhitespace)
	TokenizationLowercase  Tokenization = Tokenization(rest.PropertyTokenizationLowercase)
	TokenizationField      Tokenization = Tokenization(rest.PropertyTokenizationField)
	TokenizationGSE        Tokenization = Tokenization(rest.PropertyTokenizationGse)
	TokenizationGSE_CH     Tokenization = Tokenization(rest.PropertyTokenizationGseCh)
	TokenizationTrigram    Tokenization = Tokenization(rest.PropertyTokenizationTrigram)
	TokenizationKagomeJA   Tokenization = Tokenization(rest.PropertyTokenizationKagomeJa)
	TokenizationKagomeKR   Tokenization = Tokenization(rest.PropertyTokenizationKagomeKr)
)

type DeletionStrategy string

const (
	DeleteOnConflict      DeletionStrategy = DeletionStrategy(rest.DeleteOnConflict)
	NoAutomatedResolution DeletionStrategy = DeletionStrategy(rest.NoAutomatedResolution)
	TimeBasedResolution   DeletionStrategy = DeletionStrategy(rest.TimeBasedResolution)
)

// CreateCollectionsRequest creates a new collection in the schema.
type CreateCollectionRequest struct {
	transports.BaseEndpoint
	Collection
}

var _ transports.Endpoint = (*CreateCollectionRequest)(nil)

func (*CreateCollectionRequest) Method() string { return http.MethodPost }
func (*CreateCollectionRequest) Path() string   { return "/schema" }
func (r *CreateCollectionRequest) Body() any    { return &r.Collection }

// GetCollectionRequest by collection name.
var GetCollectionRequest = transports.IdentityEndpoint[string](http.MethodGet, "/schema/%s")

// ListCollectionsRequest fetches definitions for all collections in the schema.
var ListCollectionsRequest transports.Endpoint = transports.StaticEndpoint(http.MethodGet, "/schema")

type ListCollectionsResponse []Collection

var _ json.Unmarshaler = (*ListCollectionsResponse)(nil)

func (r *ListCollectionsResponse) UnmarshalJSON(data []byte) error {
	var schema struct {
		Collections []Collection `json:"classes"`
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		return err
	}
	*r = schema.Collections
	return nil
}

// DeleteCollectionRequest by collection name.
var DeleteCollectionRequest = transports.IdentityEndpoint[string](http.MethodDelete, "/schema/%s")

var (
	_ json.Marshaler   = (*Collection)(nil)
	_ json.Unmarshaler = (*Collection)(nil)
)

func (c *Collection) MarshalJSON() ([]byte, error) {
	properties := make([]rest.Property, 0, len(c.Properties)+len(c.References))
	for _, p := range c.Properties {
		properties = append(properties, rest.Property{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          []string{string(p.DataType)},
			NestedProperties:  p.NestedProperties.toREST(),
			Tokenization:      rest.PropertyTokenization(p.Tokenization),
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		})
	}
	for _, ref := range c.References {
		properties = append(properties, rest.Property{
			Name:     ref.Name,
			DataType: ref.Collections,
		})
	}

	out := &rest.Class{
		Class:       c.Name,
		Description: c.Description,
		Properties:  properties,
	}

	if c.Sharding != nil {
		out.ShardingConfig = map[string]any{
			"desiredCount":        c.Sharding.DesiredCount,
			"desiredVirturlCount": c.Sharding.DesiredVirtualCount,
			"virtualPerPhysical":  c.Sharding.VirtualPerPhysical,
		}
	}

	if c.Replication != nil {
		out.ReplicationConfig = rest.ReplicationConfig{
			AsyncEnabled:     c.Replication.AsyncEnabled,
			Factor:           c.Replication.Factor,
			DeletionStrategy: rest.ReplicationConfigDeletionStrategy(c.Replication.DeletionStrategy),
		}
	}

	if c.InvertedIndex != nil {
		out.InvertedIndexConfig = rest.InvertedIndexConfig{
			IndexNullState:         c.InvertedIndex.IndexNullState,
			IndexPropertyLength:    c.InvertedIndex.IndexPropertyLength,
			IndexTimestamps:        c.InvertedIndex.IndexTimestamps,
			UsingBlockMaxWAND:      c.InvertedIndex.UsingBlockMaxWAND,
			CleanupIntervalSeconds: c.InvertedIndex.CleanupIntervalSeconds,
		}

		if c.InvertedIndex.Stopwords != nil {
			out.InvertedIndexConfig.Stopwords = rest.StopwordConfig(*c.InvertedIndex.Stopwords)
		}

		if c.InvertedIndex.BM25 != nil {
			out.InvertedIndexConfig.Bm25 = rest.BM25Config{
				B:  c.InvertedIndex.BM25.B,
				K1: c.InvertedIndex.BM25.K1,
			}
		}
	}

	if c.MultiTenancy != nil {
		out.MultiTenancyConfig = rest.MultiTenancyConfig(*c.MultiTenancy)
	}

	return json.Marshal(&out)
}

func (c *Collection) UnmarshalJSON(data []byte) error {
	var class rest.Class
	if err := json.Unmarshal(data, &class); err != nil {
		return err
	}

	properties := make([]Property, 0, len(class.Properties))
	references := make([]ReferenceProperty, 0, len(class.Properties))
	for _, p := range class.Properties {
		notReference := len(p.DataType) == 1 && knownDataTypes.Contains(DataType(p.DataType[0]))
		if notReference {
			properties = append(properties, Property{
				Name:              p.Name,
				Description:       p.Description,
				DataType:          DataType(p.DataType[0]),
				NestedProperties:  makeNestedProperties(p.NestedProperties),
				Tokenization:      Tokenization(p.Tokenization),
				IndexFilterable:   p.IndexFilterable,
				IndexRangeFilters: p.IndexRangeFilters,
				IndexSearchable:   p.IndexSearchable,
			})
		} else {
			references = append(references, ReferenceProperty{
				Name:        p.Name,
				Collections: p.DataType,
			})
		}
	}

	var sharding ShardingConfig
	if len(class.ShardingConfig) > 0 {
		// In case any of the fields are not ints, the cast will return a zero value.
		// We explicitly ignore the checks _ to avoid runtime panics in such cases.
		sharding.DesiredCount, _ = class.ShardingConfig["desiredCount"].(int)
		sharding.DesiredVirtualCount, _ = class.ShardingConfig["desiredVirturlCount"].(int)
		sharding.VirtualPerPhysical, _ = class.ShardingConfig["virtualPerPhysical"].(int)
	}

	*c = Collection{
		Name:        class.Class,
		Description: class.Description,
		Properties:  properties,
		References:  references,
		Replication: &ReplicationConfig{
			AsyncEnabled:     class.ReplicationConfig.AsyncEnabled,
			Factor:           class.ReplicationConfig.Factor,
			DeletionStrategy: DeletionStrategy(class.ReplicationConfig.DeletionStrategy),
		},
		InvertedIndex: &InvertedIndexConfig{
			IndexNullState:         class.InvertedIndexConfig.IndexNullState,
			IndexPropertyLength:    class.InvertedIndexConfig.IndexPropertyLength,
			IndexTimestamps:        class.InvertedIndexConfig.IndexTimestamps,
			UsingBlockMaxWAND:      class.InvertedIndexConfig.UsingBlockMaxWAND,
			CleanupIntervalSeconds: class.InvertedIndexConfig.CleanupIntervalSeconds,
			Stopwords:              (*StopwordConfig)(&class.InvertedIndexConfig.Stopwords),
			BM25: &BM25Config{
				B:  class.InvertedIndexConfig.Bm25.B,
				K1: class.InvertedIndexConfig.Bm25.K1,
			},
		},
		Sharding: &sharding,
	}

	return nil
}

type NestedProperties []Property

func makeNestedProperties(nested []rest.NestedProperty) NestedProperties {
	if len(nested) == 0 {
		return nil
	}

	nps := make(NestedProperties, 0, len(nested))
	for _, np := range nested {
		if len(np.DataType) == 1 {
			nps = append(nps, Property{
				Name:              np.Name,
				Description:       np.Description,
				DataType:          DataType(np.DataType[0]),
				NestedProperties:  makeNestedProperties(np.NestedProperties),
				Tokenization:      Tokenization(np.Tokenization),
				IndexFilterable:   np.IndexFilterable,
				IndexRangeFilters: np.IndexRangeFilters,
				IndexSearchable:   np.IndexSearchable,
			})
		} else {
			// Invalid response -- nested property has more than 1 data type.
		}
	}
	return nps
}

func (nps NestedProperties) toREST() []rest.NestedProperty {
	if len(nps) == 0 {
		return nil
	}

	properties := make([]rest.NestedProperty, len(nps))
	for i, p := range nps {
		properties[i] = rest.NestedProperty{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          []string{string(p.DataType)},
			NestedProperties:  p.NestedProperties.toREST(),
			Tokenization:      rest.NestedPropertyTokenization(p.Tokenization),
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		}
	}

	return properties
}
