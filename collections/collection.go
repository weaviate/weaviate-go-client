package collections

import "github.com/weaviate/weaviate-go-client/v6/internal/api"

type (
	Collection struct {
		Name          string
		Description   string
		Properties    []Property
		References    []Reference
		Sharding      ShardingConfig
		Replication   ReplicationConfig
		InvertedIndex InvertedIndexConfig
		MultiTenancy  MultiTenancyConfig
	}
	Property struct {
		Name              string
		Description       string
		DataType          DataType
		NestedProperties  []Property
		Tokenization      Tokenization
		IndexInverted     bool
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
		AsyncEnabled     bool             // Enable asynchronous replication.
		DeletionStrategy DeletionStrategy // Conflict resolution strategy for deleted objects.
		Factor           int              // Number of times a collection is replicated.
	}
	InvertedIndexConfig struct {
		BM25                   BM25Config     // Tuning parameters for the BM25 algorithm.
		Stopwords              StopwordConfig // Fine-grained control over stopword list usage.
		IndexNullState         bool           // Index each object with the null state.
		IndexPropertyLength    bool           // Index length of properties.
		IndexTimestamps        bool           // Index each object by its internal timestamps.
		UsingBlockMaxWAND      bool           // Toggle UsingBlockMaxWAND usage for BM25 search.
		CleanupIntervalSeconds int64          // Asynchronous index cleanup internal.
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

func (r *Collection) toAPI() api.Collection {
	properties := make([]api.Property, len(r.Properties))
	for _, p := range r.Properties {
		properties = append(properties, api.Property{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          api.DataType(p.DataType),
			NestedProperties:  nestedProperties(p.NestedProperties).toAPI(),
			Tokenization:      api.Tokenization(p.Tokenization),
			IndexInverted:     p.IndexInverted,
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		})
	}

	references := make([]api.ReferenceProperty, len(r.References))
	for _, ref := range r.References {
		references = append(references, api.ReferenceProperty{
			Name:        ref.Name,
			Collections: ref.Collections,
		})
	}

	return api.Collection{
		Name:        r.Name,
		Description: r.Description,
		Properties:  properties,
		Sharding: api.ShardingConfig{
			DesiredCount:        r.Sharding.DesiredCount,
			DesiredVirtualCount: r.Sharding.DesiredVirtualCount,
			VirtualPerPhysical:  r.Sharding.VirtualPerPhysical,
		},
		Replication: api.ReplicationConfig{
			AsyncEnabled:     r.Replication.AsyncEnabled,
			Factor:           r.Replication.Factor,
			DeletionStrategy: api.DeletionStrategy(r.Replication.DeletionStrategy),
		},
		InvertedIndex: api.InvertedIndexConfig{
			IndexNullState:         r.InvertedIndex.IndexNullState,
			IndexPropertyLength:    r.InvertedIndex.IndexPropertyLength,
			IndexTimestamps:        r.InvertedIndex.IndexTimestamps,
			UsingBlockMaxWAND:      r.InvertedIndex.UsingBlockMaxWAND,
			CleanupIntervalSeconds: r.InvertedIndex.CleanupIntervalSeconds,
			Stopwords:              api.StopwordConfig(r.InvertedIndex.Stopwords),
			BM25: api.BM25Config{
				B:  r.InvertedIndex.BM25.B,
				K1: r.InvertedIndex.BM25.K1,
			},
		},
		MultiTenancy: api.MultiTenancyConfig(r.MultiTenancy),
	}
}

// fromAPI converts api.Collection into Collection.
func fromAPI(c api.Collection) Collection {
	properties := make([]Property, len(c.Properties))
	for _, p := range c.Properties {
		properties = append(properties, Property{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          DataType(p.DataType),
			NestedProperties:  makeNestedProperties(p.NestedProperties),
			Tokenization:      Tokenization(p.Tokenization),
			IndexInverted:     p.IndexInverted,
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		})
	}
	references := make([]Reference, len(c.Properties))
	for _, ref := range c.References {
		references = append(references, Reference{
			Name:        ref.Name,
			Collections: ref.Collections,
		})
	}

	return Collection{
		Name:        c.Name,
		Description: c.Description,
		Properties:  properties,
		References:  references,
		Replication: ReplicationConfig{
			AsyncEnabled:     c.Replication.AsyncEnabled,
			Factor:           c.Replication.Factor,
			DeletionStrategy: DeletionStrategy(c.Replication.DeletionStrategy),
		},
		InvertedIndex: InvertedIndexConfig{
			IndexNullState:         c.InvertedIndex.IndexNullState,
			IndexPropertyLength:    c.InvertedIndex.IndexPropertyLength,
			IndexTimestamps:        c.InvertedIndex.IndexTimestamps,
			UsingBlockMaxWAND:      c.InvertedIndex.UsingBlockMaxWAND,
			CleanupIntervalSeconds: c.InvertedIndex.CleanupIntervalSeconds,
			Stopwords:              StopwordConfig(c.InvertedIndex.Stopwords),
			BM25: BM25Config{
				B:  c.InvertedIndex.BM25.B,
				K1: c.InvertedIndex.BM25.K1,
			},
		},
		Sharding: ShardingConfig{
			DesiredCount:        c.Sharding.DesiredCount,
			DesiredVirtualCount: c.Sharding.DesiredVirtualCount,
			VirtualPerPhysical:  c.Sharding.VirtualPerPhysical,
		},
	}
}

type nestedProperties []Property

func makeNestedProperties(nested []api.Property) []Property {
	if len(nested) == 0 {
		return nil
	}

	nps := make(nestedProperties, len(nested))
	for _, np := range nested {
		nps = append(nps, Property{
			Name:              np.Name,
			Description:       np.Description,
			DataType:          DataType(np.DataType),
			NestedProperties:  makeNestedProperties(np.NestedProperties),
			Tokenization:      Tokenization(np.Tokenization),
			IndexFilterable:   np.IndexFilterable,
			IndexRangeFilters: np.IndexRangeFilters,
			IndexSearchable:   np.IndexSearchable,
		})
	}
	return nps
}

func (nps nestedProperties) toAPI() []api.Property {
	if len(nps) == 0 {
		return nil
	}

	properties := make([]api.Property, len(nps))
	for _, p := range nps {
		properties = append(properties, api.Property{
			Name:              p.Name,
			Description:       p.Description,
			DataType:          api.DataType(p.DataType),
			NestedProperties:  nestedProperties(p.NestedProperties).toAPI(),
			Tokenization:      api.Tokenization(p.Tokenization),
			IndexFilterable:   p.IndexFilterable,
			IndexRangeFilters: p.IndexRangeFilters,
			IndexSearchable:   p.IndexSearchable,
		})
	}

	return properties
}
