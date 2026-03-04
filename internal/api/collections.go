package api

import (
	"encoding/json"
	"net/http"
	"time"

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
		AsyncEnabled     bool                    // Enable asynchronous replication.
		AsyncReplication *AsyncReplicationConfig // Fine-tuning parameters for async replication.
		Factor           int                     // Number of times a collection is replicated.
		DeletionStrategy DeletionStrategy        // Conflict resolution strategy for deleted objects.
	}
	AsyncReplicationConfig struct {
		DiffBatchSize                   int64         // Maximum number of keys in a diff batch.
		DiffPerNodeTimeout              time.Duration // Timeout for computing a diff against a single node. Recommended unit: seconds.
		ReplicationConcurrency          int64         // Maximum number of concurrent replication workers.
		ReplicationFrequency            time.Duration // Frequency at which diff calculations are run. Recommended unit: milliseconds.
		ReplicationFrequencyPropagating time.Duration // Replication frequency during the propagating phase. Recommended unit: milliseconds.
		PrePropagationTimeout           time.Duration // Total timeout for the pre-propagation phase. Recommended unit: seconds.
		PropagationConcurrency          int64         // Maximum number of concurrent propagation workers.
		PropagationBatchSize            int64         // Maximum number of objects in a single propagation batch.
		PropagationLimit                int64         // Maximum number of objects propagated in a single replication round.
		PropagationTimeout              time.Duration // Timeout for a single propagation batch request. Recommended unit: seconds.
		PropagationDelay                time.Duration // Delay before newly added / updated objects are propagated. Recommended unit: milliseconds.
		HashTreeHeight                  int64         // Height of the hash tree used to compute the diff.
		NodePingFrequency               time.Duration // Frequency at which liveness of the target nodes is checked. Recommended unit: milliseconds.
		LoggingFrequency                time.Duration // Frequency at which replication status is logged. Recommended unit: seconds.
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

		if async := c.Replication.AsyncReplication; async != nil {
			out.ReplicationConfig.AsyncConfig = rest.ReplicationAsyncConfig{
				DiffBatchSize:               async.DiffBatchSize,
				DiffPerNodeTimeout:          int64(async.DiffPerNodeTimeout.Seconds()),
				MaxWorkers:                  async.ReplicationConcurrency,
				Frequency:                   async.ReplicationFrequency.Milliseconds(),
				FrequencyWhilePropagating:   async.ReplicationFrequencyPropagating.Milliseconds(),
				PrePropagationTimeout:       int64(async.PrePropagationTimeout.Seconds()),
				PropagationConcurrency:      async.PropagationConcurrency,
				PropagationBatchSize:        async.PropagationBatchSize,
				PropagationLimit:            async.PropagationLimit,
				PropagationTimeout:          int64(async.PropagationTimeout.Seconds()),
				PropagationDelay:            async.PropagationDelay.Milliseconds(),
				HashtreeHeight:              async.HashTreeHeight,
				AliveNodesCheckingFrequency: async.NodePingFrequency.Milliseconds(),
				LoggingFrequency:            int64(async.LoggingFrequency.Seconds()),
			}
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
				NestedProperties:  nestedPropertiesFromREST(p.NestedProperties),
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
		sharding.DesiredCount = int(class.ShardingConfig["desiredCount"].(float64))
		sharding.DesiredVirtualCount = int(class.ShardingConfig["desiredVirtualCount"].(float64))
		sharding.VirtualPerPhysical = int(class.ShardingConfig["virtualPerPhysical"].(float64))
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
			AsyncReplication: &AsyncReplicationConfig{
				DiffBatchSize:                   class.ReplicationConfig.AsyncConfig.DiffBatchSize,
				DiffPerNodeTimeout:              time.Duration(class.ReplicationConfig.AsyncConfig.DiffPerNodeTimeout) * time.Second,
				ReplicationConcurrency:          class.ReplicationConfig.AsyncConfig.MaxWorkers,
				ReplicationFrequency:            time.Duration(class.ReplicationConfig.AsyncConfig.Frequency) * time.Millisecond,
				ReplicationFrequencyPropagating: time.Duration(class.ReplicationConfig.AsyncConfig.FrequencyWhilePropagating) * time.Millisecond,
				PrePropagationTimeout:           time.Duration(class.ReplicationConfig.AsyncConfig.PrePropagationTimeout) * time.Second,
				PropagationConcurrency:          class.ReplicationConfig.AsyncConfig.PropagationConcurrency,
				PropagationBatchSize:            class.ReplicationConfig.AsyncConfig.PropagationBatchSize,
				PropagationLimit:                class.ReplicationConfig.AsyncConfig.PropagationLimit,
				PropagationTimeout:              time.Duration(class.ReplicationConfig.AsyncConfig.PropagationTimeout) * time.Second,
				PropagationDelay:                time.Duration(class.ReplicationConfig.AsyncConfig.PropagationDelay) * time.Millisecond,
				HashTreeHeight:                  class.ReplicationConfig.AsyncConfig.HashtreeHeight,
				NodePingFrequency:               time.Duration(class.ReplicationConfig.AsyncConfig.AliveNodesCheckingFrequency) * time.Millisecond,
				LoggingFrequency:                time.Duration(class.ReplicationConfig.AsyncConfig.LoggingFrequency) * time.Second,
			},
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
		MultiTenancy: &MultiTenancyConfig{
			Enabled:              class.MultiTenancyConfig.Enabled,
			AutoTenantCreation:   class.MultiTenancyConfig.AutoTenantCreation,
			AutoTenantActivation: class.MultiTenancyConfig.AutoTenantActivation,
		},
	}

	return nil
}

func nestedPropertiesFromREST(nested []rest.NestedProperty) []Property {
	if len(nested) == 0 {
		return nil
	}

	nps := make([]Property, 0, len(nested))
	for _, np := range nested {
		if len(np.DataType) != 1 {
			// Invalid response -- nested property must have exactly 1 data type.
			continue
		}

		nps = append(nps, Property{
			Name:              np.Name,
			Description:       np.Description,
			DataType:          DataType(np.DataType[0]),
			NestedProperties:  nestedPropertiesFromREST(np.NestedProperties),
			Tokenization:      Tokenization(np.Tokenization),
			IndexFilterable:   np.IndexFilterable,
			IndexRangeFilters: np.IndexRangeFilters,
			IndexSearchable:   np.IndexSearchable,
		})
	}
	return nps
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
