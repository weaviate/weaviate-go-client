package api

import (
	"net/http"
	"net/url"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

type (
	Collection struct {
		Name          string
		Description   string
		Properties    []Property
		References    []ReferenceProperty
		Sharding      ShardingConfig
		Replication   ReplicationConfig
		InvertedIndex InvertedIndexConfig
		MultiTenancy  MultiTenancyConfig
	}
	Property struct {
		Name              string
		Description       string
		DataType          DataType
		NestedProperties  NestedProperties
		Tokenization      Tokenization
		IndexInverted     bool
		IndexFilterable   bool
		IndexRangeFilters bool
		IndexSearchable   bool
	}
	ReferenceProperty struct {
		Name      string
		DataTypes []string // Collections that can be referenced.
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
		IndexNullState         bool           // Index each object with the null state.
		IndexPropertyLength    bool           // Index length of properties.
		IndexTimestamps        bool           // Index each object by its internal timestamps.
		UsingBlockMaxWAND      bool           // Toggle UsingBlockMaxWAND usage for BM25 search.
		CleanupIntervalSeconds int64          // Asynchronous index cleanup internal.
		BM25                   BM25Config     // Tuning parameters for the BM25 algorithm.
		Stopwords              StopwordConfig // Fine-grained control over stopword list usage.
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

type CreateCollectionRequest struct {
	endpoint
	Collection
}

var _ transport.Endpoint = (*CreateCollectionRequest)(nil)

func (r *CreateCollectionRequest) Method() string { return http.MethodPost }
func (r *CreateCollectionRequest) Path() string   { return "/schema" }
func (r *CreateCollectionRequest) Body() any      { return r.toREST() }

// DeleteCollectionRequest by collection name.
type DeleteCollectionRequest string

var _ transport.Endpoint = (*DeleteCollectionRequest)(nil)

func (d DeleteCollectionRequest) Method() string    { return http.MethodDelete }
func (d DeleteCollectionRequest) Path() string      { return "/schema/" + string(d) }
func (d DeleteCollectionRequest) Query() url.Values { return nil }
func (d DeleteCollectionRequest) Body() any         { return nil }

// toREST repackages Collection into [rest.Class]
// and returns a reference to it to avoid unnecessary copy.
func (c *Collection) toREST() *rest.Class {
	properties := make([]rest.Property, len(c.Properties)+len(c.References))
	for _, p := range c.Properties {
		properties = append(properties, rest.Property{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          []string{string(p.DataType)},
			NestedProperties:  p.NestedProperties.toREST(),
			Tokenization:      rest.PropertyTokenization(p.Tokenization),
			IndexInverted:     p.IndexInverted,
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		})
	}
	for _, ref := range c.References {
		properties = append(properties, rest.Property{
			Name:     ref.Name,
			DataType: ref.DataTypes,
		})
	}

	return &rest.Class{
		Class:       c.Name,
		Description: c.Description,
		Properties:  properties,
		ShardingConfig: map[string]any{
			"desiredCount":        c.Sharding.DesiredCount,
			"desiredVirturlCount": c.Sharding.DesiredVirtualCount,
			"virtualPerPhysical":  c.Sharding.VirtualPerPhysical,
		},
		ReplicationConfig: rest.ReplicationConfig{
			AsyncEnabled:     c.Replication.AsyncEnabled,
			Factor:           c.Replication.Factor,
			DeletionStrategy: rest.ReplicationConfigDeletionStrategy(c.Replication.DeletionStrategy),
		},
		InvertedIndexConfig: rest.InvertedIndexConfig{
			IndexNullState:         c.InvertedIndex.IndexNullState,
			IndexPropertyLength:    c.InvertedIndex.IndexPropertyLength,
			IndexTimestamps:        c.InvertedIndex.IndexTimestamps,
			UsingBlockMaxWAND:      c.InvertedIndex.UsingBlockMaxWAND,
			CleanupIntervalSeconds: c.InvertedIndex.CleanupIntervalSeconds,
			Stopwords:              rest.StopwordConfig(c.InvertedIndex.Stopwords),
			Bm25: rest.BM25Config{
				B:  c.InvertedIndex.BM25.B,
				K1: c.InvertedIndex.BM25.K1,
			},
		},
		MultiTenancyConfig: rest.MultiTenancyConfig(c.MultiTenancy),
	}
}

type NestedProperties []Property

func (nps NestedProperties) toREST() []rest.NestedProperty {
	if len(nps) == 0 {
		return nil
	}

	properties := make([]rest.NestedProperty, len(nps))
	for _, p := range nps {
		properties = append(properties, rest.NestedProperty{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          []string{string(p.DataType)},
			NestedProperties:  p.NestedProperties.toREST(),
			Tokenization:      rest.NestedPropertyTokenization(p.Tokenization),
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		})
	}

	return properties
}
