package collections

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type (
	Collection struct {
		Name          string
		Description   string
		Properties    []Property
		References    []Reference
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
	Reference struct {
		Name        string
		Collections []string
	}
	ShardingConfig struct {
		VirtualPerPhysical  int
		DesiredCount        int
		DesiredVirtualCount int
	}
	ReplicationConfig struct {
		AsyncEnabled     bool                    // Enable asynchronous replication.
		AsyncReplication *AsyncReplicationConfig // Fine-tune asynchronous replication.
		DeletionStrategy DeletionStrategy        // Conflict resolution strategy for deleted objects.
		Factor           int                     // Number of times a collection is replicated.
	}
	AsyncReplicationConfig api.AsyncReplicationConfig
	InvertedIndexConfig    struct {
		IndexNullState         bool            // Index each object with the null state.
		IndexPropertyLength    bool            // Index length of properties.
		IndexTimestamps        bool            // Index each object by its internal timestamps.
		UsingBlockMaxWAND      bool            // Toggle UsingBlockMaxWAND usage for BM25 search.
		CleanupIntervalSeconds int32           // Asynchronous index cleanup internal.
		BM25                   *BM25Config     // Tuning parameters for the BM25 algorithm.
		Stopwords              *StopwordConfig // Fine-grained control over stopword list usage.
	}
	BM25Config         api.BM25Config
	StopwordConfig     api.StopwordConfig
	MultiTenancyConfig api.MultiTenancyConfig
)

// DataType defines supported property data types.
type DataType api.DataType

const (
	DataTypeText           DataType = DataType(api.DataTypeText)
	DataTypeBool           DataType = DataType(api.DataTypeBool)
	DataTypeInt            DataType = DataType(api.DataTypeInt)
	DataTypeNumber         DataType = DataType(api.DataTypeNumber)
	DataTypeDate           DataType = DataType(api.DataTypeDate)
	DataTypeObject         DataType = DataType(api.DataTypeObject)
	DataTypeGeoCoordinates DataType = DataType(api.DataTypeGeoCoordinates)
	DataTypeTextArray      DataType = DataType(api.DataTypeTextArray)
	DataTypeBoolArray      DataType = DataType(api.DataTypeBoolArray)
	DataTypeIntArray       DataType = DataType(api.DataTypeIntArray)
	DataTypeNumberArray    DataType = DataType(api.DataTypeNumberArray)
	DataTypeDateArray      DataType = DataType(api.DataTypeDateArray)
	DataTypeObjectArray    DataType = DataType(api.DataTypeObjectArray)
)

type Tokenization api.Tokenization

const (
	TokenizationWord       Tokenization = Tokenization(api.TokenizationWord)
	TokenizationWhitespace Tokenization = Tokenization(api.TokenizationWhitespace)
	TokenizationLowercase  Tokenization = Tokenization(api.TokenizationLowercase)
	TokenizationField      Tokenization = Tokenization(api.TokenizationField)
	TokenizationGSE        Tokenization = Tokenization(api.TokenizationGSE)
	TokenizationGSE_CH     Tokenization = Tokenization(api.TokenizationGSE_CH)
	TokenizationTrigram    Tokenization = Tokenization(api.TokenizationTrigram)
	TokenizationKagomeJA   Tokenization = Tokenization(api.TokenizationKagomeJA)
	TokenizationKagomeKR   Tokenization = Tokenization(api.TokenizationKagomeKR)
)

type DeletionStrategy string

const (
	DeleteOnConflict      DeletionStrategy = DeletionStrategy(api.DeleteOnConflict)
	NoAutomatedResolution DeletionStrategy = DeletionStrategy(api.NoAutomatedResolution)
	TimeBasedResolution   DeletionStrategy = DeletionStrategy(api.TimeBasedResolution)
)

func collectionToAPI(c *Collection) api.Collection {
	var properties []api.Property
	if len(c.Properties) > 0 {
		properties = make([]api.Property, len(c.Properties))
		for i, p := range c.Properties {
			properties[i] = api.Property{
				Name:              p.Name,
				Description:       p.Description,
				DataType:          api.DataType(p.DataType),
				NestedProperties:  nestedPropertiesToAPI(p.NestedProperties),
				Tokenization:      api.Tokenization(p.Tokenization),
				IndexFilterable:   p.IndexFilterable,
				IndexRangeFilters: p.IndexRangeFilters,
				IndexSearchable:   p.IndexSearchable,
			}
		}
	}

	var references []api.ReferenceProperty
	if len(c.References) > 0 {
		references = make([]api.ReferenceProperty, len(c.References))
		for i, ref := range c.References {
			references[i] = api.ReferenceProperty{
				Name:        ref.Name,
				Collections: ref.Collections,
			}
		}
	}

	out := api.Collection{
		Name:         c.Name,
		Description:  c.Description,
		Properties:   properties,
		References:   references,
		MultiTenancy: (*api.MultiTenancyConfig)(c.MultiTenancy),
	}

	if c.Sharding != nil {
		out.Sharding = &api.ShardingConfig{
			DesiredCount:        c.Sharding.DesiredCount,
			DesiredVirtualCount: c.Sharding.DesiredVirtualCount,
			VirtualPerPhysical:  c.Sharding.VirtualPerPhysical,
		}
	}

	if c.Replication != nil {
		out.Replication = &api.ReplicationConfig{
			AsyncEnabled:     c.Replication.AsyncEnabled,
			Factor:           c.Replication.Factor,
			DeletionStrategy: api.DeletionStrategy(c.Replication.DeletionStrategy),
			AsyncReplication: (*api.AsyncReplicationConfig)(c.Replication.AsyncReplication),
		}
	}

	if c.InvertedIndex != nil {
		out.InvertedIndex = &api.InvertedIndexConfig{
			IndexNullState:         c.InvertedIndex.IndexNullState,
			IndexPropertyLength:    c.InvertedIndex.IndexPropertyLength,
			IndexTimestamps:        c.InvertedIndex.IndexTimestamps,
			UsingBlockMaxWAND:      c.InvertedIndex.UsingBlockMaxWAND,
			CleanupIntervalSeconds: c.InvertedIndex.CleanupIntervalSeconds,
			Stopwords:              (*api.StopwordConfig)(c.InvertedIndex.Stopwords),
			BM25:                   (*api.BM25Config)(c.InvertedIndex.BM25),
		}
	}

	return out
}

// collectionFromAPI converts api.Collection into Collection.
func collectionFromAPI(c *api.Collection) Collection {
	dev.AssertNotNil(c, "c")

	var properties []Property
	if len(c.Properties) > 0 {
		properties = make([]Property, len(c.Properties))
		for i, p := range c.Properties {
			properties[i] = Property{
				Name:              p.Name,
				Description:       p.Description,
				DataType:          DataType(p.DataType),
				NestedProperties:  nestedPropertiesFromAPI(p.NestedProperties),
				Tokenization:      Tokenization(p.Tokenization),
				IndexFilterable:   p.IndexFilterable,
				IndexRangeFilters: p.IndexRangeFilters,
				IndexSearchable:   p.IndexSearchable,
			}
		}
	}

	var references []Reference
	if len(c.References) > 0 {
		references = make([]Reference, len(c.References))
		for i, ref := range c.References {
			references[i] = Reference{
				Name:        ref.Name,
				Collections: ref.Collections,
			}
		}
	}

	var sharding *ShardingConfig
	if c.Sharding != nil {
		sharding = &ShardingConfig{
			DesiredCount:        c.Sharding.DesiredCount,
			DesiredVirtualCount: c.Sharding.DesiredVirtualCount,
			VirtualPerPhysical:  c.Sharding.VirtualPerPhysical,
		}
	}

	var replication *ReplicationConfig
	if c.Replication != nil {
		replication = &ReplicationConfig{
			AsyncEnabled:     c.Replication.AsyncEnabled,
			Factor:           c.Replication.Factor,
			DeletionStrategy: DeletionStrategy(c.Replication.DeletionStrategy),
			AsyncReplication: (*AsyncReplicationConfig)(c.Replication.AsyncReplication),
		}
	}

	var invertedIndex *InvertedIndexConfig
	if c.InvertedIndex != nil {
		invertedIndex = &InvertedIndexConfig{
			IndexNullState:         c.InvertedIndex.IndexNullState,
			IndexPropertyLength:    c.InvertedIndex.IndexPropertyLength,
			IndexTimestamps:        c.InvertedIndex.IndexTimestamps,
			UsingBlockMaxWAND:      c.InvertedIndex.UsingBlockMaxWAND,
			CleanupIntervalSeconds: c.InvertedIndex.CleanupIntervalSeconds,
			Stopwords:              (*StopwordConfig)(c.InvertedIndex.Stopwords),
			BM25:                   (*BM25Config)(c.InvertedIndex.BM25),
		}
	}

	return Collection{
		Name:          c.Name,
		Description:   c.Description,
		Properties:    properties,
		References:    references,
		Sharding:      sharding,
		Replication:   replication,
		InvertedIndex: invertedIndex,
		MultiTenancy:  (*MultiTenancyConfig)(c.MultiTenancy),
	}
}

func nestedPropertiesFromAPI(nested []api.Property) []Property {
	if len(nested) == 0 {
		return nil
	}

	out := make([]Property, len(nested))
	for i, np := range nested {
		out[i] = Property{
			Name:              np.Name,
			Description:       np.Description,
			DataType:          DataType(np.DataType),
			NestedProperties:  nestedPropertiesFromAPI(np.NestedProperties),
			Tokenization:      Tokenization(np.Tokenization),
			IndexFilterable:   np.IndexFilterable,
			IndexRangeFilters: np.IndexRangeFilters,
			IndexSearchable:   np.IndexSearchable,
		}
	}
	return out
}

func nestedPropertiesToAPI(nested []Property) []api.Property {
	if len(nested) == 0 {
		return nil
	}

	out := make([]api.Property, len(nested))
	for i, p := range nested {
		out[i] = api.Property{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          api.DataType(p.DataType),
			NestedProperties:  nestedPropertiesToAPI(p.NestedProperties),
			Tokenization:      api.Tokenization(p.Tokenization),
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		}
	}
	return out
}
