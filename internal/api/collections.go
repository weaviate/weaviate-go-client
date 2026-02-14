package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
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
		NestedProperties  []Property
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
var knownDataTypes = newSet([]DataType{
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

// DeleteCollectionRequest by collection name.
var DeleteCollectionRequest = transports.IdentityEndpoint[string](http.MethodDelete, "/schema/%s")

var _ json.Marshaler = (*Collection)(nil)

// MarshaJSON marshals Collection via [rest.Class].
func (c *Collection) MarshalJSON() ([]byte, error) {
	properties := make([]rest.Property, len(c.Properties)+len(c.References))
	for i, p := range c.Properties {
		properties[i] = rest.Property{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          []string{string(p.DataType)},
			NestedProperties:  nestedPropertiesToREST(p.NestedProperties),
			Tokenization:      rest.PropertyTokenization(p.Tokenization),
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		}
	}
	for i, ref := range c.References {
		properties[i+len(c.Properties)] = rest.Property{
			Name:     ref.Name,
			DataType: ref.Collections,
		}
	}

	out := &rest.Class{
		Class:       c.Name,
		Description: c.Description,
		Properties:  properties,
	}

	if c.Sharding != nil {
		out.ShardingConfig = map[string]any{
			"desiredCount":        c.Sharding.DesiredCount,
			"desiredVirtualCount": c.Sharding.DesiredVirtualCount,
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

func nestedPropertiesToREST(nps []Property) []rest.NestedProperty {
	if len(nps) == 0 {
		return nil
	}

	properties := make([]rest.NestedProperty, len(nps))
	for i, p := range nps {
		properties[i] = rest.NestedProperty{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          []string{string(p.DataType)},
			NestedProperties:  nestedPropertiesToREST(p.NestedProperties),
			Tokenization:      rest.NestedPropertyTokenization(p.Tokenization),
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		}
	}

	return properties
}

func newSet[Slice ~[]E, E comparable](values Slice) set[E] {
	set := make(set[E], len(values))
	for _, v := range values {
		set.Add(v)
	}
	return set
}

// set is a lightweight set implementation based on on map[string]struct{}.
// It uses struct{} as a value type, because empty structs do not use memory.
type set[E comparable] map[E]struct{}

func (s set[E]) Add(v E) {
	s[v] = struct{}{}
}

func (s set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}
